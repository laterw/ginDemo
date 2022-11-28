package main

import (
	"flow/commonUtil"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	_ "gorm.io/driver/mysql"
	"log"
	"net/http"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		// 读取存储空间大小
		ReadBufferSize: 1024,
		// 写入存储空间大小
		WriteBufferSize: 1024,
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 20 * time.Second,
	// 取消 ws 跨域校验
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandler(c *gin.Context) {
	var conn *websocket.Conn
	var err error

	conn, err = upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}

	userId := 123456

	// 将 conn 保存到字典
	commonUtil.AddClient(uint(userId), conn)
}

func main() {

	go commonUtil.InitDB()     // 初始化数据库
	go commonUtil.InitCanal()  // 初始化 canal
	go commonUtil.SetupRedis() //初始化redis
	r := gin.Default()         // 获取 gin
	r = CollectRouter(r)       //获取 所有路由,这里使用简单 路由
	r.GET("/ws", WsHandler)
	r.Run(":8888")
}
