package ws

type MsgType string

const (
	MsgType_JoinRoom  MsgType = "join_room"
	MsgType_LeaveRoom MsgType = "leave_room"
	MsgType_BroadCast MsgType = "broadcast"
	MsgType_RoomMsg   MsgType = "room-msg"
)

type ReqMsg struct {
	MsgType MsgType `json:"type"`
	Client  *Client
	Data    interface{} `json:"data"`
	RoomID  string      `json:"roomID"`
}

type ResMsg struct {
	MsgType  MsgType     `json:"type"`
	Data     interface{} `json:"data"`
	SenderID string      `json:"senderID"`
	RoomID   string      `json:"roomID"`
}

func NewRespMsg(msg *ReqMsg) *ResMsg {
	return &ResMsg{
		MsgType:  msg.MsgType,
		Data:     msg.Data,
		SenderID: msg.Client.ID,
		RoomID:   msg.RoomID,
	}
}
