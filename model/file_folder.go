package model

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-netdisk/conf"
	"go-netdisk/mysql"
	"go-netdisk/redis"
	"go-netdisk/util"
	"io/fs"
	"log"
	"os"
	"strconv"
	"time"
)

type FileFolder struct {
	Id             int `gorm:"primary_key;type:auto_increment;"`
	FileFolderName string
	ParentFolderId int
	FileStoreId    string
	Time           string
}

type FileFolders struct {
	Folders []FileFolder
}

func (ff *FileFolder) MarshalBinary() (data []byte, err error) {
	return json.Marshal(ff)
}

func (ff *FileFolder) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, ff)
}

func (ffs *FileFolders) MarshalBinary() (data []byte, err error) {
	return json.Marshal(ffs)
}

func (ffs *FileFolders) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, ffs)
}

func CreateUserRootFileFolder(id int, fileStoreId string) (FileFolder, error) {
	fileRootFolder := FileFolder{FileFolderName: fileStoreId, ParentFolderId: -1, FileStoreId: fileStoreId, Time: util.GetNowTime()}
	result := mysql.Db.Table("go_netdisk_file_folder").Create(&fileRootFolder)
	if result.Error == nil && result.RowsAffected == 1 {
		err := os.Mkdir(conf.AppConfig.Upload_location+"/"+fileRootFolder.FileFolderName, fs.ModeDir)
		if err != nil {
			log.Fatal(err)
			return fileRootFolder, err
		}
		return fileRootFolder, nil
	}
	log.Fatal(result.Error)
	return FileFolder{}, result.Error
}

func GetUserFileFolderCount(ctx context.Context, token string, fileStoreId string) (int64, error) {
	val, err := redis.Rdb.Get(ctx, "FileFolderCount_"+token).Result()
	if err == nil {
		count, e := strconv.ParseInt(val, 10, 64)
		if e == nil {
			return count, err
		}
	}
	var fileFolders FileFolders
	var fileFolderCount int64
	result := mysql.Db.Table("go_netdisk_file_folder").Find(&fileFolders.Folders, "file_store_id = ?", fileStoreId).Count(&fileFolderCount)
	if result.Error == nil {
		redis.Rdb.Set(ctx, "FileFolder_"+token, &fileFolders, time.Duration(conf.RedisConfig.TTL)*time.Second).Err()
		redis.Rdb.Set(ctx, "FileFolderCount_"+token, &fileFolderCount, time.Duration(conf.RedisConfig.TTL)*time.Second).Err()
		return fileFolderCount, nil
	}
	return fileFolderCount, result.Error
}

func GetUserFileFolders(ctx context.Context, token string, fileStoreId string, skipRedis bool) (FileFolders, error) {
	if !skipRedis {
		val, err := redis.Rdb.Get(ctx, "FileFolder_"+token).Result()
		if err == nil {
			fileFolders := new(FileFolders)
			fileFolders.UnmarshalBinary([]byte(val))
			return *fileFolders, nil
		}
	}
	var fileFolders FileFolders
	result := mysql.Db.Table("go_netdisk_file_folder").Find(&fileFolders.Folders, "file_store_id = ?", fileStoreId)
	if result.Error == nil {
		redis.Rdb.Set(ctx, "FileFolder_"+token, &fileFolders, time.Duration(conf.RedisConfig.TTL)*time.Second).Err()
		return fileFolders, nil
	}
	return fileFolders, result.Error
}

func (ffs *FileFolders) SelectFoldersByParentId(id int) FileFolders {
	var fileFolders FileFolders
	for _, ff := range ffs.Folders {
		if ff.ParentFolderId == id {
			fileFolders.Folders = append(fileFolders.Folders, ff)
		}
	}
	return fileFolders
}

func (ffs *FileFolders) SelectFolderById(id int) FileFolder {
	var fileFolder FileFolder
	for _, ff := range ffs.Folders {
		if ff.Id == id {
			fileFolder = ff
			return fileFolder
		}
	}
	fileFolder.Id = 0
	fileFolder.FileFolderName = ""
	return fileFolder
}

func CreateFileFolder(ctx context.Context, token, fileFolderName string, parentFolderId int) (FileFolders, error) {
	val, err1 := redis.Rdb.Get(ctx, "FileFolder_"+token).Result()
	fileFolders := new(FileFolders)
	fileFolders.UnmarshalBinary([]byte(val))
	fileStoreId := fileFolders.SelectFoldersByParentId(-1).Folders[0].FileStoreId
	if parentFolderId == 0 {
		for _, ff := range fileFolders.Folders {
			if ff.ParentFolderId == -1 {
				parentFolderId = ff.Id
			}
		}
	}
	newFileFolder := FileFolder{FileFolderName: fileFolderName, ParentFolderId: parentFolderId, FileStoreId: fileStoreId, Time: util.GetNowTime()}
	result := mysql.Db.Table("go_netdisk_file_folder").Create(&newFileFolder)
	if result.Error == nil && err1 == nil {
		return GetUserFileFolders(ctx, token, fileStoreId, true)
	}
	return *fileFolders, errors.New("error")
}

func GetAllParents(ctx context.Context, token string, id int) []FileFolder {
	if id == 0 {
		return []FileFolder{}
	} else {
		val, _ := redis.Rdb.Get(ctx, "FileFolder_"+token).Result()
		fileFolders := new(FileFolders)
		fileFolders.UnmarshalBinary([]byte(val))
		allParentsFolders := []FileFolder{}
		temp := id
		for {
			ff := fileFolders.SelectFolderById(temp)
			if ff.ParentFolderId == -1 {
				break
			}
			fmt.Printf("temp: %v\n", temp)
			fmt.Printf("ff.ParentFolderId: %v\n", ff.ParentFolderId)
			allParentsFolders = append(allParentsFolders, fileFolders.SelectFolderById(ff.Id))
			temp = ff.ParentFolderId
		}
		for i, j := 0, len(allParentsFolders)-1; i < j; i, j = i+1, j-1 {
			allParentsFolders[i], allParentsFolders[j] = allParentsFolders[j], allParentsFolders[i]
		}
		return allParentsFolders[:len(allParentsFolders)-1]
	}
}
