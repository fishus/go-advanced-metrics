package database

import (
	"context"
	
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/fishus/go-advanced-metrics/internal/logger"
)

type ConnPool = pgxpool.Pool

var pool *ConnPool

// Open connect to database (Postgres)
// Don't forget to call defer dbpool.Close()
func Open(ctx context.Context, dsn string) *pgxpool.Pool {
	p, err := pgxpool.New(ctx, dsn)
	if err != nil {
		logger.Log.Error(err.Error(), logger.String("event", "connect database"), logger.String("dsn", dsn))
	}
	pool = p

	return pool
}

func Pool() *pgxpool.Pool {
	return pool
}
