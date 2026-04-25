package repository

import (
	"context"
	"sync"
)

type SQLChatRepository struct {
	mu       sync.Mutex
	messages map[string][]string
}

func NewSQLChatRepository() ChatRepository {
	return &SQLChatRepository{
		messages: make(map[string][]string),
	}
}

func (r *SQLChatRepository) SaveMessage(ctx context.Context, roomID, userID, message string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.messages[roomID] = append(r.messages[roomID], message)
	return nil
}

func (r *SQLChatRepository) GetMessages(ctx context.Context, roomID string) ([]string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.messages[roomID], nil
}
