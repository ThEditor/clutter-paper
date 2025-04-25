package storage

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
)

type ClickHouseStorage struct {
	db          *sql.DB
	batchSize   int
	eventBatch  []EventData
	batchMutex  sync.Mutex
	flushTicker *time.Ticker
	done        chan bool
}

func NewClickHouseStorage(dsn string, batchSize int, flushInterval time.Duration) (*ClickHouseStorage, error) {
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ClickHouse: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping ClickHouse: %w", err)
	}

	storage := &ClickHouseStorage{
		db:         db,
		batchSize:  batchSize,
		eventBatch: make([]EventData, 0, batchSize),
		done:       make(chan bool),
	}

	if err := storage.ensureTables(); err != nil {
		return nil, fmt.Errorf("failed to ensure tables: %w", err)
	}

	if flushInterval > 0 {
		storage.flushTicker = time.NewTicker(flushInterval)
		go func() {
			for {
				select {
				case <-storage.flushTicker.C:
					if err := storage.Flush(); err != nil {
						fmt.Printf("Error flushing batch: %v\n", err)
					}
				case <-storage.done:
					return
				}
			}
		}()
	}

	return storage, nil
}

func (s *ClickHouseStorage) Close() error {
	if s.flushTicker != nil {
		s.flushTicker.Stop()
		s.done <- true
	}

	if err := s.Flush(); err != nil {
		return fmt.Errorf("failed to flush events on close: %w", err)
	}

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
	s.batchMutex.Lock()
	defer s.batchMutex.Unlock()

	s.eventBatch = append(s.eventBatch, data)

	if len(s.eventBatch) >= s.batchSize {
		return s.flushLocked()
	}

	return nil
}

func (s *ClickHouseStorage) Flush() error {
	s.batchMutex.Lock()
	defer s.batchMutex.Unlock()

	return s.flushLocked()
}

func (s *ClickHouseStorage) flushLocked() error {
	if len(s.eventBatch) == 0 {
		return nil
	}

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

	now := time.Now()
	for _, event := range s.eventBatch {
		_, err = stmt.Exec(
			event.VisitorIP,
			event.VisitorUserAgent,
			event.SiteID,
			event.Referrer,
			now,
			event.Page,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute statement: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.eventBatch = s.eventBatch[:0]

	return nil
}
