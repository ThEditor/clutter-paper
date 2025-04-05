package main

import (
	"github.com/ThEditor/clutter-paper/internal/api"
	"github.com/ThEditor/clutter-paper/internal/config"
	"github.com/ThEditor/clutter-paper/internal/storage"
)

func main() {
	cfg := config.Load()

	var store storage.Storage
	if cfg.STORAGE_MODE == "clickhouse" {
		var err error
		store, err = storage.NewClickHouseStorage(cfg.DATABASE_URL)
		if err != nil {
			panic(err)
		}
		defer store.Close()
	}

	api.Start(cfg.BIND_ADDRESS, cfg.PORT)
}
