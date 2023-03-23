package main

import (
	"go-netdisk/conf"
	"go-netdisk/mysql"
	"go-netdisk/redis"
	"go-netdisk/route"
)

func main() {
	conf.LoadConf()
	redis.SetupRedis()
	mysql.SetupDB()
	route.SetupRoute()

	defer func() {
		redis.CloseRedis()
	}()
}