package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func sendError(c *gin.Context, code int, msg string) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(code, gin.H{
		"message":   msg,
		"errorCode": code,
	})
}

func sendSuccess(c *gin.Context, op string, data interface{}) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(200, gin.H{
		"data":    data,
		"message": fmt.Sprintf("%s", op),
	})
}
