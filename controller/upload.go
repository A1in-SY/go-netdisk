package controller

import (
	"context"
	"fmt"
	"go-netdisk/conf"
	"go-netdisk/model"
	"go-netdisk/util"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	ctx := context.Background()
	token, err1 := c.Cookie("token")
	user, err2 := model.GetUser(ctx, token, "", "")
	fileFolders, _ := model.GetUserFileFolders(ctx, token, user.FileStoreId, false)

	fId, _ := strconv.Atoi(c.DefaultQuery("fId", "0"))
	currentFolder := fileFolders.SelectFolderById(fId)
	rootFolder := fileFolders.SelectFoldersByParentId(-1).Folders[0]
	var subFileFolder model.FileFolders
	if currentFolder.Id == 0 {
		currentFolder.ParentFolderId = rootFolder.Id
		subFileFolder = fileFolders.SelectFoldersByParentId(currentFolder.ParentFolderId)
	} else {
		subFileFolder = fileFolders.SelectFoldersByParentId(currentFolder.Id)
	}
	parentFolder := fileFolders.SelectFolderById(currentFolder.ParentFolderId)
	currentAllParent := model.GetAllParents(ctx, token, currentFolder.Id)
	fmt.Printf("user: %v\n", user)
	fmt.Printf("currentFolder: %v\n", currentFolder)
	fmt.Printf("subFileFolder: %v\n", subFileFolder)
	fmt.Printf("currentAllParent: %v\n", currentAllParent)
	if err1 == nil && err2 == nil {
		c.HTML(http.StatusOK, "upload.html", gin.H{
			"currAll": "active",
			"user":    user,
			"fId":     currentFolder.Id,
			"fName":   currentFolder.FileFolderName,
			// "files":            files,
			"fileFolders":      subFileFolder.Folders,
			"parentFolder":     parentFolder,
			"currentAllParent": currentAllParent,
			// "fileDetailUse":    fileDetailUse,
		})
	} else {
		c.Redirect(http.StatusMovedPermanently, "/error")
	}
}

func UploadFile(c *gin.Context) {
	ctx := context.Background()
	token, _ := c.Cookie("token")
	user, _ := model.GetUser(ctx, token, "", "")
	fId, _ := strconv.Atoi(c.DefaultQuery("fId", "0"))
	file, head, err := c.Request.FormFile("file")
	defer file.Close()
	if err != nil {
		fmt.Println("文件上传错误", err.Error())
		return
	}
	// if ok := model.CurrFileExists(Fid, head.Filename); !ok {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"code": 501,
	// 	})
	// 	return
	// }
	// if ok := model.CapacityIsEnough(head.Size, user.FileStoreId); !ok {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"code": 503,
	// 	})
	// 	return
	// }
	fileSuffix := path.Ext(head.Filename)
	storedFileName := util.GetUUID()+fileSuffix
	fileLocation := conf.AppConfig.Upload_location + "/" + user.FileStoreId + "/" + storedFileName
	newFile, err := os.Create(fileLocation)
	if err != nil {
		log.Fatalf("文件创建失败")
		return
	}
	defer newFile.Close()
	fileSize, err := io.Copy(newFile, file)
	if err != nil {
		log.Fatalf("文件创建失败文件拷贝错误")
		return
	}
	_, _ = newFile.Seek(0, 0)
	fileHash := util.GetSHA256HashCode(newFile)
	model.CreateFile(ctx, token, head.Filename, fileHash, fileSize, fId, user.FileStoreId)
	model.SubtractSize(ctx, token, fileSize, user.FileStoreId)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
	})
}
