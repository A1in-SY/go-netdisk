package model

import (
	"context"
	"encoding/json"
	"go-netdisk/conf"
	"go-netdisk/mysql"
	"go-netdisk/redis"
	"log"
	"time"
)

type FileStore struct {
	Id          string `gorm:"primaryKey"`
	UserId      int
	CurrentSize int64 `gorm:"default:0;"`
	MaxSize     int64 `gorm:"default:1048576;"`
}

func (fs *FileStore) MarshalBinary() (data []byte, err error) {
	return json.Marshal(fs)
}

func (fs *FileStore) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, fs)
}

func CreateUserFileStore(id int, fileStoreId string) (FileStore, error) {
	fileStore := FileStore{Id: fileStoreId, UserId: id}
	result := mysql.Db.Table("go_netdisk_file_store").Create(&fileStore)
	if result.Error == nil && result.RowsAffected == 1 {
		return fileStore, nil
	}
	log.Fatal(result.Error)
	return FileStore{}, result.Error
}

func GetUserFileStore(ctx context.Context, token string, id int) (FileStore, error) {
	val, err := redis.Rdb.Get(ctx, "FileStore_"+token).Result()
	if err == nil {
		fileStore := new(FileStore)
		fileStore.UnmarshalBinary([]byte(val))
		return *fileStore, err
	}
	var fileStore FileStore
	result := mysql.Db.Table("go_netdisk_file_store").Find(&fileStore, "user_id = ?", id)
	if result.Error == nil {
		redis.Rdb.Set(ctx, "FileStore_"+token, &fileStore, time.Duration(conf.RedisConfig.TTL)*time.Second).Err()
		return fileStore, nil
	}
	log.Fatal(result.Error)
	return fileStore, result.Error
}

func SubtractSize(ctx context.Context, token string, fileSize int64, fileStoreId string) {
	var fileStore FileStore
	fileStore.Id = fileStoreId
	mysql.Db.Table("go_netdisk_file_store").First(&fileStore)
	fileStore.CurrentSize = fileStore.CurrentSize + fileSize/1024
	fileStore.MaxSize = fileStore.MaxSize - fileSize/1024
	mysql.Db.Table("go_netdisk_file_store").Save(&fileStore)
}