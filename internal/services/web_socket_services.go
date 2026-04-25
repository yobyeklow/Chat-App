package services

import (
	"web_socket/internal/repository"

	"github.com/gin-gonic/gin"
)

type WebSocketService struct {
	repo repository.ChatRepository
}

func NewWebSocketService(repo repository.ChatRepository) ChatService {
	return &WebSocketService{
		repo: repo,
	}
}

func (s *WebSocketService) SendMessage(ctx *gin.Context, roomID, userID, message string) error {
	return s.repo.SaveMessage(ctx.Request.Context(), roomID, userID, message)
}

func (s *WebSocketService) GetMessages(ctx *gin.Context, roomID string) ([]string, error) {
	return s.repo.GetMessages(ctx.Request.Context(), roomID)
}
