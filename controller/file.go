package controller

import (
	"context"
	"go-netdisk/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Files(c *gin.Context) {
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
	if err1 == nil && err2 == nil {
		c.HTML(http.StatusOK, "files.html", gin.H{
			"currAll":          "active",
			"user":             user,
			"fId":              currentFolder.Id,
			"fName":            currentFolder.FileFolderName,
			// "files":            files,
			"fileFolder":       subFileFolder.Folders,
			"parentFolder":     parentFolder,
			"currentAllParent": currentAllParent,
			// "fileDetailUse":    fileDetailUse,
		})
	} else {
		c.Redirect(http.StatusMovedPermanently, "/error")
	}
}
