package conf

import (
	"fmt"

	"github.com/spf13/viper"
)

type Conf struct {
}

type mysqlConf struct {
	Host     string
	Port     string
	Name     string
	Password string
	Db       string
}

type redisConf struct {
	Host     string
	Port     string
	Password string
	Db       int
	TTL      int
}

type appConf struct {
	Upload_location string
}

var MysqlConfig mysqlConf
var RedisConfig redisConf
var AppConfig appConf

func LoadConf() {
	globalConfig := viper.New()
	globalConfig.SetConfigFile("config.yaml")
	globalConfig.AddConfigPath("./")
	err := globalConfig.ReadInConfig() // 读取配置信息
	if err != nil {                    // 读取配置信息失败
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	MysqlConfig.Host = globalConfig.GetString("mysql.host")
	MysqlConfig.Port = globalConfig.GetString("mysql.port")
	MysqlConfig.Name = globalConfig.GetString("mysql.name")
	MysqlConfig.Password = globalConfig.GetString("mysql.password")
	MysqlConfig.Db = globalConfig.GetString("mysql.db")

	RedisConfig.Host = globalConfig.GetString("redis.host")
	RedisConfig.Port = globalConfig.GetString("redis.port")
	RedisConfig.Password = globalConfig.GetString("redis.password")
	RedisConfig.Db = globalConfig.GetInt("redis.db")
	RedisConfig.TTL = globalConfig.GetInt("redis.ttl")

	AppConfig.Upload_location = globalConfig.GetString("app.upload_location")
}
