package main

import (
	"go-netdisk/model"
	"go-netdisk/redis"
	"go-netdisk/route"
)

func main() {
	redis.SetupRedis()
	model.SetupDB()
	route.SetupRoute()
}