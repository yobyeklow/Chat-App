package repository

import (
	"context"
	"web_socket/internal/database/sqlc"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error)
	FindUserByEmail(ctx context.Context, userEmail string) (sqlc.User, error)
	FindUserByUUID(ctx context.Context, userUUID uuid.UUID) (sqlc.User, error)
	SoftDeleteUser(ctx context.Context, userUuid uuid.UUID) (sqlc.User, error)
	HardDeleteUser(ctx context.Context, userUuid uuid.UUID) error
	RestoreUser(ctx context.Context, userUuid uuid.UUID) (sqlc.User, error)
}
type ChatRepository interface {
	SaveMessage(ctx context.Context, roomID, userID, message string) error
	GetMessages(ctx context.Context, roomID string) ([]string, error)
}
