package tests

import (
	"context"
	"testing"
	"time"

	"taskflow/internal/domain/todo"

	"github.com/google/uuid"
)

type benchMockTodoRepo struct {
	CreateFn func(ctx context.Context, todoItem *todo.TodoItem) error
	ListFn   func(ctx context.Context, limit, offset int) ([]*todo.TodoItem, error)
}

func (m *benchMockTodoRepo) Create(ctx context.Context, todoItem *todo.TodoItem) error {
	return m.CreateFn(ctx, todoItem)
}
func (m *benchMockTodoRepo) GetByID(ctx context.Context, id uuid.UUID) (*todo.TodoItem, error) {
	return nil, nil
}
func (m *benchMockTodoRepo) List(ctx context.Context, limit, offset int) ([]*todo.TodoItem, error) {
	return m.ListFn(ctx, limit, offset)
}
func (m *benchMockTodoRepo) Update(ctx context.Context, todoItem *todo.TodoItem) error { return nil }
func (m *benchMockTodoRepo) Delete(ctx context.Context, id uuid.UUID) error            { return nil }

type benchMockMessaging struct{}

func (m *benchMockMessaging) Publish(ctx context.Context, topic string, message interface{}) error {
	return nil
}

func (m *benchMockMessaging) PublishWithKey(ctx context.Context, topic string, key string, message interface{}) error {
	return nil
}

type benchMockCache struct{}

func (m *benchMockCache) Get(ctx context.Context, key string) (string, error) { return "", nil }
func (m *benchMockCache) Set(ctx context.Context, key string, value string, ttl int) error {
	return nil
}
func (m *benchMockCache) Delete(ctx context.Context, key string) error         { return nil }
func (m *benchMockCache) Exists(ctx context.Context, key string) (bool, error) { return false, nil }

func BenchmarkCreateTodo(b *testing.B) {
	todoRepo := &benchMockTodoRepo{
		CreateFn: func(ctx context.Context, todoItem *todo.TodoItem) error { return nil },
	}
	messaging := &benchMockMessaging{}
	cache := &benchMockCache{}
	service := todo.NewTodoService(todoRepo, messaging, cache)

	req := &todo.CreateTodoRequest{
		Description: "Benchmark todo",
		DueDate:     time.Now().Add(24 * time.Hour),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.CreateTodo(context.Background(), req)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

func BenchmarkListTodos(b *testing.B) {
	todoRepo := &benchMockTodoRepo{
		ListFn: func(ctx context.Context, limit, offset int) ([]*todo.TodoItem, error) {
			return []*todo.TodoItem{{Description: "Benchmark todo"}}, nil
		},
	}
	messaging := &benchMockMessaging{}
	cache := &benchMockCache{}
	service := todo.NewTodoService(todoRepo, messaging, cache)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.ListTodos(context.Background(), 10, 0)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}
