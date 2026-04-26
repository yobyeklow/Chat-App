package services

import (
	"web_socket/internal/database/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserService interface {
	FindUserByEmail(ctx *gin.Context, userEmail string) (sqlc.User, error)
	FindUserByUUID(ctx *gin.Context, userUUID uuid.UUID) (sqlc.User, error)
	SoftDeleteUser(ctx *gin.Context, userUuid uuid.UUID) (sqlc.User, error)
	HardDeleteUser(ctx *gin.Context, userUuid uuid.UUID) error
	RestoreUser(ctx *gin.Context, userUuid uuid.UUID) (sqlc.User, error)
}
type ChatService interface {
	SendMessage(ctx *gin.Context, roomID, userID, message string) error
	GetMessages(ctx *gin.Context, roomID string) ([]string, error)
}
type AuthService interface {
	Register(ctx *gin.Context, userInput sqlc.CreateUserParams) (sqlc.User, error)
	Login(ctx *gin.Context, email string, password string) (string, string, int, error)
	Logout(ctx *gin.Context, refreshTokenStr string) error
	RefreshToken(ctx *gin.Context, refreshTokenStr string) (string, string, int, error)
}
