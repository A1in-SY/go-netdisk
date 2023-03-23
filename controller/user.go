package controller

import (
	"context"
	"go-netdisk/conf"
	"go-netdisk/model"
	"go-netdisk/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func LoginHandler(c *gin.Context) {
	if c.Request.Method == "POST" {
		name := c.PostForm("name")
		pwd := util.MD5(c.PostForm("pwd"))
		token, _ := c.Cookie("token")
		ctx := context.Background()
		new_token := util.GetUUID()
		if user, err := model.GetUser(ctx, "", name, pwd); err == nil && name == user.UserName && pwd == user.Password {
			c.SetCookie("token", new_token, conf.RedisConfig.TTL, "/", "127.0.0.1", false, true)
			model.SetTokenUser(ctx, new_token, user)
			model.DelTokenUser(ctx, token)
			c.Redirect(http.StatusFound, "/netdisk/index")
		} else {
			c.Redirect(http.StatusMovedPermanently, "/error")
		}
	} else if c.Request.Method == "GET" {
		token, err1 := c.Cookie("token")
		ctx := context.Background()
		new_token := util.GetUUID()
		if err1 == nil {
			user, err2 := model.GetUser(ctx, token, "", "")
			if err2 == nil {
				c.SetCookie("token", new_token, conf.RedisConfig.TTL, "/", "127.0.0.1", false, true)
				model.SetTokenUser(ctx, new_token, user)
				model.DelTokenUser(ctx, token)
				c.Redirect(http.StatusFound, "/netdisk/index")
				return
			}
		}
		c.SetCookie("token", "", -1, "/", "127.0.0.1", false, true)
		c.HTML(http.StatusOK, "login.html", nil)
	} else {
		c.Redirect(http.StatusMovedPermanently, "/error")
	}
}

func RegisterHandler(c *gin.Context) {
	if c.Request.Method == "POST" {
		name := c.PostForm("name")
		pwd := util.MD5(c.PostForm("pwd"))
		ctx := context.Background()
		token := util.GetUUID()
		if !model.IfExistName(name) {
			if user, err := model.CreateUser(name, pwd, token, ctx); err == nil {
				RegisterSuccess(user)
				c.SetCookie("token", token, conf.RedisConfig.TTL, "/", "127.0.0.1", false, true)
				c.Redirect(http.StatusFound, "/login")
				return
			}
		}
		c.Redirect(http.StatusMovedPermanently, "/error")
	} else {
		c.Redirect(http.StatusMovedPermanently, "/error")
	}
}

func RegisterSuccess(user model.User) {
	model.CreateUserFileStore(user.Id, user.FileStoreId)
	model.CreateUserRootFileFolder(user.Id, user.FileStoreId)
}

func Logout(c *gin.Context) {
	token, _ := c.Cookie("token")
	c.SetCookie("token", "", -1, "/", "127.0.0.1", false, true)
	ctx := context.Background()
	model.DelTokenUser(ctx, token)
	c.Redirect(http.StatusFound, "/login")
}
