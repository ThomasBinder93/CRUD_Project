# CRUD Backend - Production-Grade Demo

A professional-grade REST API backend built with Go, showcasing best practices.

## Overview

This project demonstrates a production-ready CRUD (Create, Read, Update, Delete) API with:
- Clean, layered architecture
- Comprehensive error handling
- Input validation
- Structured logging
- Full test coverage
- Professional API documentation
- Docker containerization
- Kubernetes deployment support

## Features

✅ **RESTful API** - Follows REST conventions with proper HTTP methods and status codes  
✅ **SQLite Database** - Simple yet powerful local database with proper repository pattern  
✅ **Clean Architecture** - Organized into logical packages (config, models, database, handlers, middleware)  
✅ **Error Handling** - Consistent, standardized error responses across the API  
✅ **Input Validation** - Validates all inputs before processing  
✅ **Logging** - Structured logging using Go's built-in `slog` package  
✅ **Middleware** - CORS, logging, panic recovery, and content-type handling  
✅ **Tests** - Unit and integration tests with good coverage  
✅ **End-to-End Tests** - Playwright-based browser tests for full UI validation  
✅ **Configuration** - Environment-based configuration management  
✅ **Code Quality** - Linting configuration with golangci-lint  
✅ **Documentation** - Comprehensive API documentation  
✅ **Docker Support** - Production-ready Dockerfile with multi-stage build  
✅ **Kubernetes Ready** - Deployment and service manifests included  

## Project Structure

```
.
├── config/           # Configuration management
├── database/         # Database layer (repository pattern)
├── handlers/         # HTTP request handlers
├── middleware/       # HTTP middleware
├── models/          # Domain models with validation
├── utils/           # Utility functions and error handling
├── static/          # Frontend (HTML/JavaScript)
├── tests/           # End-to-end tests (Playwright)
├── main.go          # Application entry point
├── package.json     # Node.js dependencies for testing
├── playwright.config.js # Playwright configuration
├── Dockerfile       # Docker image configuration
├── docker-compose.yml # Local development setup
├── .golangci.yml    # Linting configuration
├── .env.example     # Environment variables example
├── Makefile         # Build and development commands
└── API_DOCUMENTATION.md # API specification
```

## Prerequisites

- Go 1.21 or higher
- Node.js 18+ and npm (for end-to-end tests)
- Docker (optional, for containerization)
- Make (optional, for convenience commands)

## Quick Start

### Local Development

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd CRUD_Project
   ```

2. Copy environment configuration:
   ```bash
   cp .env.example .env
   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

4. Run the application:
   ```bash
   make run
   # or
   go run main.go
   ```

The server will start on `http://localhost:8080`

### With Docker

```bash
docker-compose up
```

## Development Commands

### Building
```bash
make build          # Build binary
make run            # Build and run
```

### Testing
```bash
make test           # Run tests with coverage report
make test-unit      # Run unit tests only
make test-e2e       # Run end-to-end tests (requires Node.js)
```

#### End-to-End Tests
The project includes Playwright-based end-to-end tests that verify the full application functionality through the browser interface.

**Setup:**
```bash
npm install
npx playwright install
```

**Run tests:**
```bash
npm test                    # Run all tests
npm run test:headed         # Run tests with browser visible
npm run test:ui             # Run tests with Playwright UI
```

The e2e tests cover:
- Page loading and UI elements
- Creating items (valid and invalid inputs)
- Updating items
- Deleting items
- Error handling and validation feedback

### Code Quality
```bash
make lint           # Run linter
make fmt            # Format code
```

### Docker
```bash
make docker-build   # Build Docker image
make docker-run     # Run Docker container
make docker-stop    # Stop Docker container
```

### Cleanup
```bash
make clean          # Remove build artifacts
```

## API Endpoints

### Health Check
- `GET /health` - Check server health

### Items Management
- `GET /api/items` - Get all items
- `POST /api/items` - Create new item
- `GET /api/items/{id}` - Get specific item
- `PUT /api/items/{id}` - Update item
- `DELETE /api/items/{id}` - Delete item

For detailed API documentation, see [API_DOCUMENTATION.md](API_DOCUMENTATION.md)

## Configuration

Configure the application using environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | 8080 | Server port |
| `SERVER_READ_TIMEOUT` | 15 | Read timeout in seconds |
| `SERVER_WRITE_TIMEOUT` | 15 | Write timeout in seconds |
| `DB_PATH` | ./crud.db | SQLite database file path |
| `LOG_LEVEL` | info | Logging level (debug, info, warn, error) |
| `APP_VERSION` | 1.0.0-FIXED | Application version displayed in startup logs |
| `APP_ENV` | production | Runtime environment displayed in startup logs |

See `.env.example` for all available options.

## Error Handling

The API uses standardized error responses:

```json
{
  "error": {
    "code": 400,
    "message": "bad request",
    "details": "item name cannot be empty"
  }
}
```

Error codes follow HTTP standards:
- `400` - Bad Request (invalid input format)
- `404` - Not Found (resource doesn't exist)
- `422` - Unprocessable Entity (validation error)
- `500` - Internal Server Error

## Architecture Highlights

### Clean Separation of Concerns
- **config**: Centralized configuration
- **models**: Domain models with validation
- **database**: Repository pattern for data access
- **handlers**: HTTP request/response handling
- **middleware**: Cross-cutting concerns
- **utils**: Common utilities

### Database Layer
- Thread-safe operations using sync.RWMutex
- Proper error wrapping with context
- Connection pooling via sql.DB
- Automatic schema initialization

### Request Handling
- Input validation before processing
- Consistent error responses
- Request body size limits (1MB)
- Panic recovery middleware

### Logging
- Structured logging with slog
- Request/response logging middleware
- Error logging with context
- Configurable log levels

## Testing

Run the test suite:

```bash
# Run all tests with coverage
go test -v -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out
```

Tests are provided for:
- Database repository operations
- Item model validation
- HTTP handlers
- Error handling

## Deployment

### Docker Deployment

```bash
docker build -t crud-demo:latest .
docker run -p 8080:8080 crud-demo:latest
```

### Kubernetes Deployment

```bash
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
kubectl apply -f ingress.yaml
```

## Code Quality

The project uses `golangci-lint` for code quality:

```bash
golangci-lint run ./...
```

Configuration is in `.golangci.yml` and includes checks for:
- Error handling
- Code simplification
- Unused code
- Variable naming conventions
- And more...

## Production Considerations

This project demonstrates several production-grade features:

1. **Graceful Shutdown** - Server shuts down cleanly on SIGTERM
2. **Error Wrapping** - Errors are wrapped with context using `%w` format
3. **Logging** - Structured logging for observability
4. **Middleware** - CORS, recovery, and logging middleware
5. **Configuration** - Environment-based, no hardcoded values
6. **Database** - Connection pooling and thread-safe operations
7. **Testing** - Unit and integration tests
8. **Documentation** - Comprehensive API documentation

## Frontend

The project includes a basic frontend at `http://localhost:8080` with:
- Item listing
- Create new items
- Update items
- Delete items
- Real-time status messages

Built with vanilla JavaScript and Fetch API.

## Troubleshooting

### Port Already in Use
Change the port via environment variable:
```bash
SERVER_PORT=3000 go run main.go
```

### Database Issues
Delete the database file and restart:
```bash
rm crud.db
go run main.go
```

### Permission Errors (Docker)
Ensure proper permissions:
```bash
docker run -p 8080:8080 --user root crud-demo:latest
```

## Contributing

When contributing, ensure:
- Code follows Go conventions
- Tests are included for new features
- Linter passes (`make lint`)
- Code is formatted (`make fmt`)
- Documentation is updated

## Learning Resources

This project demonstrates:
- Clean architecture principles
- Go best practices
- REST API design
- Error handling patterns
- Testing strategies
- Docker containerization
- Kubernetes deployment

## License

MIT

## Next Steps for Production

To make this production-ready, consider adding:
1. **Database migrations** - Use a migration tool (golang-migrate, tern, etc.)
2. **Authentication** - JWT or OAuth2
3. **Rate limiting** - Protect against abuse
4. **Caching** - Redis for improved performance
5. **Monitoring** - Prometheus metrics and health checks
6. **API Gateway** - Kong, Nginx, or similar
7. **CI/CD** - GitHub Actions, GitLab CI, etc.
8. **Database** - Upgrade from SQLite to PostgreSQL/MySQL for production
9. **Secrets management** - Use HashiCorp Vault or similar
10. **Load testing** - Benchmark and optimize performance
kubectl apply -f ingress.yaml
```

### Docker Compose (mit DB-Volume)

```yaml
version: '3'
services:
  crud-demo:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./crud.db:/app/crud.db
```

## Projektstruktur

```
CRUD_Project/
├── main.go              # Go-Backend
├── go.mod               # Go-Module
├── crud.db              # SQLite-Datenbank (wird erstellt)
├── static/
│   ├── index.html       # Frontend-HTML
│   ├── script.js        # Frontend-JavaScript
│   └── helper.js        # Hilfsfunktionen
├── Dockerfile           # Docker-Image
├── docker-compose.yml   # Docker Compose
├── deployment.yaml      # Kubernetes Deployment
├── service.yaml         # Kubernetes Service
└── ingress.yaml         # Kubernetes Ingress
```

## Lizenz

Dieses Projekt ist unter der MIT-Lizenz lizenziert.