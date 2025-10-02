package todo

import (
	"context"
	"taskflow/internal/domain/shared"

	"github.com/google/uuid"
)

// TodoService defines the todo service interface
type TodoService interface {
	CreateTodo(ctx context.Context, req *CreateTodoRequest) (*TodoItem, error)
	GetTodo(ctx context.Context, id uuid.UUID) (*TodoItem, error)
	ListTodos(ctx context.Context, limit, offset int) ([]*TodoItem, error)
	UpdateTodo(ctx context.Context, id uuid.UUID, req *UpdateTodoRequest) (*TodoItem, error)
	DeleteTodo(ctx context.Context, id uuid.UUID) error
}

// Repository defines the todo repository interface
type Repository interface {
	Create(ctx context.Context, todo *TodoItem) error
	GetByID(ctx context.Context, id uuid.UUID) (*TodoItem, error)
	List(ctx context.Context, limit, offset int) ([]*TodoItem, error)
	Update(ctx context.Context, todo *TodoItem) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// Messaging defines the messaging interface (uses shared messaging port)
type Messaging = shared.Messaging

// Cache defines the cache interface (uses shared cache port)
type Cache = shared.Cache
