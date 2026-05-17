# Architecture Guide

This document explains the architecture and design decisions of the CRUD project.

## Overview

The CRUD project follows a **layered, domain-driven design** with clear separation of concerns. Each package has a specific responsibility and minimal dependencies on other packages.

## Package Structure

### `config/`
**Responsibility**: Configuration management

- Loads configuration from environment variables
- Provides defaults for all settings
- Centralizes all configuration logic
- No business logic

### `models/`
**Responsibility**: Domain models and validation

- Defines core data structures (Item, CreateItemRequest, UpdateItemRequest)
- Implements validation logic
- No database or HTTP knowledge
- Pure data and validation logic

### `database/`
**Responsibility**: Data persistence layer

- Repository pattern for data access abstraction
- SQLite implementation (easily swappable)
- Thread-safe operations using sync.RWMutex
- Proper error handling with context
- No business logic, pure data operations

### `handlers/`
**Responsibility**: HTTP request/response handling

- Converts HTTP requests to domain operations
- Validates input using domain models
- Converts domain models to JSON responses
- Orchestrates database and business logic calls
- No database implementation details

### `middleware/`
**Responsibility**: Cross-cutting HTTP concerns

- **logging.go**: Request/response logging middleware
- **recovery.go**: Panic recovery and error handling
- **cors.go**: CORS and content-type headers

### `utils/`
**Responsibility**: Shared utilities

- **response.go**: Standardized response formatting
- **errors.go**: API error definitions
- **decoder.go**: JSON decoding with size limits

### `tests/`
**Responsibility**: Integration tests (can be expanded)

## Design Patterns Used

### Repository Pattern
The `database` package uses the repository pattern to abstract data access:
```go
type Repository interface {
    GetAll() ([]models.Item, error)
    GetByID(id int) (*models.Item, error)
    Create(name string) (*models.Item, error)
    Update(id int, name string) (*models.Item, error)
    Delete(id int) error
}
```

Benefits:
- Easy to swap implementations (e.g., PostgreSQL, MongoDB)
- Easy to mock for testing
- Business logic doesn't depend on database details

### Middleware Chain
HTTP middleware is applied in order:
1. Recovery (catch panics)
2. CORS (allow cross-origin)
3. Content-Type (set headers)
4. Logging (request details)
5. Route handlers

### Dependency Injection
Dependencies are passed to constructors:
```go
handler := handlers.NewItemHandler(repo, logger)
```

Benefits:
- Easy to test with mocks
- Loose coupling
- Clear dependencies

## Error Handling Strategy

### Layered Error Responses

1. **Database Layer**: Wraps errors with context
   ```go
   return nil, fmt.Errorf("failed to create item: %w", err)
   ```

2. **Handler Layer**: Converts to API errors
   ```go
   utils.WriteError(w, utils.ErrorInternalServer("failed to create item"))
   ```

3. **Client**: Receives standardized error
   ```json
   {
     "error": {
       "code": 500,
       "message": "internal server error",
       "details": "failed to create item"
     }
   }
   ```

### Error Types
- **API Errors**: HTTP status with message (400, 404, 422, 500)
- **Domain Errors**: Validation errors from models
- **DB Errors**: Wrapped database errors

## Data Flow

### Create Item Flow
```
HTTP POST /api/items
  ↓
Handler.Create()
  ↓ Validate input
models.CreateItemRequest.Validate()
  ↓ Create in database
database.Repository.Create()
  ↓ Generate response
utils.WriteSuccess()
  ↓
JSON Response {data: {...}}
```

### Get Item by ID Flow
```
HTTP GET /api/items/{id}
  ↓
Handler.GetByID()
  ↓ Parse ID
strconv.Atoi()
  ↓ Query database
database.Repository.GetByID()
  ↓ Check existence
  ↓ Generate response
utils.WriteSuccess() or WriteError()
  ↓
JSON Response
```

## Thread Safety

The repository uses `sync.RWMutex` to ensure thread-safe operations:

```go
type Repository struct {
    db *sql.DB
    mu sync.RWMutex
}

func (r *Repository) GetAll() ([]models.Item, error) {
    r.mu.RLock()      // Allow concurrent reads
    defer r.mu.RUnlock()
    // ...
}

func (r *Repository) Create(name string) (*models.Item, error) {
    r.mu.Lock()        // Exclusive write lock
    defer r.mu.Unlock()
    // ...
}
```

## Testing Strategy

### Unit Tests
- Test individual functions/methods in isolation
- Use in-memory SQLite for database tests
- Mock external dependencies

### Integration Tests
- Test multiple layers together
- Use real HTTP requests via httptest
- Verify complete flows

### Test Structure
```
database/
  └── repository_test.go    # Database layer tests

handlers/
  └── item_test.go          # Handler layer tests

models/
  └── (validation in unit tests)
```

## Logging

The application uses Go's built-in `slog` package for structured logging:

```go
logger.Info("request handled",
    slog.String("method", r.Method),
    slog.String("path", r.RequestURI),
    slog.Int("status", wrapped.statusCode),
    slog.Duration("duration", duration),
)
```

Benefits:
- Structured, machine-readable logs
- Multiple output formats (text, JSON)
- Configurable log levels
- No external dependencies

## Configuration Management

Configuration is centralized in the `config` package:

```go
cfg := config.Load()
```

- Reads from environment variables
- Provides sensible defaults
- No hardcoded values
- Easy to test with different configs

## Performance Considerations

1. **Connection Pooling**: SQLite via sql.DB automatically pools connections
2. **Read/Write Locks**: RWMutex allows concurrent reads
3. **Request Validation**: Early validation prevents unnecessary DB calls
4. **Response Formatting**: Efficient JSON encoding with standard library
5. **Middleware Order**: Logging after recovery for better observability

## Scalability Patterns

To scale this application:

1. **Database**: Replace SQLite with PostgreSQL/MySQL
2. **Caching**: Add Redis for frequently accessed items
3. **Rate Limiting**: Add rate limiting middleware
4. **API Gateway**: Use Nginx/Kong for routing
5. **Load Balancing**: Deploy multiple instances
6. **Message Queue**: Add for async operations
7. **Monitoring**: Integrate Prometheus metrics

## Security Considerations

1. **Input Validation**: All inputs validated before processing
2. **Request Size Limits**: 1MB limit on request bodies
3. **Error Messages**: Detailed errors only in logs, generic in responses
4. **CORS**: Configurable (currently allows all)
5. **Dependencies**: Minimal, all from trusted sources
6. **SQL Injection**: Protected by parameterized queries

## Future Improvements

1. Implement database migrations
2. Add authentication/authorization
3. Add API versioning
4. Implement soft deletes
5. Add request/response tracing
6. Implement pagination
7. Add filtering/sorting
8. Implement API rate limiting
