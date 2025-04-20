package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	pool *pgxpool.Pool
	ctx  context.Context
}

func NewPostgresStorage(ctx context.Context, dsn string) (*PostgresStorage, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return &PostgresStorage{pool: pool, ctx: ctx}, nil
}

func (ps *PostgresStorage) SiteIDExists(siteID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM sites WHERE id = $1)`
	var exists bool

	err := ps.pool.QueryRow(ps.ctx, query, siteID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if site ID exists: %w", err)
	}

	return exists, nil
}

func (ps *PostgresStorage) Close() {
	ps.pool.Close()
}
