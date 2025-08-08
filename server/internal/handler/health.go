package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Healthz() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}
