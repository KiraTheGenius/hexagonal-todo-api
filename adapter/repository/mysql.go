package repository

import (
	"context"
	"taskflow/internal/domain/file"
	"taskflow/internal/domain/todo"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewGormConnection(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func RunGormMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&todo.TodoItem{},
		&file.File{},
	)
}

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) todo.Repository {
	return &todoRepository{db: db}
}

func (r *todoRepository) Create(ctx context.Context, todoItem *todo.TodoItem) error {
	return r.db.WithContext(ctx).Create(todoItem).Error
}

func (r *todoRepository) GetByID(ctx context.Context, id uuid.UUID) (*todo.TodoItem, error) {
	var todoItem todo.TodoItem
	err := r.db.WithContext(ctx).First(&todoItem, "id = ?", id.String()).Error
	if err != nil {
		return nil, err
	}
	return &todoItem, nil
}

func (r *todoRepository) List(ctx context.Context, limit, offset int) ([]*todo.TodoItem, error) {
	var todos []*todo.TodoItem
	err := r.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Offset(offset).Find(&todos).Error
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func (r *todoRepository) Update(ctx context.Context, todoItem *todo.TodoItem) error {
	return r.db.WithContext(ctx).Model(todoItem).Updates(todoItem).Error
}

func (r *todoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&todo.TodoItem{}, "id = ?", id.String()).Error
}
