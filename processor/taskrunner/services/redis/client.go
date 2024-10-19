package redis

import (
	"context"
	"fmt"
	"ylem_taskrunner/config"

	"github.com/go-redis/redis/v8"
)

var instance *redis.Client

func Init(ctx context.Context) {
	cfg := config.Cfg().Redis

	instance = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0,
	}).WithContext(ctx)

	_, err := instance.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
}

func Instance() *redis.Client {
	if instance == nil {
		panic("Redis client not initialized")
	}
	return instance
}
