package services

import (
	"web_socket/internal/database/sqlc"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	CreateUser(ctx *gin.Context, userInput sqlc.CreateUserParams) (sqlc.User, error)
	FindUserByEmail(ctx *gin.Context, userEmail string) (sqlc.User, error)
}
type ChatService interface {
	SendMessage(ctx *gin.Context, roomID, userID, message string) error
	GetMessages(ctx *gin.Context, roomID string) ([]string, error)
}
