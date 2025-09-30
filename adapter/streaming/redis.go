package streaming

import (
	"context"
	"encoding/json"
	"taskflow/internal/domain/shared"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient(redisURL string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
}

type redisMessaging struct {
	client *redis.Client
}

func NewRedisMessaging(client *redis.Client) shared.Messaging {
	return &redisMessaging{
		client: client,
	}
}

func (r *redisMessaging) Publish(ctx context.Context, topic string, message interface{}) error {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: topic,
		Values: map[string]interface{}{
			"data": string(messageJSON),
		},
	}).Err()
}

func (r *redisMessaging) PublishWithKey(ctx context.Context, topic string, key string, message interface{}) error {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: topic,
		ID:     key,
		Values: map[string]interface{}{
			"data": string(messageJSON),
		},
	}).Err()
}
