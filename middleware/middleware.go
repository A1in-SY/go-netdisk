package middleware

import (
	"context"
	"go-netdisk/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

//检查是否登录中间件
func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		ctx := context.Background()
		if err == nil && model.CheckToken(ctx, token) {
			c.Next()
		}
		c.Redirect(http.StatusFound, "/login")
	}
}