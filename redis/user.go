package redis

import (
	"context"
	"log"
	"time"
)

func GetUser(ctx context.Context, name string, pwd string) bool {
	val, err := Rdb.Get(ctx, name).Result()
	if err == nil && val == pwd {
		return true
	}
	return false
}

func SetUser(ctx context.Context, name string, pwd string) {
	err := Rdb.Set(ctx, name, pwd, 3600*4*time.Second).Err()
	if err != nil {
		log.Fatal(err)
		return
	}
}