package util

import (
	"github.com/gin-gonic/gin"
)

//检查是否登录中间件
func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.Redirect(http.StatusFound, "/login")
		c.Next()
	}
}