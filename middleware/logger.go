package middleware

import (
	"expense-tracker/middleware/logger"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		now := time.Now()

		params := map[string]interface{}{
			"status_code": ctx.Writer.Status(),
			"method":      ctx.Request.Method,
			"latency":     time.Since(now),
			"path":        ctx.Request.URL,
		}

		if len(ctx.Errors) > 0 && ctx.Errors[0].Err != nil {
			logger.Log.WithFields(params).Error(ctx.Errors[0].Err.Error())
		} else if ctx.Writer.Status() == 404 {
			logger.Log.WithFields(params).Warn("page not found")
		} else if ctx.Writer.Status() >= 400 {
			logger.Log.WithFields(params).Error("request failed")
		} else {
			logger.Log.WithFields(params).Info()
		}

		fmt.Println("=========================================")
	}
}
