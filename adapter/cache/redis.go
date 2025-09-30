package cache

import (
	"context"
	"encoding/json"
	"strconv"
	"taskflow/internal/domain/shared"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) shared.Cache {
	return &redisCache{
		client: client,
	}
}

func (r *redisCache) Get(ctx context.Context, key string) (string, error) {
	result := r.client.Get(ctx, key)
	if result.Err() != nil {
		if result.Err() == redis.Nil {
			return "", shared.ErrNotFound
		}
		return "", result.Err()
	}
	return result.Val(), nil
}

func (r *redisCache) Set(ctx context.Context, key string, value string, ttl int) error {
	duration := time.Duration(ttl) * time.Second
	return r.client.Set(ctx, key, value, duration).Err()
}

func (r *redisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *redisCache) Exists(ctx context.Context, key string) (bool, error) {
	result := r.client.Exists(ctx, key)
	if result.Err() != nil {
		return false, result.Err()
	}
	return result.Val() > 0, nil
}

// Helper methods for common operations

func (r *redisCache) GetJSON(ctx context.Context, key string, dest interface{}) error {
	value, err := r.Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(value), dest)
}

func (r *redisCache) SetJSON(ctx context.Context, key string, value interface{}, ttl int) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.Set(ctx, key, string(data), ttl)
}

func (r *redisCache) GetInt(ctx context.Context, key string) (int64, error) {
	value, err := r.Get(ctx, key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(value, 10, 64)
}

func (r *redisCache) SetInt(ctx context.Context, key string, value int64, ttl int) error {
	return r.Set(ctx, key, strconv.FormatInt(value, 10), ttl)
}

func (r *redisCache) Increment(ctx context.Context, key string) (int64, error) {
	result := r.client.Incr(ctx, key)
	return result.Val(), result.Err()
}

func (r *redisCache) IncrementBy(ctx context.Context, key string, value int64) (int64, error) {
	result := r.client.IncrBy(ctx, key, value)
	return result.Val(), result.Err()
}

func (r *redisCache) Decrement(ctx context.Context, key string) (int64, error) {
	result := r.client.Decr(ctx, key)
	return result.Val(), result.Err()
}

func (r *redisCache) DecrementBy(ctx context.Context, key string, value int64) (int64, error) {
	result := r.client.DecrBy(ctx, key, value)
	return result.Val(), result.Err()
}
