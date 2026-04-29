package middleware

import (
	"context"
	"web_socket/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TraceMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceId := ctx.GetHeader("X-Trace-Id")
		if traceId == "" {
			traceId = uuid.New().String()
		}
		context := context.WithValue(ctx.Request.Context(), logger.TraceIdKey, traceId)
		ctx.Request = ctx.Request.WithContext(context)
		ctx.Writer.Header().Set("X-Trace-Id", traceId)
		ctx.Set(string(logger.TraceIdKey), traceId)
		ctx.Next()
	}
}
