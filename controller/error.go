package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context) {
	c.HTML(http.StatusOK, "error.html", nil)
}