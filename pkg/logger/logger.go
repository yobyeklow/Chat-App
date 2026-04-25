package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

type contextKey string

const TraceIdKey contextKey = "trace_id"

var Log *zerolog.Logger

type LoggerConfig struct {
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	Level      string
	IsDev      string
}

type PrettyJSONWriter struct {
	Writer io.Writer
}

func InitLogger(config LoggerConfig) {
	Log = NewLogger(config)
}
func (w PrettyJSONWriter) Write(p []byte) (n int, err error) {
	var prettyJson bytes.Buffer

	if err := json.Indent(&prettyJson, p, "", " "); err != nil {
		return w.Writer.Write(p)
	}
	return w.Writer.Write(prettyJson.Bytes())
}
func NewLogger(config LoggerConfig) *zerolog.Logger {
	zerolog.TimestampFieldName = time.RFC3339
	level, err := zerolog.ParseLevel(config.Level)

	if err != nil {
		level = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(level)
	var writer io.Writer
	if config.IsDev == "developement" {
		writer = PrettyJSONWriter{
			Writer: writer,
		}
	} else {
		writer = &lumberjack.Logger{
			Filename:   config.Filename,
			MaxSize:    config.MaxAge,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		}
	}
	logger := zerolog.New(writer).With().Timestamp().Logger()
	return &logger
}
func GetTraceID(ctx context.Context) string {
	if traceId, ok := ctx.Value(TraceIdKey).(string); ok {
		return traceId
	}
	return ""
}
