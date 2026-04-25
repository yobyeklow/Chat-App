package main

import (
	"path/filepath"
	"web_socket/internal/app"
	"web_socket/internal/config"
	"web_socket/internal/utils"
	"web_socket/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
	rootDir := utils.MustGetWorkingDir()
	logFile := filepath.Join(rootDir, "internal/logs/app.log")
	logger.InitLogger(logger.LoggerConfig{
		Filename:   logFile,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     5,
		Compress:   true,
		Level:      "info",
		IsDev:      utils.GetEnv("APP_EVN", "development"),
	})
	if err := godotenv.Load(filepath.Join(rootDir, ".env")); err != nil {
		logger.Log.Warn().Msg("⚠️ No .env file found")
	} else {
		logger.Log.Info().Msg("✅ Loaded successfully .env in api proccess")
	}

	cfg := config.NewConfig()
	app, err := app.NewApplication(cfg)
	if err != nil {
		logger.Log.Fatal().Msgf("Failed to initialize application:%v", err)
	}
	if err := app.Run(); err != nil {
		logger.Log.Fatal().Msgf("Failed to run application:%v", err)
	}

}
