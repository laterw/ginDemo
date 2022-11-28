package commonUtil

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

var (
	clients = make(map[uint]map[*websocket.Conn]bool)
	mux     sync.Mutex
)

func AddClient(id uint, conn *websocket.Conn) {
	mux.Lock()
	if clients[id] == nil {
		clients[id] = make(map[*websocket.Conn]bool)
	}
	clients[id][conn] = true
	mux.Unlock()
}

func DeleteClient(id uint, conn *websocket.Conn) {
	mux.Lock()
	_ = conn.Close()
	delete(clients[id], conn)
	mux.Unlock()
}

func GetClients(id uint) (conns []*websocket.Conn) {
	mux.Lock()
	_conns, ok := clients[id]
	if ok {
		for k := range _conns {
			conns = append(conns, k)
		}
	}
	mux.Unlock()
	return
}
func IsSave(id uint) bool {
	mux.Lock()
	conns := clients[id]
	mux.Unlock()
	return conns == nil
}

func SetMessage(userId uint, content interface{}) {
	conns := GetClients(userId)
	for i := range conns {
		i := i
		err := conns[i].WriteJSON(content)
		if err != nil {
			log.Println("write json err:", err)
			DeleteClient(userId, conns[i])
		}
	}
}
