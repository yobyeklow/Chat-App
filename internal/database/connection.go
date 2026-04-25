package database

import (
	"context"
	"fmt"
	"time"
	"web_socket/internal/config"
	"web_socket/internal/database/sqlc"
	"web_socket/internal/utils"
	"web_socket/pkg/logger"
	"web_socket/pkg/pgx"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
)

var DBPool *pgxpool.Pool
var DB *sqlc.Queries

func InitDB() error {
	connStr := config.NewConfig().DNS()
	sqlLogger := utils.NewLoggerWithPath("sql.log", "info")
	conf, err := pgxpool.ParseConfig(connStr)
	fmt.Println(conf.ConnConfig.ConnString())
	if err != nil {
		return fmt.Errorf("Error parsing DB config: %v", err)
	}

	conf.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger: &pgx.PgxZerologTracer{
			Logger:         *sqlLogger,
			SlowQueryLimit: 500 * time.Millisecond,
		},
		LogLevel: tracelog.LogLevelDebug,
	}
	conf.MaxConns = 50
	conf.MinConns = 5
	conf.MaxConnIdleTime = 5 * time.Minute
	conf.MaxConnLifetime = 30 * time.Minute
	conf.HealthCheckPeriod = 1 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	DBPool, err = pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return fmt.Errorf("Error creating Database pool:%v", err)
	}
	DB = sqlc.New(DBPool)
	if err := DBPool.Ping(ctx); err != nil {
		return fmt.Errorf("DB Ping error: %v", err)
	}

	logger.Log.Info().Msg("Connected Database successfully!")
	return nil
}
