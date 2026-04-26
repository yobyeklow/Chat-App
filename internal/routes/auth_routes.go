package routes

import (
	"web_socket/internal/handlers"

	"github.com/gin-gonic/gin"
)

type AuthRoutes struct {
	handler *handlers.AuthHandler
}

func NewAuthRoutes(handler *handlers.AuthHandler) *AuthRoutes {
	return &AuthRoutes{
		handler: handler,
	}
}
func (authRoute *AuthRoutes) Register(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", authRoute.handler.Register)
		auth.POST("/login", authRoute.handler.Login)
		auth.POST("/logout", authRoute.handler.Logout)
		auth.POST("/refresh-token", authRoute.handler.RefreshToken)
	}
}
