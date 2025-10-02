package shared

import "errors"

// Domain errors
var (
	ErrNotFound     = errors.New("resource not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrConflict     = errors.New("resource conflict")
	ErrInternal     = errors.New("internal error")
	ErrValidation   = errors.New("validation failed")
	ErrTimeout      = errors.New("operation timeout")
)

// DomainError represents a domain-specific error
type DomainError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *DomainError) Error() string {
	return e.Message
}

// NewDomainError creates a new domain error
func NewDomainError(code, message, details string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Common error codes
const (
	ErrCodeNotFound     = "NOT_FOUND"
	ErrCodeInvalidInput = "INVALID_INPUT"
	ErrCodeUnauthorized = "UNAUTHORIZED"
	ErrCodeForbidden    = "FORBIDDEN"
	ErrCodeConflict     = "CONFLICT"
	ErrCodeInternal     = "INTERNAL_ERROR"
	ErrCodeValidation   = "VALIDATION_FAILED"
	ErrCodeTimeout      = "TIMEOUT"
)

// Helper functions for common errors
func NewValidationError(message string) *DomainError {
	return NewDomainError(ErrCodeValidation, message, "")
}

func NewNotFoundError(message string) *DomainError {
	return NewDomainError(ErrCodeNotFound, message, "")
}
