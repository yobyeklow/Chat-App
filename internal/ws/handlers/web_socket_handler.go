package handlers

import (
	"fmt"
	"net/http"
	"web_socket/internal/common/utils"
	"web_socket/internal/common/validation"
	"web_socket/internal/ws"
	"web_socket/internal/ws/dto"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	hub *ws.Hub
}

func NewWebSocketHandler(hub *ws.Hub) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
	}
}

func (wh *WebSocketHandler) HandleWebSocket(ctx *gin.Context) {
	var input dto.WSInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	upgrader := websocket.Upgrader{
		ReadBufferSize:  512,
		WriteBufferSize: 512,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Printf("Upgrade socket err:%v\n", err)
		return
	}

	client := ws.NewClient(conn, input.ID)
	wh.hub.Register(client)

	go client.WritePump(wh.hub)
	go client.ReadPump(wh.hub)
}
