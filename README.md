# Todo Service

A production-ready Go-based todo service built with Clean Architecture principles, featuring file uploads to S3, MySQL persistence, Redis streaming, comprehensive error handling, and structured logging.

## Features

- **Clean Architecture**: Well-organized layers with proper separation of concerns
- **File Upload**: Upload files to S3 (LocalStack for development) with validation
- **Todo Management**: Full CRUD operations for todo items
- **Streaming**: Publish todo events to Redis streams
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
   cd todo-service
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

- `GET /health` - Health check
- `POST /upload` - Upload a file
- `POST /todo` - Create a new todo item
- `GET /todo/:id` - Get a todo by ID
- `GET /todo` - List todos (supports ?limit=10&offset=0)
- `PUT /todo/:id` - Update a todo item
- `DELETE /todo/:id` - Delete a todo item

For detailed API documentation, see [docs/api.md](docs/api.md).

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
todo-service/
├── cmd/server/           # Application entry point
├── internal/
│   ├── domain/          # Business entities and interfaces
│   │   ├── entities/    # Core business objects
│   │   └── repositories/ # Repository interfaces
│   ├── service/         # Business logic and use cases
│   ├── interfaces/      # HTTP handlers and routing
│   │   └── http/
│   │       ├── handlers/ # HTTP handlers
│   │       ├── middleware.go # Middleware functions
│   │       └── router.go # Route definitions
│   └── adapter/         # External service implementations
│       ├── database/    # Database adapters
│       ├── storage/     # Storage adapters (S3)
│       └── streaming/   # Streaming adapters (Redis)
├── pkg/config/          # Configuration management
├── tests/               # Test files
├── docs/                # Documentation
├── docker-compose.yml   # Docker services definition
├── Dockerfile          # Application container
├── Makefile            # Build automation
└── README.md
```

## Architecture

This project follows Clean Architecture principles:

- **Entities**: Core business objects (TodoItem)
- **Use Cases**: Business logic and orchestration
- **Interfaces**: HTTP handlers and external service contracts
- **Infrastructure**: Database, S3, Redis implementations

Dependencies point inward, with the domain layer being completely independent of external frameworks.

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
