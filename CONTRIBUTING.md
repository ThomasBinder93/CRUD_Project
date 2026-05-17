# Contributing Guide

Thank you for your interest in contributing to the CRUD project! This guide will help you get started.

## Development Setup

### Prerequisites
- Go 1.21 or higher
- Git
- Make (optional but recommended)
- golangci-lint (for code quality checks)

### Clone and Setup
```bash
git clone <repository-url>
cd CRUD_Project
go mod tidy
```

## Development Workflow

### 1. Create a Feature Branch
```bash
git checkout -b feature/your-feature-name
```

Use clear, descriptive branch names:
- `feature/add-pagination` for new features
- `fix/incorrect-error-handling` for bug fixes
- `docs/update-readme` for documentation
- `refactor/database-layer` for refactoring

### 2. Make Your Changes

Follow Go best practices:
- Use meaningful variable/function names
- Keep functions small and focused
- Write clear comments for complex logic
- Add tests for new code

### 3. Format and Lint Your Code

```bash
# Format code
make fmt

# Run linter
make lint
```

Fix any linter warnings before submitting.

### 4. Write Tests

All new functionality should include tests:

```bash
# Run tests
make test

# Run specific test
go test -v ./database -run TestRepositoryCreate
```

Aim for at least 80% code coverage.

### 5. Commit with Clear Messages

```bash
git add .
git commit -m "feat: add pagination to items endpoint"
```

Use conventional commits:
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation
- `refactor:` - Code refactoring
- `test:` - Adding tests
- `chore:` - Maintenance

### 6. Push and Create Pull Request

```bash
git push origin feature/your-feature-name
```

Create a pull request on GitHub with:
- Clear description of changes
- Reference any related issues
- Screenshot/examples if applicable

## Code Standards

### Go Code Style
- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use `gofmt` (included in Go)
- Exported functions should have comments
- Keep lines under 120 characters

### Error Handling
Always check errors:
```go
// ✅ Good
if err != nil {
    logger.Error("operation failed", slog.String("error", err.Error()))
    return nil, fmt.Errorf("failed to do something: %w", err)
}

// ❌ Bad
_ = someFunction() // Ignoring error
```

### Testing
- Test files end with `_test.go`
- Test functions start with `Test`
- Use table-driven tests for multiple scenarios
- Mock external dependencies

Example:
```go
func TestRepositoryCreate(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid item", "Test", false},
        {"empty name", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test code
        })
    }
}
```

### Documentation
- Add comments for exported functions
- Update README for user-facing changes
- Update ARCHITECTURE.md for design changes
- Add API_DOCUMENTATION.md entries for new endpoints

Example comment:
```go
// GetByID retrieves a single item by its ID.
// Returns nil if the item doesn't exist.
// Returns an error if the database query fails.
func (r *Repository) GetByID(id int) (*models.Item, error) {
    // ...
}
```

## Adding Features

### Adding a New Endpoint

1. **Add to models/** (if needed):
```go
type NewRequest struct {
    Field string `json:"field"`
}

func (r *NewRequest) Validate() error {
    // validation logic
}
```

2. **Add to database/**:
```go
func (r *Repository) NewOperation() (*models.Item, error) {
    // implementation
}
```

3. **Add to handlers/**:
```go
func (h *ItemHandler) NewEndpoint(w http.ResponseWriter, r *http.Request) {
    // request handling
    utils.WriteSuccess(w, http.StatusOK, data)
}
```

4. **Wire in main.go**:
```go
router.HandleFunc("/api/endpoint", handler.NewEndpoint).Methods(http.MethodPost)
```

5. **Add tests** in corresponding `_test.go` file

6. **Update documentation**:
   - API_DOCUMENTATION.md
   - README.md (if public-facing)

### Adding Configuration Option

1. **Add to config/config.go**:
```go
type Config struct {
    NewOption string
}

func Load() *Config {
    return &Config{
        NewOption: getEnv("NEW_OPTION", "default"),
    }
}
```

2. **Update .env.example**:
```env
# New option description
NEW_OPTION=default_value
```

3. **Use in application**:
```go
cfg := config.Load()
value := cfg.NewOption
```

## Testing Guidelines

### Run Tests
```bash
# All tests
make test

# With coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Specific package
go test -v ./database

# Specific test
go test -v ./database -run TestCreate
```

### Coverage Requirements
- Aim for 80%+ coverage
- Focus on critical paths
- Don't test standard library
- Test error cases

### Integration Tests
For testing multiple layers:
```go
func TestEndToEnd(t *testing.T) {
    // Setup
    repo := setupTestDB(t)
    handler := handlers.NewItemHandler(repo, logger)
    
    // Execute
    req, _ := http.NewRequest(http.MethodPost, "/api/items", ...)
    rr := httptest.NewRecorder()
    handler.Create(rr, req)
    
    // Assert
    if rr.Code != http.StatusCreated {
        t.Errorf("expected 201, got %d", rr.Code)
    }
}
```

## Documentation Standards

### README Changes
If your changes affect users:
1. Update README.md with feature description
2. Update API_DOCUMENTATION.md with endpoint details
3. Update .env.example with new configuration options

### Code Comments
```go
// getEnv retrieves an environment variable with a fallback default.
// If the variable is not set, defaultValue is returned.
func getEnv(key, defaultValue string) string {
    // ...
}
```

### Commit Messages
```
feat: add user authentication

- Implement JWT token generation
- Add middleware for token validation
- Update API documentation

Closes #123
```

## Common Tasks

### Running the Application
```bash
make run
```

### Building for Production
```bash
make build
```

### Creating a Docker Image
```bash
make docker-build
```

### Checking Code Quality
```bash
make lint
make fmt
```

### Running Full Test Suite
```bash
make test
```

## Pull Request Checklist

Before submitting, ensure:
- [ ] Code follows Go conventions
- [ ] All tests pass (`make test`)
- [ ] Linter passes (`make lint`)
- [ ] Code is formatted (`make fmt`)
- [ ] Documentation is updated
- [ ] Commits have clear messages
- [ ] No hardcoded values
- [ ] Error handling is proper
- [ ] Tests cover new code
- [ ] No breaking changes (or clearly documented)

## Questions or Issues?

- Check existing issues first
- Read ARCHITECTURE.md for design decisions
- Review similar code in the codebase
- Open an issue for discussion

## Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Help other contributors
- Focus on the code, not the person

Thank you for contributing! 🎉
