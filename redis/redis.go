package redis

import "github.com/redis/go-redis/v9"

var Rdb *redis.Client

//初始化Redis
func SetupRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // 没有密码，默认值
		DB:       0,  // 默认DB 0
	})
}
