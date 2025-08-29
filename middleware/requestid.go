package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const headerXRequestID = "X-Request-ID"

const requestIDKey = "request_id"

func RequestIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rid := ctx.GetHeader(headerXRequestID)
		if rid == "" {
			rid = uuid.New().String()
		}
		ctx.Set(requestIDKey, rid)
		ctx.Header(headerXRequestID, rid)
		ctx.Next()
	}
}

func GetRequestID(ctx context.Context) string {
	if requestId, ok := ctx.Value(requestIDKey).(string); ok {
		return requestId
	}
	return ""
}
