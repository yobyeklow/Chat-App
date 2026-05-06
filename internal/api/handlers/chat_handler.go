package handlers

import (
	"net/http"
	"web_socket/internal/api/dto"
	"web_socket/internal/api/services"
	"web_socket/internal/common/utils"
	"web_socket/internal/common/validation"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChatHandler struct {
	service services.ChatService
}

func NewChatHandler(service services.ChatService) *ChatHandler {
	return &ChatHandler{
		service: service,
	}
}

func (ch *ChatHandler) SendMessage(ctx *gin.Context) {
	var inputJSON dto.SendMessageRequest
	if err := ctx.ShouldBindJSON(&inputJSON); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	var inputURI dto.SendMessageURI
	if err := ctx.ShouldBindUri(&inputURI); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	userUUIDData := ctx.GetString("user_uuid")
	userUUID, err := uuid.Parse(userUUIDData)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid User UUID", utils.ErrCodeBadRequest))
		return
	}

	messageData, err := ch.service.SendMessage(ctx, userUUID, inputURI.ConversationID, inputJSON.MessageContent, inputJSON.Attachments, inputJSON.MessageType, inputJSON.ReplyTo)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, "Sent message successfully", messageData)
}
func (ch *ChatHandler) GetConversations(ctx *gin.Context) {
	curUserUUIDData := ctx.GetString("user_uuid")
	curUserUUID, err := uuid.Parse(curUserUUIDData)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid User UUID", utils.ErrCodeBadRequest))
		return
	}
	convs, err := ch.service.GetConversations(ctx, curUserUUID)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, "Fetched conversations successfully", convs)
}
