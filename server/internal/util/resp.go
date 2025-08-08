package util

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type ErrResp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func JSON(c *gin.Context, status int, v any) {
	c.JSON(status, v)
}

func Error(c *gin.Context, status int, code, msg string) {
	c.AbortWithStatusJSON(status, ErrResp{Code: code, Message: msg})
	c.Error(errors.New(msg))
}
