package db

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"rental-management/internal/config"
)

var RedisClient *redis.Client

// InitRedis 初始化Redis连接
func InitRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx := context.Background()
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect redis: %w", err)
	}

	return RedisClient, nil
}

// GetRedis 获取Redis连接
func GetRedis() *redis.Client {
	return RedisClient
}
