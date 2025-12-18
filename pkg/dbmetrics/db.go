package dbmetrics

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *DB {
	return &DB{pool: pool}
}

func (db *DB) Query(
	ctx context.Context,
	query string,
	args ...any,
) (pgx.Rows, error) {

	start := time.Now()
	rows, err := db.pool.Query(ctx, query, args...)
	elapsed := time.Since(start).Milliseconds()

	if m := FromContext(ctx); m != nil {
		m.Record(query, elapsed)
	}

	return rows, err
}

func (db *DB) QueryRow(
	ctx context.Context,
	query string,
	args ...any,
) pgx.Row {

	start := time.Now()
	row := db.pool.QueryRow(ctx, query, args...)
	elapsed := time.Since(start).Milliseconds()

	if m := FromContext(ctx); m != nil {
		m.Record(query, elapsed)
	}

	return row
}

func (db *DB) Exec(
	ctx context.Context,
	query string,
	args ...any,
) (pgconn.CommandTag, error) {

	start := time.Now()
	tag, err := db.pool.Exec(ctx, query, args...)
	elapsed := time.Since(start).Milliseconds()

	if m := FromContext(ctx); m != nil {
		m.Record(query, elapsed)
	}

	return tag, err
}
