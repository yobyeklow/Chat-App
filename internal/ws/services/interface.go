package services

import "github.com/gin-gonic/gin"

type ChatService interface {
	SendMessage(ctx *gin.Context, roomID, userID, message string) error
	GetMessages(ctx *gin.Context, roomID string) ([]string, error)
}
