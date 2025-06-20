package streaming

import (
	"context"
	"encoding/json"
	"taskflow/internal/domain/entities"
	"taskflow/internal/domain/repositories"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient(redisURL string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
}

type streamRepository struct {
	client *redis.Client
}

func NewStreamRepository(client *redis.Client) repositories.StreamRepository {
	return &streamRepository{
		client: client,
	}
}

func (r *streamRepository) PublishTodoCreated(ctx context.Context, todo *entities.TodoItem) error {
	todoJSON, err := json.Marshal(todo)
	if err != nil {
		return err
	}

	return r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: "todo-created",
		Values: map[string]interface{}{
			"todo": string(todoJSON),
		},
	}).Err()
}
