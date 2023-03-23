package model

import (
	"context"
	"encoding/json"
	"errors"
	"go-netdisk/conf"
	"go-netdisk/mysql"
	"go-netdisk/redis"
	"go-netdisk/util"
	"path"
	"strconv"
	"strings"
	"time"
)

type File struct {
	Id             int    //文件id
	FileName       string //文件名
	FileHash       string //文件哈希值
	FileStoreId    string //文件仓库id
	FilePath       string //文件存储路径
	DownloadNum    int    //下载次数
	UploadTime     string //上传时间
	ParentFolderId int    //父文件夹id
	Size           int64  //文件大小
	SizeStr        string //文件大小单位
	Type           int    //文件类型
	Postfix        string //文件后缀
}

type Files struct {
	Fs []File
}

func (f *File) MarshalBinary() (data []byte, err error) {
	return json.Marshal(f)
}

func (f *File) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, f)
}

func (fs *Files) MarshalBinary() (data []byte, err error) {
	return json.Marshal(fs)
}

func (fs *Files) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, fs)
}

func CreateFile(ctx context.Context, token, filename, fileHash string, fileSize int64, parentFolderId int, fileStoreId string) (Files, error) {
	//更正parentFolderId
	if parentFolderId == 0 {
		val, _ := redis.Rdb.Get(ctx, "FileFolder_"+token).Result()
		fileFolders := new(FileFolders)
		fileFolders.UnmarshalBinary([]byte(val))
		for _, ff := range fileFolders.Folders {
			if ff.ParentFolderId == -1 {
				parentFolderId = ff.Id
			}
		}
	}
	val, err1 := redis.Rdb.Get(ctx, "File_"+token).Result()
	files := new(Files)
	files.UnmarshalBinary([]byte(val))
	var sizeStr string
	fileSuffix := path.Ext(filename)
	filePrefix := filename[0 : len(filename)-len(fileSuffix)]
	if fileSize < 1048576 {
		sizeStr = strconv.FormatInt(fileSize/1024, 10) + "KB"
	} else {
		sizeStr = strconv.FormatInt(fileSize/1024000, 10) + "MB"
	}
	newFile := File{
		FileName:       filePrefix,
		FileHash:       fileHash,
		FileStoreId:    fileStoreId,
		FilePath:       "",
		DownloadNum:    0,
		UploadTime:     util.GetNowTime(),
		ParentFolderId: parentFolderId,
		Size:           fileSize / 1024,
		SizeStr:        sizeStr,
		Type:           util.GetFileTypeInt(fileSuffix),
		Postfix:        strings.ToLower(fileSuffix),
	}
	result := mysql.Db.Table("go_netdisk_file").Create(&newFile)
	if result.Error == nil && err1 == nil {
		return GetUserFiles(ctx, token, fileStoreId, true)
	}
	return *files, errors.New("error")
}

func GetUserFileCount(ctx context.Context, token string, fileStoreId string) (int64, error) {
	val, err := redis.Rdb.Get(ctx, "FileCount_"+token).Result()
	if err == nil {
		count, e := strconv.ParseInt(val, 10, 64)
		if e == nil {
			return count, err
		}
	}
	var file Files
	var fileCount int64
	result := mysql.Db.Table("go_netdisk_file").Find(&file.Fs, "file_store_id = ?", fileStoreId).Count(&fileCount)
	if result.Error == nil {
		redis.Rdb.Set(ctx, "File_"+token, &file, time.Duration(conf.RedisConfig.TTL)*time.Second).Err()
		redis.Rdb.Set(ctx, "FileCount_"+token, &fileCount, time.Duration(conf.RedisConfig.TTL)*time.Second).Err()
		return fileCount, nil
	}
	return fileCount, result.Error
}

func GetUserFiles(ctx context.Context, token string, fileStoreId string, skipRedis bool) (Files, error) {
	if !skipRedis {
		val, err := redis.Rdb.Get(ctx, "File_"+token).Result()
		if err == nil {
			files := new(Files)
			files.UnmarshalBinary([]byte(val))
			return *files, nil
		}
	}
	var files Files
	result := mysql.Db.Table("go_netdisk_file").Find(&files.Fs, "file_store_id = ?", fileStoreId)
	if result.Error == nil {
		redis.Rdb.Set(ctx, "File_"+token, &files, time.Duration(conf.RedisConfig.TTL)*time.Second).Err()
		return files, nil
	}
	return files, result.Error
}