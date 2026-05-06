package repository

import (
	"context"
	"web_socket/internal/common/database/sqlc"
)

type SQLChatRepository struct {
	db sqlc.Querier
}

func NewSQLChatRepository(db sqlc.Querier) ChatRepository {
	return &SQLChatRepository{
		db: db,
	}
}
func (dr *SQLChatRepository) CreateDMConversation(ctx context.Context) (int32, error) {
	conversation, err := dr.db.CreateDMConversation(ctx)
	if err != nil {
		return 0, err
	}
	return conversation, nil
}
func (dr *SQLChatRepository) FindDMConversationn(ctx context.Context, arg sqlc.FindDMConversationParams) (int32, error) {
	conversation, err := dr.db.FindDMConversation(ctx, arg)
	if err != nil {
		return 0, err
	}
	return conversation, nil
}
func (dr *SQLChatRepository) AddParticipantToConversation(ctx context.Context, arg sqlc.AddParticipantToConversationParams) error {
	err := dr.db.AddParticipantToConversation(ctx, arg)
	if err != nil {
		return err
	}
	return nil
}
func (dr *SQLChatRepository) SendMessage(ctx context.Context, arg sqlc.SendMessageParams) (int32, error) {
	message, err := dr.db.SendMessage(ctx, arg)
	if err != nil {
		return 0, err
	}
	return message, nil
}
func (dr *SQLChatRepository) GetMessage(ctx context.Context, arg sqlc.GetMessagesParams) ([]sqlc.GetMessagesRow, error) {
	messages, err := dr.db.GetMessages(ctx, arg)
	if err != nil {
		return []sqlc.GetMessagesRow{}, err
	}
	return messages, nil
}
func (mr *SQLChatRepository) ValidateReplyMessage(ctx context.Context, arg sqlc.ValidateReplyParams) (bool, error) {
	data, err := mr.db.ValidateReply(ctx, arg)
	if err != nil {
		return data, err
	}
	return data, nil
}
func (dr *SQLChatRepository) GetDMConversations(ctx context.Context, curUserID int32) ([]sqlc.GetDMConversationRow, error) {
	data, err := dr.db.GetDMConversation(ctx, curUserID)
	if err != nil {
		return data, err
	}
	return data, nil
}
func (dr *SQLChatRepository) GetGroupConversations(ctx context.Context, curUserID int32) ([]sqlc.GetGroupConversationRow, error) {
	data, err := dr.db.GetGroupConversation(ctx, curUserID)
	if err != nil {
		return data, err
	}
	return data, nil
}
