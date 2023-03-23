package model

import (
	"context"
	"encoding/json"
	"errors"
	"go-netdisk/conf"
	"go-netdisk/mysql"
	"go-netdisk/redis"
	"go-netdisk/util"
	"log"
	"time"
)

// User struct
type User struct {
	Id           int `gorm:"primary_key;type:auto_increment;"`
	FileStoreId  string
	UserName     string
	Password     string
	RegisterTime time.Time
	ImagePath    string
}

// 序列化
func (u *User) MarshalBinary() (data []byte, err error) {
	return json.Marshal(u)
}

// 反序列化
func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

// 检查名字是否存在
func IfExistName(name string) bool {
	var user []User
	var count int64
	result := mysql.Db.Table("go_netdisk_user").Find(&user, "user_name = ?", name).Count(&count)
	if result.Error == nil && count == 1 {
		return true
	}
	return false
}

// 创建新用户并存入Redis
func CreateUser(name, pwd, token string, ctx context.Context) (User, error) {
	user := User{UserName: name, Password: pwd, RegisterTime: time.Now(), FileStoreId: util.GetUUID(), ImagePath: "/static/img/user.jpg"}
	result := mysql.Db.Table("go_netdisk_user").Create(&user)
	if result.Error == nil {
		u, _ := GetUser(ctx, "", name, pwd)
		redis.Rdb.Set(ctx, token, &u, time.Duration(conf.RedisConfig.TTL)*time.Second).Err()
		return user, nil
	}
	log.Fatal(result.Error)
	return User{}, result.Error
}

// 获取用户实例
func GetUser(ctx context.Context, token, name, pwd string) (User, error) {
	val, err := redis.Rdb.Get(ctx, token).Result()
	if err == nil {
		user := new(User)
		user.UnmarshalBinary([]byte(val))
		return *user, err
	}
	var user User
	result := mysql.Db.Table("go_netdisk_user").Find(&user, "user_name = ?", name)
	if result.Error == nil && result.RowsAffected == 1 && user.Password == pwd {
		return user, nil
	}
	return user, errors.New("error")
}

// 检查登陆状态
func CheckToken(ctx context.Context, token string) bool {
	val, err := redis.Rdb.Exists(ctx, token).Result()
	if val == int64(1) && err == nil {
		redis.UpdateTTL(ctx, token)
		return true
	}
	return false
}

// 在Reids中设置key为token值为user的json的键值对
func SetTokenUser(ctx context.Context, token string, user User) error {
	err := redis.Rdb.Set(ctx, token, &user, time.Duration(conf.RedisConfig.TTL)*time.Second).Err()
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// 更新Redis键值对
func UpdateTokenUser(ctx context.Context, token string, new_token string) error {
	err := redis.Rdb.Rename(ctx, token, new_token).Err()
	if err != nil {
		log.Fatal(err)
		return err
	}
	redis.UpdateTTL(ctx, token)
	return nil
}

// 删除Redis键值对
func DelTokenUser(ctx context.Context, token string) error {
	err := redis.Rdb.Del(ctx, token).Err()
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
