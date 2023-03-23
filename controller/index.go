package controller

import (
	"context"
	"go-netdisk/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	ctx := context.Background()
	token, err1 := c.Cookie("token")
	user, err2 := model.GetUser(ctx, token, "", "")
	userFileStore, _ := model.GetUserFileStore(ctx, token, user.Id)
	fileCount, _ := model.GetUserFileCount(ctx, token, user.FileStoreId)
	fileFolderCount, _ := model.GetUserFileFolderCount(ctx, token, user.FileStoreId)
	if err1 == nil && err2 == nil {
		c.HTML(http.StatusOK, "index.html", gin.H {
			"user": user,
			"currIndex": "active",
			"userFileStore": userFileStore,
			"fileCount": fileCount,
			"fileFolderCount": fileFolderCount,
			// "fileDetailUse": fileDetailUse,
		})
	} else {
		c.Redirect(http.StatusMovedPermanently, "/error")
	}
}