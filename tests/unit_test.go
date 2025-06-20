package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"taskflow/internal/domain/entities"
	"taskflow/internal/service"

	"github.com/google/uuid"
)

// --- Mock Repositories ---
type mockTodoRepo struct {
	CreateFn  func(ctx context.Context, todo *entities.TodoItem) error
	GetByIDFn func(ctx context.Context, id uuid.UUID) (*entities.TodoItem, error)
	ListFn    func(ctx context.Context, limit, offset int) ([]*entities.TodoItem, error)
	UpdateFn  func(ctx context.Context, todo *entities.TodoItem) error
	DeleteFn  func(ctx context.Context, id uuid.UUID) error
}

func (m *mockTodoRepo) Create(ctx context.Context, todo *entities.TodoItem) error {
	return m.CreateFn(ctx, todo)
}
func (m *mockTodoRepo) GetByID(ctx context.Context, id uuid.UUID) (*entities.TodoItem, error) {
	return m.GetByIDFn(ctx, id)
}
func (m *mockTodoRepo) List(ctx context.Context, limit, offset int) ([]*entities.TodoItem, error) {
	return m.ListFn(ctx, limit, offset)
}
func (m *mockTodoRepo) Update(ctx context.Context, todo *entities.TodoItem) error {
	return m.UpdateFn(ctx, todo)
}
func (m *mockTodoRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.DeleteFn(ctx, id)
}

type mockStreamRepo struct {
	PublishTodoCreatedFn func(ctx context.Context, todo *entities.TodoItem) error
}

func (m *mockStreamRepo) PublishTodoCreated(ctx context.Context, todo *entities.TodoItem) error {
	if m.PublishTodoCreatedFn != nil {
		return m.PublishTodoCreatedFn(ctx, todo)
	}
	return nil
}

// --- Tests ---
func TestCreateTodo_Success(t *testing.T) {
	todoRepo := &mockTodoRepo{
		CreateFn: func(ctx context.Context, todo *entities.TodoItem) error { return nil },
	}
	streamRepo := &mockStreamRepo{}
	uc := service.NewTodoService(todoRepo, streamRepo)

	req := &entities.CreateTodoRequest{
		Description: "Test todo",
		DueDate:     time.Now().Add(24 * time.Hour),
	}
	todo, err := uc.CreateTodo(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if todo.Description != req.Description {
		t.Errorf("expected description %q, got %q", req.Description, todo.Description)
	}
}

func TestCreateTodo_ValidationError(t *testing.T) {
	todoRepo := &mockTodoRepo{}
	streamRepo := &mockStreamRepo{}
	uc := service.NewTodoService(todoRepo, streamRepo)

	req := &entities.CreateTodoRequest{Description: "", DueDate: time.Now().Add(24 * time.Hour)}
	_, err := uc.CreateTodo(context.Background(), req)
	if err == nil || err.Error() != "description is required" {
		t.Errorf("expected validation error, got %v", err)
	}

	req = &entities.CreateTodoRequest{Description: "desc", DueDate: time.Now().Add(-24 * time.Hour)}
	_, err = uc.CreateTodo(context.Background(), req)
	if err == nil || err.Error() != "due date must be in the future" {
		t.Errorf("expected due date validation error, got %v", err)
	}
}

func TestCreateTodo_RepoError(t *testing.T) {
	todoRepo := &mockTodoRepo{
		CreateFn: func(ctx context.Context, todo *entities.TodoItem) error { return errors.New("db error") },
	}
	streamRepo := &mockStreamRepo{}
	uc := service.NewTodoService(todoRepo, streamRepo)

	req := &entities.CreateTodoRequest{Description: "desc", DueDate: time.Now().Add(24 * time.Hour)}
	_, err := uc.CreateTodo(context.Background(), req)
	if err == nil || err.Error() == "" {
		t.Errorf("expected repo error, got %v", err)
	}
}

func TestGetTodo_Success(t *testing.T) {
	id := uuid.New()
	todoRepo := &mockTodoRepo{
		GetByIDFn: func(ctx context.Context, tid uuid.UUID) (*entities.TodoItem, error) {
			return &entities.TodoItem{ID: tid, Description: "desc"}, nil
		},
	}
	streamRepo := &mockStreamRepo{}
	uc := service.NewTodoService(todoRepo, streamRepo)

	todo, err := uc.GetTodo(context.Background(), id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if todo.ID != id {
		t.Errorf("expected id %v, got %v", id, todo.ID)
	}
}

func TestGetTodo_RepoError(t *testing.T) {
	todoRepo := &mockTodoRepo{
		GetByIDFn: func(ctx context.Context, tid uuid.UUID) (*entities.TodoItem, error) {
			return nil, errors.New("not found")
		},
	}
	streamRepo := &mockStreamRepo{}
	uc := service.NewTodoService(todoRepo, streamRepo)

	_, err := uc.GetTodo(context.Background(), uuid.New())
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestListTodos_Success(t *testing.T) {
	todoRepo := &mockTodoRepo{
		ListFn: func(ctx context.Context, limit, offset int) ([]*entities.TodoItem, error) {
			return []*entities.TodoItem{{ID: uuid.New(), Description: "desc"}}, nil
		},
	}
	streamRepo := &mockStreamRepo{}
	uc := service.NewTodoService(todoRepo, streamRepo)

	todos, err := uc.ListTodos(context.Background(), 10, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(todos) != 1 {
		t.Errorf("expected 1 todo, got %d", len(todos))
	}
}

func TestListTodos_RepoError(t *testing.T) {
	todoRepo := &mockTodoRepo{
		ListFn: func(ctx context.Context, limit, offset int) ([]*entities.TodoItem, error) {
			return nil, errors.New("db error")
		},
	}
	streamRepo := &mockStreamRepo{}
	uc := service.NewTodoService(todoRepo, streamRepo)

	_, err := uc.ListTodos(context.Background(), 10, 0)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestUpdateTodo_Success(t *testing.T) {
	id := uuid.New()
	desc := "updated"
	todoRepo := &mockTodoRepo{
		GetByIDFn: func(ctx context.Context, tid uuid.UUID) (*entities.TodoItem, error) {
			return &entities.TodoItem{ID: tid, Description: "old"}, nil
		},
		UpdateFn: func(ctx context.Context, todo *entities.TodoItem) error { return nil },
	}
	streamRepo := &mockStreamRepo{}
	uc := service.NewTodoService(todoRepo, streamRepo)

	req := &entities.UpdateTodoRequest{Description: &desc}
	todo, err := uc.UpdateTodo(context.Background(), id, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if todo.Description != desc {
		t.Errorf("expected description %q, got %q", desc, todo.Description)
	}
}

func TestUpdateTodo_NotFound(t *testing.T) {
	todoRepo := &mockTodoRepo{
		GetByIDFn: func(ctx context.Context, tid uuid.UUID) (*entities.TodoItem, error) {
			return nil, errors.New("not found")
		},
	}
	streamRepo := &mockStreamRepo{}
	uc := service.NewTodoService(todoRepo, streamRepo)

	req := &entities.UpdateTodoRequest{Description: new(string)}
	_, err := uc.UpdateTodo(context.Background(), uuid.New(), req)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestUpdateTodo_RepoError(t *testing.T) {
	id := uuid.New()
	desc := "desc"
	todoRepo := &mockTodoRepo{
		GetByIDFn: func(ctx context.Context, tid uuid.UUID) (*entities.TodoItem, error) {
			return &entities.TodoItem{ID: tid, Description: "old"}, nil
		},
		UpdateFn: func(ctx context.Context, todo *entities.TodoItem) error { return errors.New("db error") },
	}
	streamRepo := &mockStreamRepo{}
	uc := service.NewTodoService(todoRepo, streamRepo)

	req := &entities.UpdateTodoRequest{Description: &desc}
	_, err := uc.UpdateTodo(context.Background(), id, req)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestDeleteTodo_Success(t *testing.T) {
	id := uuid.New()
	todoRepo := &mockTodoRepo{
		GetByIDFn: func(ctx context.Context, tid uuid.UUID) (*entities.TodoItem, error) {
			return &entities.TodoItem{ID: tid}, nil
		},
		DeleteFn: func(ctx context.Context, tid uuid.UUID) error { return nil },
	}
	streamRepo := &mockStreamRepo{}
	uc := service.NewTodoService(todoRepo, streamRepo)

	err := uc.DeleteTodo(context.Background(), id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDeleteTodo_NotFound(t *testing.T) {
	todoRepo := &mockTodoRepo{
		GetByIDFn: func(ctx context.Context, tid uuid.UUID) (*entities.TodoItem, error) {
			return nil, errors.New("not found")
		},
	}
	streamRepo := &mockStreamRepo{}
	uc := service.NewTodoService(todoRepo, streamRepo)

	err := uc.DeleteTodo(context.Background(), uuid.New())
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestDeleteTodo_RepoError(t *testing.T) {
	id := uuid.New()
	todoRepo := &mockTodoRepo{
		GetByIDFn: func(ctx context.Context, tid uuid.UUID) (*entities.TodoItem, error) {
			return &entities.TodoItem{ID: tid}, nil
		},
		DeleteFn: func(ctx context.Context, tid uuid.UUID) error { return errors.New("db error") },
	}
	streamRepo := &mockStreamRepo{}
	uc := service.NewTodoService(todoRepo, streamRepo)

	err := uc.DeleteTodo(context.Background(), id)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
