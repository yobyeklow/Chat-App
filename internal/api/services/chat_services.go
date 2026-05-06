package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"
	"web_socket/internal/api/repository"
	"web_socket/internal/common/database/sqlc"
	"web_socket/internal/common/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type chatService struct {
	repo     repository.ChatRepository
	groupSvc GroupService
	userSvc  UserService
}

func NewMessageServices(repo repository.ChatRepository, groupSvc GroupService, userSvc UserService) ChatService {
	return &chatService{
		repo:     repo,
		groupSvc: groupSvc,
		userSvc:  userSvc,
	}
}
func (dm *chatService) StartDirectMessage(ctx *gin.Context, sender uuid.UUID, receiver uuid.UUID) error {
	context := ctx.Request.Context()
	if sender == receiver {
		return utils.NewError("Cannot start DM with yourself", utils.ErrCodeBadRequest)
	}
	senderData, err := dm.userSvc.FindUserByUUID(ctx, sender)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NewError("Sender not found", utils.ErrCodeBadRequest)
		}
		return utils.WrapError("Failed to fetch Sender data", utils.ErrCodeInternal, err)
	}
	receiverData, err := dm.userSvc.FindUserByUUID(ctx, receiver)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NewError("Receiver not found", utils.ErrCodeBadRequest)
		}
		return utils.WrapError("Failed to fetch Receiver data", utils.ErrCodeInternal, err)
	}
	findDMArg := sqlc.FindDMConversationParams{
		User1ID: senderData.UserID,
		User2ID: receiverData.UserID,
	}
	convID, err := dm.repo.FindDMConversationn(context, findDMArg)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return utils.WrapError("Failed to find dm conversation", utils.ErrCodeInternal, err)

		} else if errors.Is(err, sql.ErrNoRows) {
			convID, err = dm.repo.CreateDMConversation(context)
			if err != nil {
				return utils.WrapError("Failed to create new conversation", utils.ErrCodeInternal, err)
			}
		}
	}

	for _, uid := range []int32{senderData.UserID, receiverData.UserID} {
		addParticipantArg := sqlc.AddParticipantToConversationParams{
			ConversationID: convID,
			UserID:         uid,
		}
		err = dm.repo.AddParticipantToConversation(context, addParticipantArg)
		if err != nil {
			return utils.WrapError("Failed to add participant", utils.ErrCodeInternal, err)
		}
	}
	return nil
}
func (dm *chatService) SendMessage(ctx *gin.Context, senderUUID uuid.UUID, conversationID int32, content string,
	attachments json.RawMessage, msgType int32, replyTo *int32) (int32, error) {
	context := ctx.Request.Context()
	_, err := dm.userSvc.FindUserByUUID(ctx, senderUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, utils.NewError("User not found", utils.ErrCodeBadRequest)
		}
		return 0, utils.WrapError("Failed to fetch User data", utils.ErrCodeInternal, err)
	}
	if replyTo != nil {
		validArg := sqlc.ValidateReplyParams{
			MessageID:      *replyTo,
			ConversationID: conversationID,
		}
		valid, err := dm.repo.ValidateReplyMessage(context, validArg)
		if err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) && pqErr.Code == "23503" {
				return 0, utils.NewError("Replied message was deleted", utils.ErrCodeBadRequest)
			}
			return 0, utils.WrapError("Failed to validate reply message", utils.ErrCodeInternal, err)
		}
		if !valid {
			return 0, utils.NewError("Replied message not found or has been deleted", utils.ErrCodeBadRequest)
		}
	}

	arg := sqlc.SendMessageParams{
		UserUuid:       senderUUID,
		ConversationID: conversationID,
		Content:        content,
		Attachments:    attachments,
		MessageType:    msgType,
		ReplyTo:        *replyTo,
	}
	messageId, err := dm.repo.SendMessage(context, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, utils.NewError("Not a participant or conversation not found", utils.ErrCodeBadRequest)
		}
		return 0, utils.WrapError("Failed to send message", utils.ErrCodeInternal, err)
	}
	return messageId, nil
}
func (dm *chatService) GetMessages(ctx *gin.Context, conversationID int32, cursorTime *time.Time, cursorID *int32, limit int32) ([]sqlc.GetMessagesRow, error) {
	context := ctx.Request.Context()

	arg := sqlc.GetMessagesParams{
		ConversationID: conversationID,
		CursorTime:     utils.ToTimestamptz(cursorTime),
		CursorID:       cursorID,
		Limitarg:       limit,
	}
	messages, err := dm.repo.GetMessage(context, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []sqlc.GetMessagesRow{}, utils.NewError("Not a participant or conversation not found", utils.ErrCodeBadRequest)
		}
		return []sqlc.GetMessagesRow{}, utils.WrapError("Failed to send message", utils.ErrCodeInternal, err)
	}

	return messages, nil
}

type Conversation struct {
	ID            int32  `json:"id"`
	Type          int32  `json:"type"`
	Title         string `json:"title"`
	Avatar        string `json:"avatar" binding:"omitempty"`
	LastMessageAt string `json:"last_message_at"`
	UnreadCount   int32  `json:"unread_count"`
}

func (dm *chatService) GetConversations(ctx *gin.Context, curUserUUID uuid.UUID) ([]Conversation, error) {
	context := ctx.Request.Context()
	user, err := dm.userSvc.FindUserByUUID(ctx, curUserUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Conversation{}, utils.NewError("User not found", utils.ErrCodeBadRequest)
		}
		return []Conversation{}, utils.WrapError("Failed to fetch User data", utils.ErrCodeInternal, err)
	}
	var (
		dmConvs []sqlc.GetDMConversationRow
		grConvs []sqlc.GetGroupConversationRow
	)
	dmConvs, err = dm.repo.GetDMConversations(context, user.UserID)
	if err != nil {
		return []Conversation{}, utils.WrapError("Failed to get direct message conversations", utils.ErrCodeInternal, err)
	}
	grConvs, err = dm.repo.GetGroupConversations(context, user.UserID)
	if err != nil {
		return []Conversation{}, utils.WrapError("Failed to get groups conversations", utils.ErrCodeInternal, err)
	}
	var conversations []Conversation

	for _, c := range dmConvs {
		conversations = append(conversations, Conversation{
			ID:            c.ConversationID,
			Type:          1,
			Title:         c.UserEmail,
			LastMessageAt: c.LastMessageAt.Time.Format("2006-01-02T15:04:05Z07:00"),
			UnreadCount:   int32(c.UnreadCount),
		})
	}

	for _, c := range grConvs {
		conversations = append(conversations, Conversation{
			ID:            c.ConversationID,
			Type:          2,
			Title:         c.Title,
			LastMessageAt: c.LastMessageAt.Time.Format("2006-01-02T15:04:05Z07:00"),
			UnreadCount:   int32(c.UnreadCount),
		})
	}
	return conversations, nil
}
