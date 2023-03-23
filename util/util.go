package util

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

// 生成一个UUID
func GetUUID() string {
	return uuid.New().String()
}

// MD5加密
func MD5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

// 获取当前时间
func GetNowTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// 获取文件类型
func GetFileTypeInt(filePrefix string) int {
	filePrefix = strings.ToLower(filePrefix)
	if filePrefix == ".doc" || filePrefix == ".docx" || filePrefix == ".txt" || filePrefix == ".pdf" {
		return 1
	}
	if filePrefix == ".jpg" || filePrefix == ".png" || filePrefix == ".gif" || filePrefix == ".jpeg" {
		return 2
	}
	if filePrefix == ".mp4" || filePrefix == ".avi" || filePrefix == ".mov" || filePrefix == ".rmvb" || filePrefix == ".rm" {
		return 3
	}
	if filePrefix == ".mp3" || filePrefix == ".cda" || filePrefix == ".wav" || filePrefix == ".wma" || filePrefix == ".ogg" {
		return 4
	}
	return 5
}

// SHA256生成哈希值
func GetSHA256HashCode(file *os.File) string {
	hash := sha256.New()
	_, _ = io.Copy(hash, file)
	bytes := hash.Sum(nil)
	hashCode := hex.EncodeToString(bytes)
	return hashCode
}
