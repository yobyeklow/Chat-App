package dto

type WSMessageRequest struct {
	Type   string `json:"type"`
	RoomID string `json:"room_id" binding:"omitempty"`
	Data   string `json:"data"`
}

type WSMessageResponse struct {
	Type   string `json:"type"`
	RoomID string `json:"room_id" binding:"omitempty"`
	UserID string `json:"user_id" binding:"omitempty"`
	Data   string `json:"data"`
}

type WSRoomJoinRequest struct {
	Type   string `json:"type"` // "join_room"
	RoomID string `json:"room_id"`
}

type WSRoomLeaveRequest struct {
	Type   string `json:"type"` // "leave_room"
	RoomID string `json:"room_id"`
}
type WSInput struct {
	ID string `json:"user_id" binding:"required"`
}
