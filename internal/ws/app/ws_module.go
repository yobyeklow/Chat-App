package app

import (
	"web_socket/internal/ws"
	"web_socket/internal/ws/handlers"
	"web_socket/internal/ws/routes"
)

type WSModule struct {
	route routes.Routes
}

func NewWSModule(hub *ws.Hub) *WSModule {
	wsHandler := handlers.NewWebSocketHandler(hub)
	wsRoutes := routes.NewWSRoutes(wsHandler)
	return &WSModule{
		route: wsRoutes,
	}
}

func (module *WSModule) Routes() routes.Routes {
	return module.route
}
