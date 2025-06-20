package tests

import (
	"context"
	"testing"
	"time"

	"taskflow/internal/domain/entities"
	"taskflow/internal/service"

	"github.com/google/uuid"
)

type benchMockTodoRepo struct {
	CreateFn func(ctx context.Context, todo *entities.TodoItem) error
	ListFn   func(ctx context.Context, limit, offset int) ([]*entities.TodoItem, error)
}

func (m *benchMockTodoRepo) Create(ctx context.Context, todo *entities.TodoItem) error {
	return m.CreateFn(ctx, todo)
}
func (m *benchMockTodoRepo) GetByID(ctx context.Context, id uuid.UUID) (*entities.TodoItem, error) {
	return nil, nil
}
func (m *benchMockTodoRepo) List(ctx context.Context, limit, offset int) ([]*entities.TodoItem, error) {
	return m.ListFn(ctx, limit, offset)
}
func (m *benchMockTodoRepo) Update(ctx context.Context, todo *entities.TodoItem) error { return nil }
func (m *benchMockTodoRepo) Delete(ctx context.Context, id uuid.UUID) error            { return nil }

type benchMockStreamRepo struct{}

func (m *benchMockStreamRepo) PublishTodoCreated(ctx context.Context, todo *entities.TodoItem) error {
	return nil
}

func BenchmarkCreateTodo(b *testing.B) {
	todoRepo := &benchMockTodoRepo{
		CreateFn: func(ctx context.Context, todo *entities.TodoItem) error { return nil },
	}
	streamRepo := &benchMockStreamRepo{}
	uc := service.NewTodoService(todoRepo, streamRepo)

	req := &entities.CreateTodoRequest{
		Description: "Benchmark todo",
		DueDate:     time.Now().Add(24 * time.Hour),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := uc.CreateTodo(context.Background(), req)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

func BenchmarkListTodos(b *testing.B) {
	todoRepo := &benchMockTodoRepo{
		ListFn: func(ctx context.Context, limit, offset int) ([]*entities.TodoItem, error) {
			return []*entities.TodoItem{{Description: "Benchmark todo"}}, nil
		},
	}
	streamRepo := &benchMockStreamRepo{}
	uc := service.NewTodoService(todoRepo, streamRepo)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := uc.ListTodos(context.Background(), 10, 0)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}
