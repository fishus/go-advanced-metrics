package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/fishus/go-advanced-metrics/internal/logger"
)

type Connector interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Ping(context.Context) error
	Close()
}

var _ Connector = (*pgxpool.Pool)(nil)

var pool Connector

var ErrNotConnected = errors.New("database connection not established")

var ErrNoRows = pgx.ErrNoRows

// Open connect to database (Postgres)
// Don't forget to call defer dbpool.Close()
func Open(ctx context.Context, dsn string) Connector {
	p, err := pgxpool.New(ctx, dsn)

	if err != nil {
		logger.Log.Error(err.Error(), logger.String("event", "connect database"), logger.String("dsn", dsn))
	}
	pool = p

	return pool
}

func Pool() (Connector, error) {
	if pool == nil {
		return nil, ErrNotConnected
	}
	return pool, nil
}
