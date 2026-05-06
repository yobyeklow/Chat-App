package routes

import (
	handlers "web_socket/internal/api/handlers"
	"web_socket/internal/common/middleware"

	"github.com/gin-gonic/gin"
)

type ChatRoutes struct {
	handler *handlers.ChatHandler
}

func NewChatRoutes(handler *handlers.ChatHandler) *ChatRoutes {
	return &ChatRoutes{
		handler: handler,
	}
}
func (chatRoute *ChatRoutes) Register(r *gin.RouterGroup) {
	chat := r.Group("/conversations")
	chat.Use(middleware.AuthMiddleware())
	{
		chat.GET("/", chatRoute.handler.GetConversations)
		chat.POST("/:conversation_id/messages", chatRoute.handler.SendMessage)
		chat.GET("/:conversation_id/messages", chatRoute.handler.GetConversations)
	}
}
