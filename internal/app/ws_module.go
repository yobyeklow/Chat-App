package app

import (
	"web_socket/internal/handlers"
	"web_socket/internal/routes"
	"web_socket/internal/ws"
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
