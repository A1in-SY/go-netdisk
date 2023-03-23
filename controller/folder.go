package controller

import (
	"context"
	"go-netdisk/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddFileFolder(c *gin.Context) {
	fileFolderName := c.PostForm("fileFolderName")
	parentFolderId := c.PostForm("parentFolderId")
	token, _ := c.Cookie("token")
	ctx := context.Background()
	id, _ := strconv.Atoi(parentFolderId)
	model.CreateFileFolder(ctx, token, fileFolderName, id)
}
