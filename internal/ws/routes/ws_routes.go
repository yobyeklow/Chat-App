package routes

import (
	"web_socket/internal/ws/handlers"

	"github.com/gin-gonic/gin"
)

type WSRoutes struct {
	handler *handlers.WebSocketHandler
}

func NewWSRoutes(handler *handlers.WebSocketHandler) *WSRoutes {
	return &WSRoutes{
		handler: handler,
	}
}

func (wsRoute *WSRoutes) Register(r *gin.RouterGroup) {
	ws := r.Group("/ws")
	ws.GET("/chat", func(ctx *gin.Context) {
		wsRoute.handler.HandleWebSocket(ctx)
	})
}
