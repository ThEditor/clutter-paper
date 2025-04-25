package main

import (
	"context"
	"time"

	"github.com/ThEditor/clutter-paper/internal/api"
	"github.com/ThEditor/clutter-paper/internal/config"
	"github.com/ThEditor/clutter-paper/internal/storage"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	store, err := storage.NewClickHouseStorage(cfg.DATABASE_URL, 500, 15*time.Second)
	if err != nil {
		panic(err)
	}
	defer store.Close()

	redisStore, err := storage.NewRedisStorage(ctx, cfg.REDIS_URL)
	if err != nil {
		panic(err)
	}
	defer redisStore.Close()

	postgresStore, err := storage.NewPostgresStorage(ctx, cfg.POSTGRES_URL)
	if err != nil {
		panic(err)
	}
	defer postgresStore.Close()

	api.Start(cfg.BIND_ADDRESS, cfg.PORT, store, redisStore, postgresStore)
}
