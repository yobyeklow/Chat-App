package routes

import (
	"web_socket/internal/middleware"
	"web_socket/internal/utils"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type Routes interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, routes ...Routes) {
	recoverLogger := utils.NewLoggerWithPath("recovery.log", "warning")
	rateLimitLogger := utils.NewLoggerWithPath("ratelimit.log", "warning")
	httpLogger := utils.NewLoggerWithPath("http.log", "info")
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(
		middleware.RateLimiterMiddleware(rateLimitLogger),
		middleware.TraceMiddleware(),
		middleware.LoggerMiddleware(httpLogger),
		middleware.RecoveryMiddleware(recoverLogger),
	)
	api := r.Group("/api/v1")

	for _, route := range routes {
		route.Register(api)
	}
	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{
			"Error": "Not Found",
			"path":  ctx.Request.URL.Path,
		})
	})
}
