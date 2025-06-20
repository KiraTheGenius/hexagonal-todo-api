package repositories

import (
	"context"
	"io"
)

type FileRepository interface {
	Upload(ctx context.Context, filename string, content io.Reader, contentType string) (string, error)
	Download(ctx context.Context, fileID string) (io.ReadCloser, error)
	Delete(ctx context.Context, fileID string) error
}
