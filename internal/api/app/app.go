package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web_socket/internal/api/routes"
	"web_socket/internal/common/config"
	"web_socket/internal/common/database"
	"web_socket/internal/common/validation"
	"web_socket/internal/ws"
	"web_socket/pkg/auth"
	"web_socket/pkg/cache"
	"web_socket/pkg/logger"

	"web_socket/internal/common/database/sqlc"

	"github.com/gin-gonic/gin"
)

type Module interface {
	Routes() routes.Routes
}
type Application struct {
	cfg     *config.Config
	router  *gin.Engine
	modules []Module
}
type ModuleContext struct {
	db  sqlc.Querier
	hub *ws.Hub
}

func NewApplication(cfg *config.Config) (*Application, error) {
	if err := validation.InitValidator(); err != nil {
		logger.Log.Fatal().Msgf("Validator init failed %v", err)
		return nil, err
	}
	r := gin.Default()
	r.LoadHTMLGlob("front-end/*.html")
	r.Static("/static", "./front-end")
	//Connect DB
	if err := database.InitDB(); err != nil {
		logger.Log.Fatal().Msgf("Database init failed %v", err)
		return nil, err
	}
	redisClient := config.NewRedisClient()
	cacheService := cache.NewRedisCacheService(redisClient)
	tokenService := auth.NewJWTService(cacheService)
	ctx := &ModuleContext{
		db:  database.DB,
		hub: ws.NewHub(),
	}
	modules := []Module{
		NewUserModule(ctx),
		NewAuthModule(ctx, tokenService, cacheService),
		NewGroupModule(ctx, tokenService),
	}
	routes.RegisterRoutes(r, tokenService, cacheService, getModulesRoute(modules)...)
	return &Application{
		router:  r,
		cfg:     cfg,
		modules: modules,
	}, nil

}
func getModulesRoute(modules []Module) []routes.Routes {
	routeList := make([]routes.Routes, len(modules))
	for i, module := range modules {
		routeList[i] = module.Routes()
	}
	return routeList
}
func (app *Application) Run() error {
	server := &http.Server{
		Addr:    app.cfg.ServerAddress,
		Handler: app.router,
	}

	quitSrv := make(chan os.Signal, 1)
	signal.Notify(quitSrv, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	logger.Log.Info().Msgf("Server is running at %s", app.cfg.ServerAddress)
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()
	<-quitSrv
	logger.Log.Info().Msg("Shutdown signal received!")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Log.Fatal().Err(err).Msg("Server forced to shutdown")
	}
	logger.Log.Info().Msg("Server shutdown!")
	return nil
}
