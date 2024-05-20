package web

import (
	"github.com/gin-gonic/gin"
)

// ? ==================== Structs ====================

type ErrorResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ? ==================== Functions ====================

// SuccessResponseBody returns a response with a status code passed by parameter
func SuccessResponseBody(ctx *gin.Context, status int, data interface{}) {
	ctx.JSON(status, data)
}

// * =========== *

// ErrorResponseBody returns a response with a status code passed by parameter
func ErrorResponseBody(ctx *gin.Context, status int, code string, message string) {
	ctx.JSON(status, ErrorResponse{
		Status:  status,
		Code:    code,
		Message: message,
	})
}
