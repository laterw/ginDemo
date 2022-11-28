package model

import (
	"flag"
	"flow/commonUtil"
	result "flow/entity"
	"flow/redis"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
)

func GetflowTatal() result.Result {
	DB := commonUtil.GetDB()
	var result result.Result
	DB.Raw("select count(1) as total from act_hi_procinst aht   ").Scan(&result)
	return result
}

func GetOldTotal(c *gin.Context) result.Result {
	dateTime := c.Param("dateTime")
	DB := commonUtil.GetDB()
	var result result.Result
	DB.Raw("select count(1) as total from act_hi_procinst aht where date_format(aht.START_TIME_,'%Y-%m-%d') =  ?", dateTime).Scan(&result)
	return result
}

func GetLeft1(c *gin.Context) []result.Result {
	DB := commonUtil.GetDB()
	var result []result.Result
	DB.Raw("select TENANT_ID_ as tenant ,count(1) as total from act_hi_procinst ahp   group by ahp.TENANT_ID_ ").Scan(&result)
	return result
}
func GetLeft2(c *gin.Context) []result.Result {
	DB := commonUtil.GetDB()
	var result []result.Result
	DB.Raw("select START_USER_ID_ as start_user,count(1) as total from act_hi_procinst ahp  group by ahp.START_USER_ID_ ").Scan(&result)
	for i := range result {
		result[i].StartUser = redis.Redis.Get(result[i].StartUser)
	}
	return result
}

func GetLeft3(c *gin.Context) result.Result {
	busssinessKey := c.Param("busssinessKey")
	DB := commonUtil.GetDB()
	var result result.Result
	DB.Debug().Raw("select START_USER_ID_  as start_user,date_format(START_TIME_,'%Y-%m-%d %H:%i:%s') as start_time,id_  as id, BUSINESS_KEY_  as buss_key from act_hi_procinst ahp where START_TIME_ > (select ahp2.START_TIME_  from act_hi_procinst ahp2 where ahp2.BUSINESS_KEY_=?)  order by ahp.START_TIME_  desc  limit 1 ", busssinessKey).Scan(&result)
	return result
}

func GetLeft3Init(c *gin.Context) []result.Result {
	DB := commonUtil.GetDB()
	var result []result.Result
	DB.Debug().Raw("select START_USER_ID_  as start_user,date_format(START_TIME_,'%Y-%m-%d %H:%i:%s') as start_time,id_  as id, BUSINESS_KEY_  as buss_key from act_hi_procinst ahp  order by ahp.START_TIME_  desc  limit 9 ").Scan(&result)
	for i := range result {
		result[i].StartUser = redis.Redis.Get(result[i].StartUser)
	}
	return result
}

func TestWebs(c *gin.Context) {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	var dialer *websocket.Dialer
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.WriteMessage(1, []byte("heartbeat"))
	if err != nil {
		fmt.Println("read:", err)
		return
	}
}

var addr = flag.String("addr", "localhost:8888", "http service address")

var upgrader = websocket.Upgrader{
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options
