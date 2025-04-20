package common

import (
	"fmt"

	"github.com/ThEditor/clutter-paper/internal/storage"
)

type Server struct {
	Clickhouse *storage.ClickHouseStorage
	Redis      *storage.RedisStorage
	Postgres   *storage.PostgresStorage
}

func CheckSiteID(siteID string, s *Server) error {
	exists, err := s.Redis.SiteIDExists(siteID)
	if err != nil {
		return fmt.Errorf("redis check failed: %v", err)
	}

	if exists {
		return nil
	}

	exists, err = s.Postgres.SiteIDExists(siteID)
	if err != nil {
		return fmt.Errorf("postgres check failed: %v", err)
	}

	if exists {
		if err := s.Redis.AddSiteID(siteID); err != nil {
			return fmt.Errorf("failed to cache site ID: %v", err)
		}
	}

	return nil
}
