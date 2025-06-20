package entities

import (
	"time"

	"github.com/google/uuid"
)

type TodoItem struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Description string    `json:"description" db:"description"`
	DueDate     time.Time `json:"dueDate" db:"due_date"`
	FileID      *string   `json:"fileId,omitempty" db:"file_id"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type CreateTodoRequest struct {
	Description string    `json:"description" binding:"required"`
	DueDate     time.Time `json:"dueDate" binding:"required"`
	FileID      *string   `json:"fileId,omitempty"`
}

type UpdateTodoRequest struct {
	Description *string    `json:"description,omitempty"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
	FileID      *string    `json:"fileId,omitempty"`
}
