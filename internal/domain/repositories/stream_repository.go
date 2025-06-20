package repositories

import (
	"context"
	"taskflow/internal/domain/entities"
)

type StreamRepository interface {
	PublishTodoCreated(ctx context.Context, todo *entities.TodoItem) error
}
