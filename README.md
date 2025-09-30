# Hexagonal Todo API

A Go-based todo service demonstrating **Hexagonal Architecture** (Ports & Adapters) with file uploads, MySQL persistence, and Redis streaming.

## 🏗️ Hexagonal Architecture

This project showcases clean architecture principles with clear separation between business logic and infrastructure concerns.

### Architecture Layers

```
┌─────────────────────────────────────────────────────────────┐
│                    HTTP Adapter                             │
│  ┌─────────────────┐  ┌─────────────────┐                   │
│  │   Todo Handler  │  │   File Handler  │                   │
│  └─────────────────┘  └─────────────────┘                   │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                    Domain Layer                             │
│  ┌─────────────────┐  ┌─────────────────┐                   │
│  │   Todo Service  │  │   File Service  │                   │
│  └─────────────────┘  └─────────────────┘                   │
│  ┌─────────────────────────────────────────────────────────┐│
│  │              Shared Domain                              ││
│  │  ┌─────────────┐ ┌─────────────┐ ┌──────────────┐       ││
│  │  │   Ports     │ │   Errors    │ │ Value Objects│       ││
│  │  └─────────────┘ └─────────────┘ └──────────────┘       ││
│  └─────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                  Driven Adapters                            │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────┐│
│  │   MySQL     │ │     S3      │ │   Redis     │ │  Cache  ││
│  │ Repository  │ │   Storage   │ │ Messaging   │ │ Adapter ││
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────┘│
└─────────────────────────────────────────────────────────────┘
```

### Key Principles

1. **Domain-Centric**: Business logic is independent of external frameworks
2. **Dependency Inversion**: Dependencies point inward toward the domain
3. **Ports & Adapters**: External systems accessed through interfaces
4. **Testability**: Easy to mock dependencies for unit testing

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.21+

### Run with Docker
```bash
# Clone and start services
git clone <your-repo>
cd hexagonal-todo-api
docker-compose up -d

# The API will be available at http://localhost:8080
```

### Run Locally
```bash
# Install dependencies
go mod tidy

# Set environment variables
export DATABASE_URL="mysql://user:password@localhost:3306/taskflow"
export REDIS_URL="localhost:6379"
export S3_ENDPOINT="http://localhost:4566"
export S3_ACCESS_KEY_ID="test"
export S3_SECRET_ACCESS_KEY="test"
export S3_BUCKET="taskflow-files"

# Run the application
go run cmd/server/main.go
```

## 📁 Project Structure

```
hexagonal-todo-api/
├── cmd/server/           # Application entry point
├── internal/
│   └── domain/          # Domain layer (business logic)
│       ├── todo/        # Todo domain
│       ├── file/        # File domain
│       └── shared/      # Shared domain utilities
├── adapter/             # Adapters (infrastructure)
│   ├── http/           # HTTP adapter (handlers, router)
│   ├── repository/     # Database adapters
│   ├── storage/        # File storage adapters
│   ├── streaming/      # Event streaming adapters
│   └── cache/          # Caching adapters
├── pkg/                # Shared packages
│   ├── config/         # Configuration
│   └── middleware/     # Reusable middleware
└── tests/              # Test files
```

## 🧪 Testing

```bash
# Run unit tests
go test ./tests/...

# Run benchmarks
go test -bench=. ./tests/...

# Run with coverage
go test -cover ./tests/...
```

## 📚 Learning Resources

### Hexagonal Architecture Concepts

1. **Ports**: Interfaces defined by the domain
2. **Adapters**: Implementations of ports
3. **Domain Services**: Business logic implementation
4. **Dependency Injection**: Wiring dependencies in main.go

### Key Files to Study

- `internal/domain/todo/service.go` - Business logic
- `internal/domain/shared/ports.go` - Shared interfaces
- `adapter/repository/mysql.go` - Database adapter
- `cmd/server/main.go` - Dependency injection

## 🔧 Configuration

Environment variables:
- `DATABASE_URL`: MySQL connection string
- `REDIS_URL`: Redis connection string
- `S3_ENDPOINT`: S3-compatible storage endpoint
- `S3_ACCESS_KEY_ID`: S3 access key
- `S3_SECRET_ACCESS_KEY`: S3 secret key
- `S3_BUCKET`: S3 bucket name
- `PORT`: Server port (default: 8080)

## 📖 API Documentation

See [API Documentation](docs/api.md) for detailed endpoint information.

## 🤝 Contributing

This is an educational project. Feel free to:
- Study the hexagonal architecture implementation
- Experiment with different adapters
- Add new features following the same patterns
- Improve the documentation

## 📄 License

MIT License - feel free to use this for learning and educational purposes.