package db

import (
	"context"
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping pgx pool: %w", err)
	}

	return pool, nil
}

func RunMigrations(pool *pgxpool.Pool) error {
	goose.SetBaseFS(migrations)

	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	err = goose.Up(stdlib.OpenDBFromPool(pool), "migrations")
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
