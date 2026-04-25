package routes

import (
	"net/http"
	"web_socket/internal/handlers"

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
	//Public Route
	users.POST("/create", userRoute.handler.CreateUser)

	// Serve static HTML for WebSocket testing
	r.GET("/test-ws", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})
}
