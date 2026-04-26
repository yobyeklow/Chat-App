package routes

import (
	"net/http"
	"web_socket/internal/handlers"
	"web_socket/internal/middleware"

	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	handler *handlers.UserHandler
}

func NewUserRoutes(handler *handlers.UserHandler) *UserRoutes {
	return &UserRoutes{
		handler: handler,
	}
}
func (userRoute *UserRoutes) Register(r *gin.RouterGroup) {
	users := r.Group("/users")
	users.Use(middleware.AuthMiddleware())
	{
		users.GET("/:uuid", userRoute.handler.FindUserByUUID)
		users.PUT("/:uuid/restore", userRoute.handler.RestoreUser)
		users.DELETE("/:uuid", userRoute.handler.SoftDeleteUser)
		users.DELETE("/:uuid/clean", userRoute.handler.HardDeleteUser)
	}

	// Serve static HTML for WebSocket testing
	r.GET("/test-ws", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})
}
