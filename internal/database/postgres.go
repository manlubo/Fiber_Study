package database

import (
	"context"
	"fmt"
	"time"

	"study/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(cfg *config.Postgres) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)

	// 연결 풀 생성
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	// 연결 풀 설정
	pool.Config().MaxConns = int32(cfg.MaxOpenConns)
	pool.Config().MinConns = int32(cfg.MaxIdleConns)
	pool.Config().MaxConnLifetime = 30 * time.Minute

	// 연결 테스트
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}
