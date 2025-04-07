package storage

import (
	"database/sql"
	"fmt"
	"time"

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

	storage := &ClickHouseStorage{db: db}
	if err := storage.ensureTables(); err != nil {
		return nil, fmt.Errorf("failed to ensure tables: %w", err)
	}

	return storage, nil
}

func (s *ClickHouseStorage) Close() error {
	return s.db.Close()
}

func (s *ClickHouseStorage) ensureTables() error {
	tableSchemas := []string{
		`
		CREATE TABLE IF NOT EXISTS events (
			id UUID DEFAULT generateUUIDv4(),
			visitor_ip String,
			visitor_user_agent String,
			site_id String,
			referrer String,
			created_on DateTime,
			page String,
			PRIMARY KEY (id)
		) ENGINE = MergeTree()
		`,
		// others
	}

	for _, schema := range tableSchemas {
		if _, err := s.db.Exec(schema); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	return nil
}

type EventData struct {
	VisitorIP        string
	VisitorUserAgent string
	SiteID           string
	Referrer         string
	Page             string
}

func (s *ClickHouseStorage) InsertEvent(data EventData) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	stmt, err := tx.Prepare(`
        INSERT INTO events (
            visitor_ip,
            visitor_user_agent,
            site_id,
            referrer,
            created_on,
            page
        ) VALUES (
            ?, ?, ?, ?, ?, ?
        )
    `)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		data.VisitorIP,
		data.VisitorUserAgent,
		data.SiteID,
		data.Referrer,
		time.Now(),
		data.Page,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
