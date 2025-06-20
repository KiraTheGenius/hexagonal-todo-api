package database

import (
	"context"
	"taskflow/internal/domain/entities"
	"taskflow/internal/domain/repositories"

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
	return db.AutoMigrate(&entities.TodoItem{})
}

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) repositories.TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) Create(ctx context.Context, todo *entities.TodoItem) error {
	return r.db.WithContext(ctx).Create(todo).Error
}

func (r *todoRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.TodoItem, error) {
	var todo entities.TodoItem
	err := r.db.WithContext(ctx).First(&todo, "id = ?", id.String()).Error
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *todoRepository) List(ctx context.Context, limit, offset int) ([]*entities.TodoItem, error) {
	var todos []*entities.TodoItem
	err := r.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Offset(offset).Find(&todos).Error
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func (r *todoRepository) Update(ctx context.Context, todo *entities.TodoItem) error {
	return r.db.WithContext(ctx).Model(todo).Updates(todo).Error
}

func (r *todoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.TodoItem{}, "id = ?", id.String()).Error
}
