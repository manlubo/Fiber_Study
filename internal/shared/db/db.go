package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Execer interface {
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
}

type DB interface {
	Execer
	Begin(ctx context.Context) (Tx, error)
}

type Tx interface {
	Execer
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
