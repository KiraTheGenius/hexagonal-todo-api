package repository

import (
	"context"
	"taskflow/internal/domain/file"

	"gorm.io/gorm"
)

type fileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) file.Repository {
	return &fileRepository{db: db}
}

func (r *fileRepository) Create(ctx context.Context, file *file.File) error {
	return r.db.WithContext(ctx).Create(file).Error
}

func (r *fileRepository) GetByID(ctx context.Context, id string) (*file.File, error) {
	var fileItem file.File
	err := r.db.WithContext(ctx).First(&fileItem, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &fileItem, nil
}

func (r *fileRepository) Update(ctx context.Context, file *file.File) error {
	return r.db.WithContext(ctx).Model(file).Updates(file).Error
}

func (r *fileRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&file.File{}, "id = ?", id).Error
}

func (r *fileRepository) List(ctx context.Context, limit, offset int) ([]*file.File, error) {
	var files []*file.File
	err := r.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Offset(offset).Find(&files).Error
	if err != nil {
		return nil, err
	}
	return files, nil
}
