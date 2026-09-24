package db

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Connect(ctx context.Context) (*pgxpool.Pool, error) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = "postgres://golink:golink@127.0.0.1:5433/golink?sslmode=disable"
	}
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}

func Migrate(pool *pgxpool.Pool) error {
	dir := os.Getenv("MIGRATIONS_DIR")
	if dir == "" {
		dir = "migrations"
	}
	sqlDB := stdlib.OpenDBFromPool(pool)
	defer sqlDB.Close()
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(sqlDB, dir)
}

// SeedIfEmpty loads the four seed files once. Fixed UUIDs make a second run a no-op via ON CONFLICT.
func SeedIfEmpty(ctx context.Context, pool *pgxpool.Pool) error {
	var n int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM ev_models`).Scan(&n); err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	dir := os.Getenv("SEED_DIR")
	if dir == "" {
		dir = filepath.Join("..", "seed")
	}
	for _, name := range []string{"ev_models.sql", "price_listings.sql", "charging_stations.sql", "import_policy.sql"} {
		b, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return fmt.Errorf("seed %s: %w", name, err)
		}
		if _, err := pool.Exec(ctx, string(b)); err != nil {
			return fmt.Errorf("seed %s: %w", name, err)
		}
	}
	return nil
}
