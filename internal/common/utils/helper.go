package utils

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"web_socket/pkg/logger"

	"github.com/rs/zerolog"
)

func GetEnv(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}
func GetEnvInt(key string, defaultVal int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return intVal
}
func NewLoggerWithPath(fileName string, level string) *zerolog.Logger {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("❌ Unable to get working dir:", err)
	}
	logDir := filepath.Join(cwd, "internal/logs/", fileName)
	config := logger.LoggerConfig{
		Level:      level,
		Filename:   logDir,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     5,
		Compress:   true,
		IsDev:      GetEnv("APP_STATUS", "developement"),
	}
	return logger.NewLogger(config)
}
func MustGetWorkingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Unable to get working dir:", err)
	}
	return dir
}
func MapRoleText(status int) string {
	switch status {
	case 1:
		return "User"
	case 2:
		return "Adminstrator"
	default:
		return "None"
	}
}
