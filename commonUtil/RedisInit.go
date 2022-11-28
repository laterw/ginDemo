package commonUtil

import (
	"flow/redis"
	"fmt"
)

// SetupRedis 初始化 Redis
func SetupRedis() {
	//viper.Get("redis.port")
	// 建立 Redis 连接
	redis.ConnectRedis(
		fmt.Sprintf("%v:%v", "127.0.01", "6379"),
		"",
		"",
		5,
	)
}
