package utils

import "github.com/gin-gonic/gin"

func ApiResponse(ctx *gin.Context, statusCode int, data interface{}) {
	ctx.JSON(statusCode, data)
	if statusCode >= 400 {
		ctx.AbortWithStatus(statusCode)
	}
}
