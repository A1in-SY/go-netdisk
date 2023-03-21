package controller

import (
	"context"
	"fmt"
	"go-netdisk/model"
	"go-netdisk/redis"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func LoginHandler(c *gin.Context) {
	if c.Request.Method == "POST" {
		name := c.PostForm("name")
		pwd := c.PostForm("pwd")
		ctx := context.Background()
		if redis.GetUser(ctx, name, pwd) {
			c.Redirect(http.StatusFound, "/netdisk/index")
		} else if model.GetUser(name, pwd) {
			redis.SetUser(ctx, name, pwd)
			c.Redirect(http.StatusFound, "/netdisk/index")
		} else {
			c.Redirect(http.StatusFound, "/error")
		}
	} else if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "login.html", nil)
	} else {
		c.Redirect(http.StatusFound, "/error")
	}
}

func RegisterHandler(c *gin.Context) {
	name := c.PostForm("name")
	pw := c.PostForm("pw")
	fmt.Printf("name: %v\n", name)
	fmt.Printf("pw: %v\n", pw)
	c.HTML(http.StatusOK, "index.html", nil)
}
