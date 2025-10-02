package file

import (
	"context"
	"io"
	"taskflow/internal/domain/shared"
)

// FileService defines the file service interface
type FileService interface {
	UploadFile(ctx context.Context, req *CreateFileRequest, content io.Reader) (*UploadResponse, error)
	GetFile(ctx context.Context, fileID string) (*File, error)
	DownloadFile(ctx context.Context, fileID string) (io.ReadCloser, error)
	DeleteFile(ctx context.Context, fileID string) error
	ListFiles(ctx context.Context, limit, offset int) ([]*File, error)
	UpdateFile(ctx context.Context, fileID string, req *UpdateFileRequest) (*File, error)
}

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
