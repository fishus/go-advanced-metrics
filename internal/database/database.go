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
	Close()
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	AcquireFunc(ctx context.Context, f func(*pgxpool.Conn) error) error
	AcquireAllIdle(ctx context.Context) []*pgxpool.Conn
	Reset()
	Config() *pgxpool.Config
	Stat() *pgxpool.Stat
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	Ping(ctx context.Context) error
}

var _ Connector = (*pgxpool.Pool)(nil)

var pool Connector

type NamedArgs = pgx.NamedArgs
type StatementDescription = pgconn.StatementDescription
type Rows = pgx.Rows

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
