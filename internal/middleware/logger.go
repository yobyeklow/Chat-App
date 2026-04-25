package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"
	"web_socket/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type CustomReponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *CustomReponseWriter) Write(data []byte) (n int, err error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

func LoggerMiddleware(httpLogger *zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		contentType := ctx.GetHeader("Content-Type")
		requestBody := make(map[string]any)

		var formFiles []map[string]any

		if strings.HasPrefix(contentType, "multipart/form-data") {
			if err := ctx.Request.ParseMultipartForm(32 << 20); err == nil && ctx.Request.MultipartForm != nil {
				for key, val := range ctx.Request.MultipartForm.Value {
					if len(val) > 1 {
						requestBody[key] = val[0]
					} else {
						requestBody[key] = val
					}
				}
				for field, files := range ctx.Request.MultipartForm.File {
					for _, file := range files {
						formFiles = append(formFiles, map[string]any{
							"field":        field,
							"filename":     file.Filename,
							"size":         formatFileSize(file.Size),
							"content_type": file.Header.Get("Content-Type"),
						})
					}
				}
			}
		} else {
			bodyBytes, err := io.ReadAll(ctx.Request.Body)
			if err != nil {
				httpLogger.Error().Err(err).Msg("Failed to read request body")
			}

			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			if strings.HasPrefix(contentType, "application/json") {
				_ = json.Unmarshal(bodyBytes, &requestBody)
			} else {
				values, _ := url.ParseQuery(string(bodyBytes))
				for key, vals := range values {
					if len(vals) > 1 {
						requestBody[key] = vals[0]
					} else {
						requestBody[key] = vals
					}
				}
			}
			customWriter := &CustomReponseWriter{
				ResponseWriter: ctx.Writer,
				body:           bytes.NewBufferString(""),
			}
			ctx.Writer = customWriter
			ctx.Next()

			duration := time.Since(start)
			statusCode := ctx.Writer.Status()
			responseContentType := ctx.Writer.Header().Get("Content-Type")
			responseBodyRaw := customWriter.body.String()

			var responseBodyParsed any
			if strings.HasPrefix(responseContentType, "image/") {
				responseBodyParsed = "[BINARY DATA]"
			} else if strings.HasPrefix(responseContentType, "application/json") ||
				strings.HasPrefix(strings.TrimSpace(responseBodyRaw), "{") ||
				strings.HasPrefix(strings.TrimSpace(responseBodyRaw), "[") {
				if err := json.Unmarshal([]byte(responseBodyRaw), &responseBodyParsed); err != nil {
					responseBodyParsed = responseBodyRaw
				}
			} else {
				responseBodyParsed = responseBodyRaw
			}
			logEvent := httpLogger.Info()
			if statusCode >= 500 {
				logEvent = httpLogger.Error()
			} else if statusCode >= 400 {
				logEvent = httpLogger.Warn()
			}

			logEvent.Str("trace_id", logger.GetTraceID(ctx.Request.Context())).
				Str("method", ctx.Request.Method).
				Str("path", ctx.Request.URL.Path).
				Str("query", ctx.Request.URL.RawQuery).
				Str("client_ip", ctx.ClientIP()).
				Str("user_agent", ctx.GetHeader("User-Agent")).
				Str("referer", ctx.Request.Referer()).
				Str("protocol", ctx.Request.Proto).
				Str("host", ctx.Request.Host).
				Str("remote_addr", ctx.Request.RemoteAddr).
				Str("request_uri", ctx.Request.RequestURI).
				Interface("headers", ctx.Request.Header).
				Interface("request_body", requestBody).
				Int("status_code", statusCode).
				Interface("response_body", responseBodyParsed).
				Int64("duration_ms", duration.Milliseconds()).
				Msg("HTTP Request Log")

		}
	}
}
func formatFileSize(size int64) string {
	switch {
	case size >= 1<<20:
		return fmt.Sprintf("%.2f MB", float64(size)/(1<<20))
	case size >= 1<<10:
		return fmt.Sprintf("%.2f KB", float64(size)/(1<<10))
	default:
		return fmt.Sprintf("%d B", size)
	}
}
