package routes

import (
	"net/http"
	"web_socket/internal/handlers"

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

	r.GET("/test-ws", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})
}
