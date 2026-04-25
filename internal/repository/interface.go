package repository

import (
	"context"
	"web_socket/internal/database/sqlc"
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error)
	FindUserByEmail(ctx context.Context, userEmail string) (sqlc.User, error)
}
type ChatRepository interface {
	SaveMessage(ctx context.Context, roomID, userID, message string) error
	GetMessages(ctx context.Context, roomID string) ([]string, error)
}
