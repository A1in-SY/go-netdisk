package mysql

import (
	"fmt"
	"go-netdisk/conf"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

// 初始化Mysql数据库
func SetupDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.MysqlConfig.Name, conf.MysqlConfig.Password, conf.MysqlConfig.Host, conf.MysqlConfig.Port, conf.MysqlConfig.Db)
	fmt.Printf("dsn: %v\n", dsn)
	var err error
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("连接Mysql数据库失败...")
	}
}
