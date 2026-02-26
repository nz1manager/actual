package utils

import "github.com/gin-gonic/gin"

type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

func SuccessResponse(c *gin.Context, status int, data interface{}) {
    c.JSON(status, Response{
        Success: true,
        Data:    data,
    })
}

func ErrorResponse(c *gin.Context, status int, message string) {
    c.JSON(status, Response{
        Success: false,
        Error:   message,
    })
}

func MessageResponse(c *gin.Context, status int, message string) {
    c.JSON(status, Response{
        Success: true,
        Message: message,
    })
}
