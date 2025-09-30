# TaskFlow Architecture Documentation

## Overview

TaskFlow is built using **Hexagonal Architecture** (also known as Ports & Adapters), a design pattern that promotes clean separation between business logic and external concerns. This architecture makes the application more testable, maintainable, and flexible.

## Architecture Principles

### 1. Domain-Centric Design
The domain layer (`internal/domain/`) contains the core business logic and is completely independent of external frameworks, databases, or UI technologies.

### 2. Dependency Inversion
Dependencies point inward toward the domain layer. External systems are accessed through interfaces (ports) defined by the domain.

### 3. Separation of Concerns
Each layer has a single responsibility:
- **Domain**: Business logic and rules
- **Adapters**: External system integrations
- **Interfaces**: User/system interactions

## Architecture Layers

### Domain Layer (Hexagon Core)

Located in `internal/domain/`, this is the heart of the application.

#### Structure
```
internal/domain/
├── todo/           # Todo domain
│   ├── model.go    # Domain models and DTOs
│   ├── port.go     # Port interfaces
│   └── service.go  # Domain services
└── file/           # File domain
    ├── port.go     # Port interfaces
    └── service.go  # Domain services
```

#### Components

**Models** (`model.go`):
- `TodoItem`: Core business entity
- `CreateTodoRequest`: Input DTO
- `UpdateTodoRequest`: Update DTO

**Ports** (`port.go`):
- Define interfaces that the domain needs
- Act as contracts for external dependencies
- Enable dependency inversion

**Services** (`service.go`):
- Contain business logic and use cases
- Orchestrate domain operations
- Implement business rules and validation

### Adapters (Outside Hexagon)

Located in `adapter/`, these implement the ports defined by the domain.

#### Repository Adapter (`adapter/repository/`)
- **File**: `mysql.go`
- **Purpose**: Implements `todo.Repository` interface
- **Technology**: MySQL with GORM
- **Responsibilities**:
  - Database connection management
  - CRUD operations for todos
  - Data mapping between domain and persistence models

#### Storage Adapter (`adapter/storage/`)
- **File**: `s3.go`
- **Purpose**: Implements `file.Repository` interface
- **Technology**: AWS S3 (LocalStack for development)
- **Responsibilities**:
  - File upload to S3
  - File download from S3
  - File deletion from S3

#### Streaming Adapter (`adapter/streaming/`)
- **File**: `redis.go`
- **Purpose**: Implements `todo.StreamRepository` interface
- **Technology**: Redis Streams
- **Responsibilities**:
  - Publishing todo events to streams
  - Event serialization
  - Asynchronous event processing

### Interface Adapters (Outside Hexagon)

Located in `internal/interfaces/`, these handle external communication.

#### HTTP Adapter (`internal/interfaces/http/`)
- **Files**: `handlers/`, `router.go`, `middleware.go`
- **Purpose**: REST API interface
- **Technology**: Gin framework
- **Responsibilities**:
  - HTTP request/response handling
  - Input validation
  - Error handling and formatting
  - CORS and middleware

## Data Flow

### Request Flow
```
HTTP Request → HTTP Handler → Domain Service → Repository Adapter → Database
                    ↓              ↓
              Domain Service → Streaming Adapter → Redis
```

### Response Flow
```
Database → Repository Adapter → Domain Service → HTTP Handler → HTTP Response
```

## Port Interfaces

### Todo Domain Ports

#### Repository Interface
```go
type Repository interface {
    Create(ctx context.Context, todo *TodoItem) error
    GetByID(ctx context.Context, id uuid.UUID) (*TodoItem, error)
    List(ctx context.Context, limit, offset int) ([]*TodoItem, error)
    Update(ctx context.Context, todo *TodoItem) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

#### StreamRepository Interface
```go
type StreamRepository interface {
    PublishTodoCreated(ctx context.Context, todo *TodoItem) error
}
```

### File Domain Ports

#### Repository Interface
```go
type Repository interface {
    Upload(ctx context.Context, filename string, content io.Reader, contentType string) (string, error)
    Download(ctx context.Context, fileID string) (io.ReadCloser, error)
    Delete(ctx context.Context, fileID string) error
}
```

## Dependency Injection

The main application (`cmd/server/main.go`) orchestrates the dependency injection:

```go
// 1. Create infrastructure adapters
db := repository.NewGormConnection(cfg.DatabaseURL)
todoRepo := repository.NewTodoRepository(db)
fileRepo := storage.NewFileRepository(s3Client, bucket)
streamRepo := streaming.NewStreamRepository(redisClient)

// 2. Create domain services
todoService := todo.NewTodoService(todoRepo, streamRepo)
fileService := file.NewFileService(fileRepo)

// 3. Create interface adapters
todoHandler := handlers.NewTodoHandler(todoService)
fileHandler := handlers.NewFileHandler(fileService)

// 4. Start HTTP server
router.SetupRouter(todoHandler, fileHandler)
```

## Testing Strategy

### Unit Testing
- **Domain Services**: Test business logic with mocked ports
- **Adapters**: Test integration with external systems
- **Handlers**: Test HTTP request/response handling

### Test Structure
```
tests/
├── unit_test.go      # Domain service tests
├── benchmark_test.go # Performance tests
└── integration_test.go # End-to-end tests
```

### Mocking Strategy
- Domain services use mock implementations of ports
- Adapters can be tested with test databases/containers
- Handlers can be tested with mock services

## Benefits of This Architecture

### 1. Testability
- Domain logic can be tested in isolation
- Easy to mock external dependencies
- Fast unit tests without external systems

### 2. Maintainability
- Clear separation of concerns
- Easy to understand and modify
- Changes in one layer don't affect others

### 3. Flexibility
- Easy to swap implementations (e.g., different databases)
- Can add new interfaces without changing domain
- Supports multiple delivery mechanisms

### 4. Independence
- Domain logic is framework-agnostic
- Can change external systems without affecting business logic
- Technology choices are isolated to adapters

## Configuration Management

Configuration is managed through the `pkg/config` package:

- Environment-based configuration
- Validation on startup
- Type-safe configuration structs
- Support for different environments

## Error Handling

### Domain Layer
- Business rule violations return domain errors
- No external system dependencies in error handling

### Adapter Layer
- External system errors are wrapped and returned
- Logging of external system interactions

### Interface Layer
- HTTP status code mapping
- User-friendly error messages
- Request tracing and logging

## Future Enhancements

### Potential Additions
- **Authentication Adapter**: JWT/OAuth2 implementation
- **Caching Adapter**: Redis caching layer
- **Metrics Adapter**: Prometheus metrics collection
- **Message Queue Adapter**: RabbitMQ/Kafka integration
- **GraphQL Interface**: Alternative to REST API

### Scalability Considerations
- Horizontal scaling of domain services
- Database sharding strategies
- Event-driven architecture expansion
- Microservices decomposition

## Best Practices

### Domain Layer
- Keep business logic pure and testable
- Use value objects for complex data
- Implement proper validation
- Avoid external dependencies

### Adapter Layer
- Implement ports faithfully
- Handle external system errors gracefully
- Use appropriate logging levels
- Implement retry mechanisms where needed

### Interface Layer
- Validate input thoroughly
- Provide clear error messages
- Use appropriate HTTP status codes
- Implement proper middleware

This architecture provides a solid foundation for building maintainable, testable, and scalable applications while keeping the business logic clean and independent of external concerns.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                    HTTP Interface Layer                     │
│  ┌─────────────────┐  ┌─────────────────┐                   │
│  │  Todo Handler   │  │  File Handler   │                   │
│  └─────────────────┘  └─────────────────┘                   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Domain Layer (Hexagon)                   │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐  │
│  │  Todo Service   │  │  File Service   │  │   Shared    │  │
│  │                 │  │                 │  │   Domain    │  │
│  │  ┌───────────┐  │  │  ┌───────────┐  │  │             │  │
│  │  │   Ports   │  │  │  │   Ports   │  │  │  ┌───────┐  │  │
│  │  │           │  │  │  │           │  │  │  │ Ports │  │  │
│  │  │ Repository│  │  │  │ Repository│  │  │  │Cache  │  │  │
│  │  │ Messaging │  │  │  │ Storage   │  │  │  │Storage│  │  │
│  │  │ Cache     │  │  │  │           │  │  │  │Msg    │  │  │
│  │  └───────────┘  │  │  └───────────┘  │  │  └───────┘  │  │
│  └─────────────────┘  └─────────────────┘  └─────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              ▲
                              │
┌─────────────────────────────────────────────────────────────┐
│                  Infrastructure Adapters                    │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │   MySQL     │  │     S3      │  │   Redis     │          │
│  │  Adapter    │  │  Adapter    │  │  Adapter    │          │
│  │             │  │             │  │             │          │
│  │ Repository  │  │ Storage     │  │ Messaging   │          │
│  │             │  │             │  │ Cache       │          │
│  └─────────────┘  └─────────────┘  └─────────────┘          │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                  External Systems                           │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │   MySQL     │  │     S3      │  │   Redis     │          │
│  │  Database   │  │   Storage   │  │  Streams    │          │
│  │             │  │             │  │  Cache      │          │
│  └─────────────┘  └─────────────┘  └─────────────┘          │
└─────────────────────────────────────────────────────────────┘
```

**Key Points:**
- **Domain Layer** is at the center (hexagon)
- **Dependencies point inward** toward the domain
- **Adapters implement ports** defined by the domain
- **External systems** are accessed only through adapters
- **Shared Domain** provides common utilities and interfaces
- **All adapters** are grouped together for consistency

## Recent Architecture Improvements (v2.0.0)

### Centralized Shared Ports
- **Location**: `internal/domain/shared/ports.go`
- **Interfaces**: `Cache`, `Storage`, `Messaging`, `EventStore`
- **Benefits**: Reusable across domains, consistent interface definitions

### Enhanced Folder Structure
- **Moved**: `internal/interfaces/http` → `adapter/http`
- **Result**: All adapters grouped in `adapter/` directory
- **Structure**: `adapter/http`, `adapter/repository`, `adapter/storage`, `adapter/streaming`, `adapter/cache`

### Shared Domain Utilities
- **Errors**: `internal/domain/shared/errors.go` - Domain-specific error types
- **Value Objects**: `internal/domain/shared/value_objects.go` - Reusable value objects
- **Benefits**: Consistent error handling and data types across domains

### Reusable Middleware Package
- **Location**: `pkg/middleware/`
- **Components**: CORS, Logging, Recovery, Timeout
- **Benefits**: Reusable across projects, configurable middleware

### Enhanced Port Interfaces
- **Todo Domain**: Uses `Messaging` and `Cache` from shared ports
- **File Domain**: Uses `Storage` from shared ports
- **Benefits**: Consistent interfaces, easier to swap implementations

### Improved Adapter Implementations
- **S3 Storage**: Implements `shared.Storage` interface
- **Redis Messaging**: Implements `shared.Messaging` interface
- **Redis Cache**: New adapter implementing `shared.Cache` interface
- **File Repository**: New adapter for file metadata persistence

### Updated Service Architecture
- **Todo Service**: Now uses `Messaging` and `Cache` interfaces
- **File Service**: Enhanced with proper domain models and storage integration
- **Benefits**: Better separation of concerns, more testable
