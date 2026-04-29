package app

import (
	"web_socket/internal/api/handlers"
	"web_socket/internal/api/repository"
	"web_socket/internal/api/routes"
	"web_socket/internal/api/services"
)

type UserModule struct {
	route routes.Routes
}

func NewUserModule(ctx *ModuleContext) *UserModule {
	userRepo := repository.NewSqlUserRepository(ctx.db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)
	userRoutes := routes.NewUserRoutes(userHandler)
	return &UserModule{
		route: userRoutes,
	}
}
func (module *UserModule) Routes() routes.Routes {
	return module.route
}
