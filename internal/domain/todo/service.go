package todo

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"taskflow/internal/domain/shared"

	"github.com/google/uuid"
)

type todoService struct {
	todoRepo  Repository
	messaging Messaging
	cache     Cache
	logger    *slog.Logger
}

func NewTodoService(todoRepo Repository, messaging Messaging, cache Cache) TodoService {
	return &todoService{
		todoRepo:  todoRepo,
		messaging: messaging,
		cache:     cache,
		logger:    slog.Default(),
	}
}

func (s *todoService) CreateTodo(ctx context.Context, req *CreateTodoRequest) (*TodoItem, error) {
	if req.Description == "" {
		return nil, shared.NewValidationError("description is required")
	}

	if req.DueDate.Before(time.Now()) {
		return nil, shared.NewValidationError("due date must be in the future")
	}

	todo := &TodoItem{
		ID:          uuid.New(),
		Description: req.Description,
		DueDate:     req.DueDate,
		FileID:      req.FileID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.todoRepo.Create(ctx, todo); err != nil {
		s.logger.Error("failed to create todo", "error", err, "todo_id", todo.ID)
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}

	s.logger.Info("todo created successfully", "todo_id", todo.ID, "description", todo.Description)

	// Publish to messaging system
	go func() {
		if err := s.messaging.Publish(context.Background(), "todo.created", todo); err != nil {
			s.logger.Error("failed to publish todo created event", "error", err, "todo_id", todo.ID)
		} else {
			s.logger.Info("todo created event published", "todo_id", todo.ID)
		}
	}()

	return todo, nil
}

func (s *todoService) GetTodo(ctx context.Context, id uuid.UUID) (*TodoItem, error) {
	todo, err := s.todoRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get todo", "error", err, "todo_id", id)
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	s.logger.Info("todo retrieved", "todo_id", id)
	return todo, nil
}

func (s *todoService) ListTodos(ctx context.Context, limit, offset int) ([]*TodoItem, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	todos, err := s.todoRepo.List(ctx, limit, offset)
	if err != nil {
		s.logger.Error("failed to list todos", "error", err, "limit", limit, "offset", offset)
		return nil, fmt.Errorf("failed to list todos: %w", err)
	}

	s.logger.Info("todos listed", "count", len(todos), "limit", limit, "offset", offset)
	return todos, nil
}

func (s *todoService) UpdateTodo(ctx context.Context, id uuid.UUID, req *UpdateTodoRequest) (*TodoItem, error) {
	// Get existing todo
	existing, err := s.todoRepo.GetByID(ctx, id)
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

	if err := s.todoRepo.Update(ctx, existing); err != nil {
		s.logger.Error("failed to update todo", "error", err, "todo_id", id)
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}

	s.logger.Info("todo updated", "todo_id", id)
	return existing, nil
}

func (s *todoService) DeleteTodo(ctx context.Context, id uuid.UUID) error {
	// Check if todo exists
	if _, err := s.todoRepo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("todo not found: %w", err)
	}

	if err := s.todoRepo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete todo", "error", err, "todo_id", id)
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	s.logger.Info("todo deleted", "todo_id", id)
	return nil
}
