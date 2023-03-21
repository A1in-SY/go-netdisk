package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Help(c *gin.Context) {
	c.HTML(http.StatusOK, "help.html", nil)
}