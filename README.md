# TaskFlow

A production-ready Go-based todo service built with **Hexagonal Architecture** (Ports & Adapters), featuring file uploads to S3, MySQL persistence, Redis streaming, comprehensive error handling, and structured logging.

## Features

- **Hexagonal Architecture**: Clean separation between domain logic and infrastructure
- **File Upload**: Upload files to S3 (LocalStack for development) with validation
- **Todo Management**: Full CRUD operations for todo items
- **Event Streaming**: Publish todo events to Redis streams
- **Database**: MySQL with automated migrations
- **Testing**: Comprehensive unit tests and benchmarks
- **Docker**: Full containerization with Docker Compose
- **API Documentation**: Comprehensive API documentation
- **Middleware**: CORS, request logging, timeout handling
- **Configuration**: Environment-based configuration with validation

## Prerequisites

- Docker & Docker Compose
- Make
- Go 1.24+ (for local development)

## Quick Start

1. **Clone and run the service:**
   ```bash
   git clone <repository-url>
   cd taskflow
   make run
   ```

2. **Setup S3 bucket (after services are up):**
   ```bash
   make setup-localstack
   ```

3. **Test the API:**
   ```bash
   # Health check
   curl http://localhost:8080/health

   # Upload a file
   curl -X POST -F "file=@example.txt" http://localhost:8080/upload

   # Create a todo
   curl -X POST http://localhost:8080/todo \
     -H "Content-Type: application/json" \
     -d '{
       "description": "Complete the assignment",
       "dueDate": "2024-12-31T23:59:59Z",
       "fileId": "optional-file-id-from-upload"
     }'

   # List todos
   curl http://localhost:8080/todo?limit=5

   # Update a todo
   curl -X PUT http://localhost:8080/todo/todo-id \
     -H "Content-Type: application/json" \
     -d '{"description": "Updated description"}'

   # Delete a todo
   curl -X DELETE http://localhost:8080/todo/todo-id
   ```

## API Endpoints

### Health
- `GET /health` - Health check

### File Management
- `POST /upload` - Upload a file
- `GET /file/:id` - Get file metadata
- `GET /file/:id/download` - Download a file
- `DELETE /file/:id` - Delete a file

### Todo Management
- `POST /todo` - Create a new todo item
- `GET /todo/:id` - Get a todo by ID
- `GET /todo` - List todos (supports ?limit=10&offset=0)
- `PUT /todo/:id` - Update a todo item
- `DELETE /todo/:id` - Delete a todo item

For detailed API documentation, see [docs/api.md](docs/api.md).

For comprehensive architecture documentation, see [docs/architecture.md](docs/architecture.md).

For development guidelines and setup instructions, see [docs/development.md](docs/development.md).

For a quick reference guide, see [docs/quick-reference.md](docs/quick-reference.md).

For a complete list of changes and version history, see [CHANGELOG.md](CHANGELOG.md).

## Development Commands

```bash
# Run services locally (without app container)
make dev

# Run tests
make test

# Run benchmarks
make benchmark

# Clean up containers and volumes
make clean

# Update Go modules
make mod

# Run linter
make lint
```

## Project Structure

```
taskflow/
├── cmd/server/           # Application entry point
├── adapter/              # All adapters (infrastructure)
│   ├── cache/            # Cache adapters (Redis)
│   ├── http/             # HTTP interface adapters
│   │   ├── handlers/     # HTTP handlers
│   │   ├── middleware.go # Middleware functions
│   │   └── router.go     # Route definitions
│   ├── repository/       # Database adapters (MySQL)
│   ├── storage/          # Storage adapters (S3)
│   └── streaming/        # Event streaming adapters (Redis)
├── internal/
│   └── domain/           # Business logic (hexagon core)
│       ├── shared/       # Shared domain utilities
│       │   ├── ports.go  # Shared port interfaces
│       │   ├── errors.go # Domain errors
│       │   └── value_objects.go # Value objects
│       ├── todo/         # Todo domain
│       │   ├── model.go  # Domain models
│       │   ├── port.go   # Port interfaces
│       │   └── service.go # Domain services
│       └── file/         # File domain
│           ├── model.go  # Domain models
│           ├── port.go   # Port interfaces
│           └── service.go # Domain services
├── pkg/
│   ├── config/          # Configuration management
│   └── middleware/      # Reusable middleware
│       ├── cors.go      # CORS middleware
│       ├── logging.go   # Request logging
│       ├── recovery.go  # Panic recovery
│       └── timeout.go   # Request timeout
├── tests/               # Test files
├── docs/                # Documentation
├── docker-compose.yml   # Docker services definition
├── Dockerfile          # Application container
├── Makefile            # Build automation
└── README.md
```

## Architecture

This project follows **Hexagonal Architecture** (Ports & Adapters) principles:

### Core Concepts

- **Domain Layer** (`internal/domain/`): Contains the business logic and domain models
  - **Models**: Core business objects (TodoItem, File)
  - **Ports**: Interfaces that define what the domain needs (Repository, StreamRepository)
  - **Services**: Business logic and use cases

- **Adapters** (`adapter/`): Infrastructure implementations that adapt external systems
  - **Repository Adapter**: MySQL database implementation
  - **Storage Adapter**: S3 file storage implementation
  - **Streaming Adapter**: Redis event streaming implementation

- **Interface Adapters** (`internal/interfaces/`): External interface implementations
  - **HTTP Adapter**: REST API handlers and routing

### Key Benefits

- **Dependency Inversion**: Domain depends on abstractions, not concrete implementations
- **Testability**: Easy to mock dependencies for unit testing
- **Flexibility**: Easy to swap out implementations (e.g., different databases)
- **Maintainability**: Clear boundaries between business logic and infrastructure
- **Independence**: Domain layer is completely independent of external frameworks

### Dependency Flow

```
HTTP Interface → Domain Services → Port Interfaces
     ↓                ↓                ↓
Infrastructure Adapters → External Systems
```

Dependencies point inward toward the domain layer, ensuring business logic remains pure and testable.

## Hexagonal Architecture Implementation

### Domain Layer Structure

Each domain (`todo`, `file`) follows a consistent pattern:

```
internal/domain/{domain}/
├── model.go    # Domain models and DTOs
├── port.go     # Port interfaces (what the domain needs)
└── service.go  # Domain services (business logic)
```

### Port Interfaces

**Todo Domain Ports:**
- `Repository`: Data persistence operations
- `StreamRepository`: Event publishing operations

**File Domain Ports:**
- `Repository`: File storage operations

### Adapter Implementations

**Repository Adapter** (`adapter/repository/mysql.go`):
- Implements `todo.Repository` interface
- Handles MySQL database operations
- Uses GORM for ORM functionality

**Storage Adapter** (`adapter/storage/s3.go`):
- Implements `file.Repository` interface
- Handles S3 file operations
- Supports LocalStack for development

**Streaming Adapter** (`adapter/streaming/redis.go`):
- Implements `todo.StreamRepository` interface
- Handles Redis stream publishing
- Publishes todo events asynchronously

### Interface Adapters

**HTTP Adapter** (`internal/interfaces/http/`):
- REST API handlers
- Request/response mapping
- Error handling and validation
- Middleware for cross-cutting concerns

### Dependency Injection

The main application (`cmd/server/main.go`) wires everything together:

1. Creates infrastructure adapters
2. Injects them into domain services
3. Passes domain services to interface adapters
4. Starts the HTTP server

This ensures loose coupling and makes testing straightforward.

## Key Improvements

### Error Handling
- Custom error types for different scenarios
- Proper error wrapping and context
- Consistent error responses across all endpoints

### Logging
- Structured logging with slog
- Request tracing with unique request IDs
- Different log levels for different environments

### Validation
- Input validation at multiple layers
- Configuration validation on startup
- File upload validation (type and size)

### Middleware
- CORS support
- Request logging
- Timeout handling
- Panic recovery

### API Design
- RESTful endpoints
- Consistent response formats
- Proper HTTP status codes
- Pagination support

## Testing

The project includes:
- Unit tests for all layers
- Mocked external services
- Benchmarks for performance testing
- Integration tests for database operations

Run tests with:
```bash
make test
make benchmark
```

## Configuration

Environment variables:
- `PORT`: Server port (default: 8080)
- `DATABASE_URL`: MySQL connection string
- `REDIS_URL`: Redis connection string
- `AWS_REGION`: AWS region for S3
- `S3_BUCKET`: S3 bucket name
- `AWS_ACCESS_KEY_ID`: AWS access key
- `AWS_SECRET_ACCESS_KEY`: AWS secret key
- `S3_ENDPOINT`: S3 endpoint (for LocalStack)
- `ENVIRONMENT`: Environment (development/staging/production)
- `LOG_LEVEL`: Log level (debug/info/warn/error)

## Production Considerations

- Add authentication and authorization
- Implement rate limiting
- Add metrics and monitoring
- Use proper secrets management
- Set up CI/CD pipelines
- Add health checks for all dependencies
- Implement circuit breakers for external services
- Add distributed tracing
- Set up proper backup strategies

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

## License

This project is licensed under the MIT License.
