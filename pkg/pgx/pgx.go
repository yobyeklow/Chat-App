package pgx

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
	"web_socket/pkg/logger"

	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
)

type PgxZerologTracer struct {
	Logger         zerolog.Logger
	SlowQueryLimit time.Duration
}
type QueryInfo struct {
	QueryName     string
	OperationType string
	CleanSQL      string
	OrginalSQL    string
}

// -- name: CreateUser :one
var (
	sqlcNameRegex = regexp.MustCompile(`-- name:\s*(\w+)\s*:(\w+)`)
	spaceRegex    = regexp.MustCompile(`\s+`)
	commentRegex  = regexp.MustCompile(`-- [^\r\n]*`)
)

func parseSQL(sql string) QueryInfo {
	info := QueryInfo{
		OrginalSQL: sql,
	}

	matches := sqlcNameRegex.FindStringSubmatch(sql)
	if len(matches) == 3 {
		info.QueryName = matches[1]
		info.OperationType = strings.ToUpper(matches[2])
	}

	cleanSQL := commentRegex.ReplaceAllString(sql, "")
	cleanSQL = strings.TrimSpace(cleanSQL)
	cleanSQL = spaceRegex.ReplaceAllString(cleanSQL, "")

	info.CleanSQL = cleanSQL
	return info
}

func formatArg(arg any) string {
	val := reflect.ValueOf(arg)
	if arg == nil || (val.Kind() == reflect.Pointer && val.IsNil()) {
		return "NULL"
	}

	if val.Kind() == reflect.Pointer {
		arg = val.Elem().Interface()
	}

	switch v := arg.(type) {
	case string:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''"))
	case bool:
		return fmt.Sprintf("%t", v)
	case int8, int16, int32, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case time.Time:
		return fmt.Sprintf("'%s'", v.Format("2006-01-02T15:04:05Z07:00"))
	default:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(fmt.Sprintf("%v", v), "'", "''"))

	}
}

func replacePlaceHolders(sql string, args []any) string {
	for index, argument := range args {
		placeholder := fmt.Sprintf("$%d", index+1)
		sql = strings.ReplaceAll(sql, placeholder, formatArg(argument))
	}
	return sql
}
func (pgxTracer *PgxZerologTracer) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	sql, _ := data["sql"].(string)
	args, _ := data["args"].([]any)
	duration, _ := data["time"].(time.Duration)

	queryInfo := parseSQL(sql)
	var finalSql string
	if len(args) > 0 {
		finalSql = replacePlaceHolders(queryInfo.CleanSQL, args)
	} else {
		finalSql = queryInfo.CleanSQL
	}

	baseLogger := pgxTracer.Logger.With().
		Str("tracer_id", logger.GetTraceID(ctx)).
		Dur("duration", duration).
		Str("sql_original", queryInfo.OrginalSQL).
		Str("query_name", queryInfo.QueryName).
		Str("operation_time", queryInfo.OperationType).
		Str("sql", finalSql).
		Interface("args", args)

	logger := baseLogger.Logger()
	if msg == "Query" && duration > pgxTracer.SlowQueryLimit {
		logger.Warn().Str("event", "Slow Query").Msg("Slow SQL Query")
		return
	}
	if msg == "Query" {
		logger.Warn().Str("event", "Query").Msg("Excuted SQL")
		return
	}
}
