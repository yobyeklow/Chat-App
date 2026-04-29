package repository

import (
	"context"
)

type ChatRepository interface {
	SaveMessage(ctx context.Context, roomID, userID, message string) error
	GetMessages(ctx context.Context, roomID string) ([]string, error)
}
