package utils

import "github.com/gin-gonic/gin"

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func Success(c *gin.Context, status int, message string, data any) {
	c.JSON(status, Response{
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, Response{
		Message: message,
	})
}
