package file

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
)

type fileService struct {
	fileRepo Repository
	storage  Storage
}

func NewFileService(fileRepo Repository, storage Storage) FileService {
	return &fileService{
		fileRepo: fileRepo,
		storage:  storage,
	}
}

func (s *fileService) UploadFile(ctx context.Context, req *CreateFileRequest, content io.Reader) (*UploadResponse, error) {
	// Upload to storage
	storageKey, err := s.storage.Upload(ctx, req.Filename, content, req.ContentType)
	if err != nil {
		return nil, err
	}

	// Create file entity
	file := &File{
		ID:          uuid.New(),
		Filename:    req.Filename,
		ContentType: req.ContentType,
		Size:        req.Size,
		StorageKey:  storageKey,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save to repository
	if err := s.fileRepo.Create(ctx, file); err != nil {
		// Clean up storage if repository save fails
		s.storage.Delete(ctx, storageKey)
		return nil, err
	}

	// Get URL if available
	url, _ := s.storage.GetURL(ctx, storageKey)

	return &UploadResponse{
		FileID: file.ID.String(),
		URL:    url,
	}, nil
}

func (s *fileService) GetFile(ctx context.Context, fileID string) (*File, error) {
	return s.fileRepo.GetByID(ctx, fileID)
}

func (s *fileService) DownloadFile(ctx context.Context, fileID string) (io.ReadCloser, error) {
	// Get file metadata
	file, err := s.fileRepo.GetByID(ctx, fileID)
	if err != nil {
		return nil, err
	}

	// Download from storage
	return s.storage.Download(ctx, file.StorageKey)
}

func (s *fileService) DeleteFile(ctx context.Context, fileID string) error {
	// Get file metadata
	file, err := s.fileRepo.GetByID(ctx, fileID)
	if err != nil {
		return err
	}

	// Delete from storage
	if err := s.storage.Delete(ctx, file.StorageKey); err != nil {
		return err
	}

	// Delete from repository
	return s.fileRepo.Delete(ctx, fileID)
}

func (s *fileService) ListFiles(ctx context.Context, limit, offset int) ([]*File, error) {
	return s.fileRepo.List(ctx, limit, offset)
}

func (s *fileService) UpdateFile(ctx context.Context, fileID string, req *UpdateFileRequest) (*File, error) {
	// Get existing file
	file, err := s.fileRepo.GetByID(ctx, fileID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Filename != nil {
		file.Filename = *req.Filename
	}
	if req.ContentType != nil {
		file.ContentType = *req.ContentType
	}
	file.UpdatedAt = time.Now()

	// Save changes
	if err := s.fileRepo.Update(ctx, file); err != nil {
		return nil, err
	}

	return file, nil
}
