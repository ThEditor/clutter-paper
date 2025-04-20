package storage

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisStorage struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStorage(ctx context.Context, connUrl string) (*RedisStorage, error) {
	opt, err := redis.ParseURL(connUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	client := redis.NewClient(opt)

	_, err1 := client.Ping(ctx).Result()
	if err1 != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &RedisStorage{
		client: client,
		ctx:    ctx,
	}, nil
}

func (r *RedisStorage) AddSiteID(siteID string) error {
	return r.client.SAdd(r.ctx, "site_ids", siteID).Err()
}

func (r *RedisStorage) SiteIDExists(siteID string) (bool, error) {
	return r.client.SIsMember(r.ctx, "site_ids", siteID).Result()
}

func (r *RedisStorage) Close() error {
	return r.client.Close()
}
