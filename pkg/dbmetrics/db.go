package dbmetrics

import (
	"context"
	"study/internal/shared/db"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxDB struct {
	pool *pgxpool.Pool
}

type PgxTx struct {
	tx pgx.Tx
}

func New(pool *pgxpool.Pool) *PgxDB {
	return &PgxDB{pool: pool}
}

func (db *PgxDB) Query(
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

func (db *PgxDB) QueryRow(
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

func (db *PgxDB) Exec(
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

func (db *PgxDB) Begin(ctx context.Context) (db.Tx, error) {
	start := time.Now()

	tx, err := db.pool.Begin(ctx)
	elapsed := time.Since(start).Milliseconds()

	if m := FromContext(ctx); m != nil {
		m.Record("BEGIN", elapsed)
	}

	if err != nil {
		return nil, err
	}

	return &PgxTx{tx: tx}, nil
}

func (t *PgxTx) Query(
	ctx context.Context,
	query string,
	args ...any,
) (pgx.Rows, error) {

	start := time.Now()
	rows, err := t.tx.Query(ctx, query, args...)
	elapsed := time.Since(start).Milliseconds()

	if m := FromContext(ctx); m != nil {
		m.Record(query, elapsed)
	}

	return rows, err
}

func (t *PgxTx) QueryRow(
	ctx context.Context,
	query string,
	args ...any,
) pgx.Row {

	start := time.Now()
	row := t.tx.QueryRow(ctx, query, args...)
	elapsed := time.Since(start).Milliseconds()

	if m := FromContext(ctx); m != nil {
		m.Record(query, elapsed)
	}

	return row
}

func (t *PgxTx) Exec(
	ctx context.Context,
	query string,
	args ...any,
) (pgconn.CommandTag, error) {

	start := time.Now()
	tag, err := t.tx.Exec(ctx, query, args...)
	elapsed := time.Since(start).Milliseconds()

	if m := FromContext(ctx); m != nil {
		m.Record(query, elapsed)
	}

	return tag, err
}

func (t *PgxTx) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *PgxTx) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}
