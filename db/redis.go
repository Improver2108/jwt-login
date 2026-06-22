package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func NewRedis() *Redis {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		log.Fatal("No redis addr provided")
	}
	rdb := redis.NewClient(&redis.Options{Addr: addr})

	return &Redis{Client: rdb}
}

func (r *Redis) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return r.Client.Set(ctx, key, value, ttl).Err()
}

func (r *Redis) Del(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

func (r *Redis) Get(ctx context.Context, key string) (any, error) {
	return r.Client.Get(ctx, key).Result()
}
