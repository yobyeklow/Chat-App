package ws

type Room struct {
	clients map[string]*Client
	ID      string
}

func NewRoom(id string) *Room {
	return &Room{
		ID:      id,
		clients: map[string]*Client{},
	}
}
