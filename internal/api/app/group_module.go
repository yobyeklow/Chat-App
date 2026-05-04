package app

import (
	"web_socket/internal/api/handlers"
	"web_socket/internal/api/repository"
	"web_socket/internal/api/routes"
	"web_socket/internal/api/services"
	"web_socket/pkg/auth"
)

type GroupModule struct {
	route routes.Routes
}

func NewGroupModule(ctx *ModuleContext, tokenService auth.TokenService) *GroupModule {
	groupRepo := repository.NewSQLGroupRepository(ctx.db)
	groupService := services.NewGroupService(groupRepo)
	groupHandler := handlers.NewGroupHandler(groupService, tokenService)
	groupRoutes := routes.NewGroupRoutes(groupHandler)
	return &GroupModule{
		route: groupRoutes,
	}
}
func (module *GroupModule) Routes() routes.Routes {
	return module.route
}
