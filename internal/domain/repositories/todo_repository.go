package repositories

import (
	"context"
	"github.com/google/uuid"
	"taskflow/internal/domain/entities"
)

type TodoRepository interface {
	Create(ctx context.Context, todo *entities.TodoItem) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.TodoItem, error)
	List(ctx context.Context, limit, offset int) ([]*entities.TodoItem, error)
	Update(ctx context.Context, todo *entities.TodoItem) error
	Delete(ctx context.Context, id uuid.UUID) error
}
