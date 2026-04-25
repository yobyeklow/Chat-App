package ws

import (
	"context"
	"fmt"
	"sync"
)

type Hub struct {
	clients      map[string]*Client
	rooms        map[string]*Room
	mu           *sync.RWMutex
	registerCH   chan *Client
	unregisterCH chan *Client
	broadcastCH  chan *ReqMsg
	joinRoomCH   chan *ReqMsg
	leaveRoomCH  chan *ReqMsg
	roomMsgCH    chan *ReqMsg
}

func NewHub() *Hub {
	return &Hub{
		clients:      make(map[string]*Client),
		rooms:        make(map[string]*Room),
		mu:           new(sync.RWMutex),
		registerCH:   make(chan *Client, 64),
		unregisterCH: make(chan *Client, 64),
		broadcastCH:  make(chan *ReqMsg, 64),
		joinRoomCH:   make(chan *ReqMsg, 64),
		leaveRoomCH:  make(chan *ReqMsg, 64),
		roomMsgCH:    make(chan *ReqMsg, 64),
	}
}

func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client.ID] = client
}

func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, client.ID)
	for _, room := range h.rooms {
		_, ok := room.clients[client.ID]
		if ok {
			delete(room.clients, client.ID)
		}
	}
	fmt.Printf("Client left the server cID=%s\n", client.ID)
}
func (h *Hub) JoinRoom(msg *ReqMsg) {
	roomID := msg.RoomID
	room, ok := h.rooms[roomID]
	if !ok {
		room = NewRoom(roomID)
		h.rooms[roomID] = room
	}

	room.clients[msg.Client.ID] = msg.Client
	fmt.Printf("Client joined the room %s, cID = %s \n", roomID, msg.Client.ID)
}
func (h *Hub) LeaveRoom(msg *ReqMsg) {
	roomID := msg.RoomID
	room, ok := h.rooms[roomID]
	if !ok {
		fmt.Printf("Cannot leave room that does not exist rID = %s, cID = %s\n", roomID, msg.Client.ID)
	}

	delete(room.clients, msg.Client.ID)
	fmt.Printf("Client left the room %s, cID = %s \n", roomID, msg.Client.ID)
}
func (h *Hub) Broadcast(msg *ReqMsg) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	cls := map[string]*Client{}
	for id, client := range h.clients {
		if id != msg.Client.ID {
			cls[id] = client
		}
	}
	go h.SendMessage(msg, cls)
	fmt.Println("Broadcast was sent")
}
func (h *Hub) RoomMsg(msg *ReqMsg) {
	cls := map[string]*Client{}
	roomId := msg.RoomID
	room, ok := h.rooms[roomId]
	if !ok {
		fmt.Printf("The room does not exists, cannot send msg into it\n")
		return
	}

	_, ok = room.clients[msg.Client.ID]
	if !ok {
		fmt.Printf("The clientID = %s does not belong to the room = %s, cannot send msg into it\n", msg.Client.ID, room.ID)
		return
	}
	cls = map[string]*Client{}
	for id, client := range room.clients {
		if id != msg.Client.ID {
			cls[id] = client
		}
	}
	go h.SendMessage(msg, cls)
	fmt.Printf("The clientID = %s sent msg to the room = %s\n", msg.Client.ID, room.ID)
}
func (h *Hub) SendMessage(msg *ReqMsg, cls map[string]*Client) {
	res := NewRespMsg(msg)
	for _, client := range cls {
		client.msgCH <- res
	}
	cls = nil
}

func (h *Hub) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case client := <-h.registerCH:
			h.Register(client)
		case client := <-h.unregisterCH:
			h.Unregister(client)
		case msg := <-h.joinRoomCH:
			h.JoinRoom(msg)
		case msg := <-h.leaveRoomCH:
			go h.LeaveRoom(msg)
		case msg := <-h.roomMsgCH:
			go h.RoomMsg(msg)
		case msg := <-h.broadcastCH:
			h.Broadcast(msg)
		}
	}
}
