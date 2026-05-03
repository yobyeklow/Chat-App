package middleware

import (
	"net/http"
	"strings"
	"web_socket/pkg/auth"
	"web_socket/pkg/cache"

	"github.com/gin-gonic/gin"
)

var (
	jwtService   auth.TokenService
	cacheService cache.RedisService
)

func InitAuthMiddlware(jwt auth.TokenService, cache cache.RedisService) {
	jwtService = jwt
	cacheService = cache
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"Error": "Authorization header missing or invalid",
			})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		_, claims, err := jwtService.ParseToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"Error": "Authorization header missing or invalid!",
			})
			return
		}
		if jti, ok := claims["jti"].(string); ok {
			redis_key := "blacklist:" + jti
			exists, err := cacheService.Exists(redis_key)
			if err == nil && exists {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"Error": "Token revoked!",
				})
				return
			}
		}
		payload, err := jwtService.DecryptAccessTokenPayload(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"Error": "Authorization header missing or invalid!",
			})
			return
		}
		ctx.Set("user_uuid", payload.UserUUID)
		ctx.Set("user_email", payload.Email)
		ctx.Set("user_role", payload.Role)
		ctx.Next()
	}
}
