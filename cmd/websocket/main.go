package main

import (
	"context"
	"net/http"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
	"web_socket/internal/common/utils"
	"web_socket/internal/ws"
	"web_socket/internal/ws/handlers"
	"web_socket/internal/ws/routes"
	"web_socket/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	rootDir := utils.MustGetWorkingDir()

	logFile := filepath.Join(rootDir, "internal/logs/app.log")

	logger.InitLogger(logger.LoggerConfig{
		Level:      "info",
		Filename:   logFile,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     5,
		Compress:   true,
		IsDev:      utils.GetEnv("APP_EVN", "development"),
	})

	if err := godotenv.Load(filepath.Join(rootDir, ".env")); err != nil {
		logger.Log.Warn().Msg("⚠️ No .env file found")
	} else {
		logger.Log.Info().Msg("✅ Loaded successfully .env in websocket process")
	}

	hub := ws.NewHub()
	if hub == nil {
		logger.Log.Fatal().Msg("Failed to create WebSocket hub")
	}

	wsHandler := handlers.NewWebSocketHandler(hub)
	wsRoutes := routes.NewWSRoutes(wsHandler)

	r := gin.Default()
	wsRoutes.Register(r.Group("/api/v1"))

	wsAddr := utils.GetEnv("WS_SERVER_HOST", "localhost") + ":" + utils.GetEnv("WS_SERVER_PORT", "8081")
	server := &http.Server{
		Addr:    wsAddr,
		Handler: r,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := hub.Run(ctx); err != nil && err != context.Canceled {
			logger.Log.Error().Err(err).Msg("WebSocket hub failed")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Log.Info().Msgf("WebSocket server starting on %s", wsAddr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Log.Error().Err(err).Msg("WebSocket server failed")
		}
	}()

	<-ctx.Done()
	logger.Log.Info().Msg("Received shutdown signal")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error().Err(err).Msg("Server shutdown failed")
	}

	wg.Wait()
	logger.Log.Info().Msg("Main process terminated")
}
