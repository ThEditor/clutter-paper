package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/ClickHouse/clickhouse-go"
)

type ClickHouseStorage struct {
	db *sql.DB
}

func NewClickHouseStorage(dsn string) (*ClickHouseStorage, error) {
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ClickHouse: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping ClickHouse: %w", err)
	}

	return &ClickHouseStorage{db: db}, nil
}

func (s *ClickHouseStorage) Close() error {
	return s.db.Close()
}
