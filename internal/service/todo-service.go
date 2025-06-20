package service

import (
	"context"
	"fmt"
	"log/slog"
	"taskflow/internal/domain/entities"
	"taskflow/internal/domain/repositories"
	"time"

	"github.com/google/uuid"
)

type TodoService interface {
	CreateTodo(ctx context.Context, req *entities.CreateTodoRequest) (*entities.TodoItem, error)
	GetTodo(ctx context.Context, id uuid.UUID) (*entities.TodoItem, error)
	ListTodos(ctx context.Context, limit, offset int) ([]*entities.TodoItem, error)
	UpdateTodo(ctx context.Context, id uuid.UUID, req *entities.UpdateTodoRequest) (*entities.TodoItem, error)
	DeleteTodo(ctx context.Context, id uuid.UUID) error
}

type todoService struct {
	todoRepo   repositories.TodoRepository
	streamRepo repositories.StreamRepository
	logger     *slog.Logger
}

func NewTodoService(todoRepo repositories.TodoRepository, streamRepo repositories.StreamRepository) TodoService {
	return &todoService{
		todoRepo:   todoRepo,
		streamRepo: streamRepo,
		logger:     slog.Default(),
	}
}

func (uc *todoService) CreateTodo(ctx context.Context, req *entities.CreateTodoRequest) (*entities.TodoItem, error) {
	if req.Description == "" {
		return nil, fmt.Errorf("description is required")
	}

	if req.DueDate.Before(time.Now()) {
		return nil, fmt.Errorf("due date must be in the future")
	}

	todo := &entities.TodoItem{
		ID:          uuid.New(),
		Description: req.Description,
		DueDate:     req.DueDate,
		FileID:      req.FileID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.todoRepo.Create(ctx, todo); err != nil {
		uc.logger.Error("failed to create todo", "error", err, "todo_id", todo.ID)
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}

	uc.logger.Info("todo created successfully", "todo_id", todo.ID, "description", todo.Description)

	// Publish to stream
	go func() {
		if err := uc.streamRepo.PublishTodoCreated(context.Background(), todo); err != nil {
			uc.logger.Error("failed to publish to stream", "error", err, "todo_id", todo.ID)
		} else {
			uc.logger.Info("todo event published to stream", "todo_id", todo.ID)
		}
	}()

	return todo, nil
}

func (uc *todoService) GetTodo(ctx context.Context, id uuid.UUID) (*entities.TodoItem, error) {
	todo, err := uc.todoRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("failed to get todo", "error", err, "todo_id", id)
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	uc.logger.Info("todo retrieved", "todo_id", id)
	return todo, nil
}

func (uc *todoService) ListTodos(ctx context.Context, limit, offset int) ([]*entities.TodoItem, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	todos, err := uc.todoRepo.List(ctx, limit, offset)
	if err != nil {
		uc.logger.Error("failed to list todos", "error", err, "limit", limit, "offset", offset)
		return nil, fmt.Errorf("failed to list todos: %w", err)
	}

	uc.logger.Info("todos listed", "count", len(todos), "limit", limit, "offset", offset)
	return todos, nil
}

func (uc *todoService) UpdateTodo(ctx context.Context, id uuid.UUID, req *entities.UpdateTodoRequest) (*entities.TodoItem, error) {
	// Get existing todo
	existing, err := uc.todoRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("todo not found: %w", err)
	}

	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.DueDate != nil {
		existing.DueDate = *req.DueDate
	}
	if req.FileID != nil {
		existing.FileID = req.FileID
	}

	existing.UpdatedAt = time.Now()

	if err := uc.todoRepo.Update(ctx, existing); err != nil {
		uc.logger.Error("failed to update todo", "error", err, "todo_id", id)
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}

	uc.logger.Info("todo updated", "todo_id", id)
	return existing, nil
}

func (uc *todoService) DeleteTodo(ctx context.Context, id uuid.UUID) error {
	// Check if todo exists
	if _, err := uc.todoRepo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("todo not found: %w", err)
	}

	if err := uc.todoRepo.Delete(ctx, id); err != nil {
		uc.logger.Error("failed to delete todo", "error", err, "todo_id", id)
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	uc.logger.Info("todo deleted", "todo_id", id)
	return nil
}
