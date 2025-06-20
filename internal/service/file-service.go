package service

import (
	"context"
	"io"
	"taskflow/internal/domain/repositories"
)

type FileService interface {
	UploadFile(ctx context.Context, filename string, content io.Reader, contentType string) (string, error)
}

type fileService struct {
	fileRepo repositories.FileRepository
}

func NewFileService(fileRepo repositories.FileRepository) FileService {
	return &fileService{
		fileRepo: fileRepo,
	}
}

func (uc *fileService) UploadFile(ctx context.Context, filename string, content io.Reader, contentType string) (string, error) {
	return uc.fileRepo.Upload(ctx, filename, content, contentType)
}
