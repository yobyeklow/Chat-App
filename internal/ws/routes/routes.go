package routes

import (
	"web_socket/internal/common/middleware"
	"web_socket/internal/common/utils"
	"web_socket/pkg/auth"
	"web_socket/pkg/cache"
	"web_socket/pkg/logger"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type Routes interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, authService auth.TokenService, cacheService cache.RedisService, routes ...Routes) {
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
	middleware.InitAuthMiddlware(authService, cacheService)
	logger.Log.Info().Msgf("Registered %d route groups", len(routes))
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
