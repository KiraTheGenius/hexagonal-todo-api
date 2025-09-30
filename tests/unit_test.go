package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"taskflow/internal/domain/todo"

	"github.com/google/uuid"
)

// --- Mock Repositories ---
type mockTodoRepo struct {
	CreateFn  func(ctx context.Context, todoItem *todo.TodoItem) error
	GetByIDFn func(ctx context.Context, id uuid.UUID) (*todo.TodoItem, error)
	ListFn    func(ctx context.Context, limit, offset int) ([]*todo.TodoItem, error)
	UpdateFn  func(ctx context.Context, todoItem *todo.TodoItem) error
	DeleteFn  func(ctx context.Context, id uuid.UUID) error
}

func (m *mockTodoRepo) Create(ctx context.Context, todoItem *todo.TodoItem) error {
	return m.CreateFn(ctx, todoItem)
}
func (m *mockTodoRepo) GetByID(ctx context.Context, id uuid.UUID) (*todo.TodoItem, error) {
	return m.GetByIDFn(ctx, id)
}
func (m *mockTodoRepo) List(ctx context.Context, limit, offset int) ([]*todo.TodoItem, error) {
	return m.ListFn(ctx, limit, offset)
}
func (m *mockTodoRepo) Update(ctx context.Context, todoItem *todo.TodoItem) error {
	return m.UpdateFn(ctx, todoItem)
}
func (m *mockTodoRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.DeleteFn(ctx, id)
}

type mockMessaging struct {
	PublishFn func(ctx context.Context, topic string, message interface{}) error
}

func (m *mockMessaging) Publish(ctx context.Context, topic string, message interface{}) error {
	if m.PublishFn != nil {
		return m.PublishFn(ctx, topic, message)
	}
	return nil
}

func (m *mockMessaging) PublishWithKey(ctx context.Context, topic string, key string, message interface{}) error {
	return m.Publish(ctx, topic, message)
}

type mockCache struct{}

func (m *mockCache) Get(ctx context.Context, key string) (string, error)              { return "", nil }
func (m *mockCache) Set(ctx context.Context, key string, value string, ttl int) error { return nil }
func (m *mockCache) Delete(ctx context.Context, key string) error                     { return nil }
func (m *mockCache) Exists(ctx context.Context, key string) (bool, error)             { return false, nil }

// --- Tests ---
func TestCreateTodo_Success(t *testing.T) {
	todoRepo := &mockTodoRepo{
		CreateFn: func(ctx context.Context, todoItem *todo.TodoItem) error { return nil },
	}
	messaging := &mockMessaging{}
	cache := &mockCache{}
	service := todo.NewTodoService(todoRepo, messaging, cache)

	req := &todo.CreateTodoRequest{
		Description: "Test todo",
		DueDate:     time.Now().Add(24 * time.Hour),
	}
	todoItem, err := service.CreateTodo(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if todoItem.Description != req.Description {
		t.Errorf("expected description %q, got %q", req.Description, todoItem.Description)
	}
}

func TestCreateTodo_ValidationError(t *testing.T) {
	todoRepo := &mockTodoRepo{}
	messaging := &mockMessaging{}
	cache := &mockCache{}
	service := todo.NewTodoService(todoRepo, messaging, cache)

	req := &todo.CreateTodoRequest{Description: "", DueDate: time.Now().Add(24 * time.Hour)}
	_, err := service.CreateTodo(context.Background(), req)
	if err == nil || err.Error() != "description is required" {
		t.Errorf("expected validation error, got %v", err)
	}

	req = &todo.CreateTodoRequest{Description: "desc", DueDate: time.Now().Add(-24 * time.Hour)}
	_, err = service.CreateTodo(context.Background(), req)
	if err == nil || err.Error() != "due date must be in the future" {
		t.Errorf("expected due date validation error, got %v", err)
	}
}

func TestCreateTodo_RepoError(t *testing.T) {
	todoRepo := &mockTodoRepo{
		CreateFn: func(ctx context.Context, todoItem *todo.TodoItem) error { return errors.New("db error") },
	}
	messaging := &mockMessaging{}
	cache := &mockCache{}
	service := todo.NewTodoService(todoRepo, messaging, cache)

	req := &todo.CreateTodoRequest{Description: "desc", DueDate: time.Now().Add(24 * time.Hour)}
	_, err := service.CreateTodo(context.Background(), req)
	if err == nil || err.Error() == "" {
		t.Errorf("expected repo error, got %v", err)
	}
}

func TestGetTodo_Success(t *testing.T) {
	id := uuid.New()
	todoRepo := &mockTodoRepo{
		GetByIDFn: func(ctx context.Context, tid uuid.UUID) (*todo.TodoItem, error) {
			return &todo.TodoItem{ID: tid, Description: "desc"}, nil
		},
	}
	messaging := &mockMessaging{}
	cache := &mockCache{}
	service := todo.NewTodoService(todoRepo, messaging, cache)

	todoItem, err := service.GetTodo(context.Background(), id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if todoItem.ID != id {
		t.Errorf("expected id %v, got %v", id, todoItem.ID)
	}
}

func TestGetTodo_RepoError(t *testing.T) {
	todoRepo := &mockTodoRepo{
		GetByIDFn: func(ctx context.Context, tid uuid.UUID) (*todo.TodoItem, error) {
			return nil, errors.New("not found")
		},
	}
	messaging := &mockMessaging{}
	cache := &mockCache{}
	service := todo.NewTodoService(todoRepo, messaging, cache)

	_, err := service.GetTodo(context.Background(), uuid.New())
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestListTodos_Success(t *testing.T) {
	todoRepo := &mockTodoRepo{
		ListFn: func(ctx context.Context, limit, offset int) ([]*todo.TodoItem, error) {
			return []*todo.TodoItem{{ID: uuid.New(), Description: "desc"}}, nil
		},
	}
	messaging := &mockMessaging{}
	cache := &mockCache{}
	service := todo.NewTodoService(todoRepo, messaging, cache)

	todos, err := service.ListTodos(context.Background(), 10, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(todos) != 1 {
		t.Errorf("expected 1 todo, got %d", len(todos))
	}
}

func TestListTodos_RepoError(t *testing.T) {
	todoRepo := &mockTodoRepo{
		ListFn: func(ctx context.Context, limit, offset int) ([]*todo.TodoItem, error) {
			return nil, errors.New("db error")
		},
	}
	messaging := &mockMessaging{}
	cache := &mockCache{}
	service := todo.NewTodoService(todoRepo, messaging, cache)

	_, err := service.ListTodos(context.Background(), 10, 0)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestUpdateTodo_Success(t *testing.T) {
	id := uuid.New()
	desc := "updated"
	todoRepo := &mockTodoRepo{
		GetByIDFn: func(ctx context.Context, tid uuid.UUID) (*todo.TodoItem, error) {
			return &todo.TodoItem{ID: tid, Description: "old"}, nil
		},
		UpdateFn: func(ctx context.Context, todoItem *todo.TodoItem) error { return nil },
	}
	messaging := &mockMessaging{}
	cache := &mockCache{}
	service := todo.NewTodoService(todoRepo, messaging, cache)

	req := &todo.UpdateTodoRequest{Description: &desc}
	todoItem, err := service.UpdateTodo(context.Background(), id, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if todoItem.Description != desc {
		t.Errorf("expected description %q, got %q", desc, todoItem.Description)
	}
}

func TestUpdateTodo_NotFound(t *testing.T) {
	todoRepo := &mockTodoRepo{
		GetByIDFn: func(ctx context.Context, tid uuid.UUID) (*todo.TodoItem, error) {
			return nil, errors.New("not found")
		},
	}
	messaging := &mockMessaging{}
	cache := &mockCache{}
	service := todo.NewTodoService(todoRepo, messaging, cache)

	req := &todo.UpdateTodoRequest{Description: new(string)}
	_, err := service.UpdateTodo(context.Background(), uuid.New(), req)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestUpdateTodo_RepoError(t *testing.T) {
	id := uuid.New()
	desc := "desc"
	todoRepo := &mockTodoRepo{
		GetByIDFn: func(ctx context.Context, tid uuid.UUID) (*todo.TodoItem, error) {
			return &todo.TodoItem{ID: tid, Description: "old"}, nil
		},
		UpdateFn: func(ctx context.Context, todoItem *todo.TodoItem) error { return errors.New("db error") },
	}
	messaging := &mockMessaging{}
	cache := &mockCache{}
	service := todo.NewTodoService(todoRepo, messaging, cache)

	req := &todo.UpdateTodoRequest{Description: &desc}
	_, err := service.UpdateTodo(context.Background(), id, req)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestDeleteTodo_Success(t *testing.T) {
	id := uuid.New()
	todoRepo := &mockTodoRepo{
		GetByIDFn: func(ctx context.Context, tid uuid.UUID) (*todo.TodoItem, error) {
			return &todo.TodoItem{ID: tid}, nil
		},
		DeleteFn: func(ctx context.Context, tid uuid.UUID) error { return nil },
	}
	messaging := &mockMessaging{}
	cache := &mockCache{}
	service := todo.NewTodoService(todoRepo, messaging, cache)

	err := service.DeleteTodo(context.Background(), id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDeleteTodo_NotFound(t *testing.T) {
	todoRepo := &mockTodoRepo{
		GetByIDFn: func(ctx context.Context, tid uuid.UUID) (*todo.TodoItem, error) {
			return nil, errors.New("not found")
		},
	}
	messaging := &mockMessaging{}
	cache := &mockCache{}
	service := todo.NewTodoService(todoRepo, messaging, cache)

	err := service.DeleteTodo(context.Background(), uuid.New())
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestDeleteTodo_RepoError(t *testing.T) {
	id := uuid.New()
	todoRepo := &mockTodoRepo{
		GetByIDFn: func(ctx context.Context, tid uuid.UUID) (*todo.TodoItem, error) {
			return &todo.TodoItem{ID: tid}, nil
		},
		DeleteFn: func(ctx context.Context, tid uuid.UUID) error { return errors.New("db error") },
	}
	messaging := &mockMessaging{}
	cache := &mockCache{}
	service := todo.NewTodoService(todoRepo, messaging, cache)

	err := service.DeleteTodo(context.Background(), id)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
