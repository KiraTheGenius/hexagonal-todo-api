package shared

import (
	"context"
	"io"
)

// Cache defines the interface for caching operations
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl int) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}

// Storage defines the interface for file storage operations
type Storage interface {
	Upload(ctx context.Context, filename string, content io.Reader, contentType string) (string, error)
	Download(ctx context.Context, fileID string) (io.ReadCloser, error)
	Delete(ctx context.Context, fileID string) error
	GetURL(ctx context.Context, fileID string) (string, error)
}

// Messaging defines the interface for message publishing
type Messaging interface {
	Publish(ctx context.Context, topic string, message interface{}) error
	PublishWithKey(ctx context.Context, topic string, key string, message interface{}) error
}

// EventStore defines the interface for event storage
type EventStore interface {
	Append(ctx context.Context, streamID string, events []Event) error
	GetEvents(ctx context.Context, streamID string, fromVersion int) ([]Event, error)
}

// Event represents a domain event
type Event struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Version   int         `json:"version"`
	Timestamp int64       `json:"timestamp"`
}
