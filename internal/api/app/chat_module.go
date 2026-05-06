package app

import (
	"web_socket/internal/api/handlers"
	"web_socket/internal/api/repository"
	"web_socket/internal/api/routes"
	"web_socket/internal/api/services"
	"web_socket/pkg/auth"
	"web_socket/pkg/cache"
)

type ChatModule struct {
	routes routes.Routes
}

func NewChatModule(ctx *ModuleContext, tokenService auth.TokenService, cache cache.RedisService) *ChatModule {
	groupRepo := repository.NewSQLGroupRepository(ctx.db)
	groupSrv := services.NewGroupService(groupRepo)
	userRepo := repository.NewSqlUserRepository(ctx.db)
	userSrv := services.NewUserService(userRepo)
	chatRepo := repository.NewSQLChatRepository(ctx.db)
	chatService := services.NewMessageServices(chatRepo, groupSrv, userSrv)
	chatHandler := handlers.NewChatHandler(chatService)
	chatRoutes := routes.NewChatRoutes(chatHandler)
	return &ChatModule{
		routes: chatRoutes,
	}
}
func (module *ChatModule) Routes() routes.Routes {
	return module.routes
}
