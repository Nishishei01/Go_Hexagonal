package redisAdapter

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Nishishei01/Go_Hexagonal/internal/ports"
	"github.com/redis/go-redis/v9"
)

type RedisCacheAdapter struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCacheAdapter(client *redis.Client) ports.CacheRepository {
	return &RedisCacheAdapter{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *RedisCacheAdapter) Set(key string, value any, expiration time.Duration) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(r.ctx, key, p, expiration).Err()
}

func (r *RedisCacheAdapter) Get(key string, dest any) error {
	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (r *RedisCacheAdapter) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}
