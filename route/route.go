package route

import (
	"go-netdisk/controller"
	"go-netdisk/middleware"
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
	netdisk.Use(middleware.CheckLogin())
	{
		netdisk.GET("/index", controller.Index)
		netdisk.GET("/logout", controller.Logout)
		netdisk.GET("/help", controller.Help)
		netdisk.GET("/files", controller.Files)
		netdisk.GET("/upload", controller.Upload)
	}
	{
		netdisk.POST("/addFolder", controller.AddFileFolder)
		netdisk.POST("/uploadFile", controller.UploadFile)
	}

	if err := r.Run(":80"); err != nil {
		log.Fatal("服务器启动失败...")
	}
}