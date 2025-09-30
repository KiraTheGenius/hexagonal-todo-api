package file

import (
	"context"
	"taskflow/internal/domain/shared"
)

// Repository defines the file repository interface
type Repository interface {
	Create(ctx context.Context, file *File) error
	GetByID(ctx context.Context, id string) (*File, error)
	Update(ctx context.Context, file *File) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*File, error)
}

// Storage defines the file storage interface (uses shared storage port)
type Storage = shared.Storage
