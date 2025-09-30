package file

import (
	"time"

	"github.com/google/uuid"
)

// File represents a file entity in the domain
type File struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Filename    string    `json:"filename" db:"filename"`
	ContentType string    `json:"contentType" db:"content_type"`
	Size        int64     `json:"size" db:"size"`
	StorageKey  string    `json:"storageKey" db:"storage_key"`
	URL         string    `json:"url,omitempty" db:"url"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

// CreateFileRequest represents the request to create a file
type CreateFileRequest struct {
	Filename    string `json:"filename" binding:"required"`
	ContentType string `json:"contentType" binding:"required"`
	Size        int64  `json:"size" binding:"required,min=1"`
}

// UpdateFileRequest represents the request to update a file
type UpdateFileRequest struct {
	Filename    *string `json:"filename,omitempty"`
	ContentType *string `json:"contentType,omitempty"`
}

// FileMetadata represents file metadata without content
type FileMetadata struct {
	ID          uuid.UUID `json:"id"`
	Filename    string    `json:"filename"`
	ContentType string    `json:"contentType"`
	Size        int64     `json:"size"`
	URL         string    `json:"url,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// UploadResponse represents the response after file upload
type UploadResponse struct {
	FileID string `json:"fileId"`
	URL    string `json:"url,omitempty"`
}
