package util

import (
	"expense-tracker/dto"

	"github.com/gin-gonic/gin"
)

func ResponseMsg(ctx *gin.Context, success bool, err, data any, statusCode int) {
	ctx.JSON(statusCode, dto.Response{
		Success: success,
		Error:   err,
		Data:    data,
	})
}
