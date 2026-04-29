package app

import (
	"web_socket/internal/api/handlers"
	"web_socket/internal/api/repository"
	"web_socket/internal/api/routes"
	"web_socket/internal/api/services"
	"web_socket/pkg/auth"
	"web_socket/pkg/cache"
)

type AuthModule struct {
	routes routes.Routes
}

func NewAuthModule(ctx *ModuleContext, tokenService auth.TokenService, cache cache.RedisService) *AuthModule {
	userRepo := repository.NewSqlUserRepository(ctx.db)
	authService := services.NewAuthServices(userRepo, tokenService, cache)
	authHandler := handlers.NewAuthHandler(authService)
	authRoutes := routes.NewAuthRoutes(authHandler)
	return &AuthModule{
		routes: authRoutes,
	}
}
func (module *AuthModule) Routes() routes.Routes {
	return module.routes
}
