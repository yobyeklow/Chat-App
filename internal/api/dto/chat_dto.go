package dto

import "encoding/json"

type SendMessageRequest struct {
	MessageContent string          `json:"message_content" binding:"required"`
	Attachments    json.RawMessage `json:"attachments"`
	MessageType    int32           `json:"message_type" binding:"omitempty,min=1,max=3"`
	ReplyTo        *int32          `json:"reply_to"`
}
type SendMessageURI struct {
	ConversationID int32 `uri:"conversation_id" binding:"required"`
}
