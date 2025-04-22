package redis

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/redis/go-redis/v9"
	"time"
)

type Client struct {
	*redis.Client
}

func NewClient(cfg *config.RedisConfig) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %v", err)
	}

	return &Client{rdb}, nil
}
