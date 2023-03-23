package redis

import (
	"context"
	"go-netdisk/conf"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

//初始化Redis
func SetupRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     conf.RedisConfig.Host+":"+conf.RedisConfig.Port,
		Password: conf.RedisConfig.Password, // 没有密码，默认值
		DB:       conf.RedisConfig.Db,  // 默认DB 0
	})
}

func CloseRedis() {
	Rdb.Close()
}

func UpdateTTL(ctx context.Context, token string) error {
	err := Rdb.Expire(ctx, token, time.Duration(conf.RedisConfig.TTL)*time.Second).Err()
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}