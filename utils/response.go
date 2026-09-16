package utils

import "github.com/gin-gonic/gin"

// Response is the uniform JSON envelope used by the helpers below.
type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Success writes a JSON response with a message and optional data payload.
func Success(c *gin.Context, status int, message string, data any) {
	c.JSON(status, Response{
		Message: message,
		Data:    data,
	})
}

// Error writes a JSON error response with a message.
func Error(c *gin.Context, status int, message string) {
	c.JSON(status, Response{
		Message: message,
	})
}
