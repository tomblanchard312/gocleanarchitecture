# Copilot Instructions for Go Clean Architecture

This is a Go Clean Architecture blog API following Uncle Bob's principles. The project strictly enforces dependency inversion and separation of concerns.

## Architecture & Data Flow

**Dependency Direction (inward only):** Entities ← Use Cases ← Interfaces ← Frameworks

- **Entry Point:** `cmd/main.go` - Dependency injection and server startup
- **HTTP Layer:** `frameworks/web/router.go` + `interfaces/blog_post_controller.go`
- **Business Logic:** `usecases/blog_post_usecase.go`
- **Domain:** `entities/blog_post.go` (pure Go structs, no external deps)
- **Data Access:** Repository interface in `interfaces/`, implementations in `frameworks/db/`

## Key Conventions

**Repository Pattern:** Interface defined in `interfaces/blog_post_repository.go`, implemented by:

- `frameworks/db/in_memory_blog_post_repository.go` (for tests)
- `frameworks/db/sqlite/blog_post_repository.go` (production)

**Error Handling:** Use custom types from `errors/errors.go`:

```go
errors.New("message")           // Simple error
errors.Wrap(err, "context")     // Wrapped error
```

**Logging:** Structured logging via `frameworks/logger/logger.go`:

```go
logger.Error("message", logger.Field("key", value))
```

**Testing:** Mock repositories and loggers. See `tests/usecases/blog_post_usecase_test.go` for patterns.

## Development Workflows

**Build:** `go build ./cmd` (from project root)
**Run:** `go run ./cmd` or `./cmd/gocleanarchitecture.exe`
**Test:** `go test ./...` (runs all tests)
**Coverage:** `go test ./tests/... -coverprofile=coverage.out && go tool cover -html=coverage.out`

**Note:** SQLite tests require CGO. On Windows, use: `set CGO_ENABLED=1 && go test ./...`

## Adding Features

**New Repository Implementation:**

1. Implement `interfaces.BlogPostRepository` in `frameworks/db/<type>/`
2. Wire in `cmd/main.go` dependency injection
3. Add tests in `tests/frameworks/db/`

**New HTTP Endpoint:**

1. Add method to `interfaces/blog_post_controller.go`
2. Register route in `frameworks/web/router.go`
3. Add use case method if needed in `usecases/blog_post_usecase.go`

**New Entity/Domain Logic:**

1. Add to `entities/` (keep pure - no framework imports)
2. Update repository interface if persistence needed
3. Add use case methods for business rules

## Critical Files for Changes

- `cmd/main.go` - All dependency wiring happens here
- `interfaces/blog_post_repository.go` - Contract for all data implementations
- `frameworks/web/router.go` - HTTP route definitions
- `config/config.go` - Environment configuration (DB path, server port, logging)

## Project-Specific Notes

- Uses Gorilla Mux for routing, Viper for config, SQLite3 for persistence
- Tests use in-memory repository to avoid DB dependencies
- Middleware for logging/recovery in `frameworks/web/middleware/`
- Custom logger supports structured fields: `logger.Field("key", value)`
- BlogPost entity has string ID, Title, Content + timestamps
