package route

import (
	"go-netdisk/controller"
	"go-netdisk/util"
	"log"

	"github.com/gin-gonic/gin"
)

//初始化路由
func SetupRoute() {
	r := gin.Default()
	r.LoadHTMLGlob("view/*")
	r.Static("/static", "./static")

	r.Any("/login", controller.LoginHandler)
	r.Any("/register", controller.RegisterHandler)
	r.GET("/error", controller.Error)

	netdisk := r.Group("/netdisk")
	netdisk.Use(util.CheckLogin())
	{
		netdisk.GET("/index", controller.Index)
	}

	if err := r.Run(":80"); err != nil {
		log.Fatal("服务器启动失败...")
	}
}