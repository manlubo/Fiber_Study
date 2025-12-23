package database

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"study/internal/config"

	"study/pkg/util"

	"github.com/exaring/otelpgx"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateDsn(cfg *config.Postgres) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)
}

func NewPostgres(cfg *config.Postgres) (*pgxpool.Pool, error) {
	dsn := CreateDsn(cfg)

	// config 파싱
	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	// tracer 설정
	poolCfg.ConnConfig.Tracer = otelpgx.NewTracer()

	// 연결 풀 설정
	poolCfg.MaxConns = int32(cfg.MaxOpenConns)
	poolCfg.MinConns = int32(cfg.MaxIdleConns)
	poolCfg.MaxConnLifetime = 30 * time.Minute

	// 연결 풀 생성
	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, err
	}

	// 연결 테스트
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}

func RunMigration(dbURL string) error {
	path := filepath.ToSlash(util.GetPath("migrations"))

	m, err := migrate.New(
		"file://"+path,
		dbURL,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
