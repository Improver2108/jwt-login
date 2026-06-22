package cache

import (
	"context"
	"time"
)

type JTIStore interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Del(ctx context.Context, key string) error
	Get(ctx context.Context, key string) (any, error)
}

type JTICache struct {
	cache JTIStore
}

func NewJTICache(cache JTIStore) *JTICache {
	return &JTICache{cache: cache}
}

func (j *JTICache) SetJTI(ctx context.Context, key, userID string, ttl time.Duration) error {
	return j.cache.Set(ctx, key, userID, ttl)
}

func (j *JTICache) DelJTI(ctx context.Context, key string) error {
	return j.cache.Del(ctx, key)
}

func (j *JTICache) GetJTI(ctx context.Context, key string) (string, error) {
	value, err := j.cache.Get(ctx, key)

	if err != nil {
		return "", err
	}

	res, ok := value.(string)
	if !ok {
		return "", err
	}

	return res, nil
}
