package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id           int
	OpenId       string
	FileStoreId  int
	UserName     string
	Password     string
	RegisterTime time.Time
	ImagePath    string
}

func GetUser(name, pwd string) bool {
	var user User
	result := Db.Table("go_netdisk_user").Where("user_name = ?", name).First(&user)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) && user.Password == pwd{
		return true
	}
	return false
}
