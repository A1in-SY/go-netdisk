package model

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

//初始化Mysql数据库
func SetupDB() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/go_netdisk?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("连接Mysql数据库失败...")
	}
}
