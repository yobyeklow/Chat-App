package ws

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID    string
	mu    *sync.RWMutex
	Conn  *websocket.Conn
	msgCH chan *ResMsg
	done  chan struct{}
}

func NewClient(conn *websocket.Conn, ID string) *Client {
	return &Client{
		ID:    ID,
		msgCH: make(chan *ResMsg, 64),
		mu:    new(sync.RWMutex),
		Conn:  conn,
		done:  make(chan struct{}),
	}
}

func (c *Client) ReadPump(hub *Hub) {
	defer func() {
		close(c.done)
		hub.unregisterCH <- c
	}()

	for {
		_, payload, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				return
			}
			return
		}

		msg := new(ReqMsg)
		if err := json.Unmarshal(payload, msg); err != nil {
			fmt.Printf("Unable to unmarshal the msg %v \n", err)
			continue
		}

		msg.Client = c
		switch msg.MsgType {
		case MsgType_JoinRoom:
			hub.joinRoomCH <- msg
		case MsgType_LeaveRoom:
			hub.leaveRoomCH <- msg
		case MsgType_RoomMsg:
			hub.roomMsgCH <- msg
		default:
			fmt.Println("Unkown msg type -> ignoring it!")
		}
	}
}

func (c *Client) WritePump(hub *Hub) {
	defer c.Conn.Close()
	for {
		select {
		case <-c.done:
			return
		case msg := <-c.msgCH:
			err := c.Conn.WriteJSON(msg)
			if err != nil {
				fmt.Printf("Error sending msg to clientID =%s\n", c.ID)
				return
			}
		}
	}
}
