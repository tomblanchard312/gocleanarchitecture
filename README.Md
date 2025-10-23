# Go Clean Architecture Example

[![Go Version](https://img.shields.io/github/go-mod/go-version/tomblanchard312/gocleanarchitecture)](https://github.com/tomblanchard312/gocleanarchitecture/blob/main/go.mod)
[![Last Commit](https://img.shields.io/github/last-commit/tomblanchard312/gocleanarchitecture)](https://github.com/tomblanchard312/gocleanarchitecture/commits/main)
[![Open Issues](https://img.shields.io/github/issues/tomblanchard312/gocleanarchitecture)](https://github.com/tomblanchard312/gocleanarchitecture/issues)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

This is an example project demonstrating Clean Architecture principles in Go. It is a simple blog application that allows you to create and retrieve blog posts. The project structure is designed to separate concerns and make the application easy to maintain and test.

## Project Structure

The project is organized into the following layers:

1. **Entities**: Core business logic.
2. **Use Cases (Interactors)**: Application-specific business rules.
3. **Interface Adapters**: Adapters for the user interface, database, web controllers, etc.
4. **Frameworks and Drivers**: Framework-specific code, like routing, databases, etc.

### Directory Structure

```
.
│   .gitignore
│   coverage
│   go.mod
│   go.sum
│   gocleanarchitecture.exe
│   GoCleanArchitecture.Md
│   LICENSE
│   README.Md
│
├───.vscode
│       settings.json
│
├───cmd
│       blog.db
│       blogposts.json
│       gocleanarchitecture.exe
│       main.go
│       server.log
│
├───config
│       config.go
│
├───entities
│       blog_post.go
│
├───errors
│       errors.go
│
├───frameworks
│   ├───db
│   │   │   in_memory_blog_post_repository.go
│   │   │
│   │   └───sqlite
│   │           blog_post_repository.go
│   │           sqlite.go
│   │
│   ├───logger
│   │       logger.go
│   │
│   └───web
│       │   router.go
│       │
│       └───middleware
│               logging.go
│               recovery.go
│
├───interfaces
│       blog_post_controller.go
│       blog_post_repository.go
│
├───tests
│   ├───entities
│   │       blog_post_test.go
│   │
│   ├───errors
│   │       errors_test.go
│   │
│   ├───frameworks
│   │   ├───db
│   │   │       in_memory_blog_post_repository_test.go
│   │   │       sqlite_blog_post_repository_test.go
│   │   │
│   │   └───logger
│   │           logger_test.go
│   │
│   ├───interfaces
│   │       blog_post_controller_test.go
│   │
│   └───usecases
│           blog_post_usecase_test.go
│
└───usecases
        blog_post_usecase.go
```

## Why Clean Architecture?

Clean Architecture promotes separation of concerns, making your codebase easier to maintain and test. Here are some benefits:

- **Testability**: Each layer can be tested independently.
- **Maintainability**: Changes in one layer have minimal impact on other layers.
- **Scalability**: Easily extend functionality by adding new layers or components.
- **Flexibility**: Swap out frameworks or tools with minimal changes to the overall architecture.

## Getting Started

## Features

- CRUD operations for blog posts
- **User Authentication & Authorization**: JWT-based authentication with secure password hashing
- Clean Architecture implementation
- **Multiple Database Backends**: SQLite, Supabase, In-Memory
- Custom error handling with domain validation
- Structured logging
- Middleware for logging, recovery, and authentication
- Rich domain model with business rules
- User profile management with avatar support

## Prerequisites

- [Go](https://golang.org/dl/) 1.17 or higher
- **Database Options**:
  - SQLite3 (default, no setup required)
  - Supabase account and project (for cloud database)
  - In-memory (for testing, no setup required)

### Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/tomblanchard312/gocleanarchitecture.git
    cd gocleanarchitecture
    ```

2. Install dependencies:

    ```sh
    go mod tidy
    ```

### Running the Application

1. Build the application:

    ```sh
    cd cmd
    go build -o gocleanarchitecture.exe ./cmd
    ```

2. Run the application:

    ```sh
    ./gocleanarchitecture.exe
    ```

3. The server will start on port 8080. You can test the endpoints using curl or any API testing tool.

## API Endpoints

### Authentication Endpoints (Public)

- `POST /auth/register`: Register a new user account
- `POST /auth/login`: Authenticate and receive a JWT token
- `GET /auth/users/{username}`: Get public user profile by username

### Authentication Endpoints (Protected - Requires JWT Token)

- `GET /auth/profile`: Get authenticated user's profile
- `PUT /auth/profile`: Update authenticated user's profile
- `POST /auth/change-password`: Change user password

### Blog Post Endpoints (Public - Read Only)

- `GET /blogposts`: Retrieve all blog posts
- `GET /blogposts/{id}`: Retrieve a specific blog post

### Blog Post Endpoints (Protected - Requires JWT Token)

- `POST /blogposts`: Create a new blog post
- `PUT /blogposts/{id}`: Update a blog post
- `DELETE /blogposts/{id}`: Delete a blog post

### Example: Register a User

```sh
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john_doe","email":"john@example.com","password":"securepass123","full_name":"John Doe"}'
```

### Example: Login

```sh
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email_or_username":"john@example.com","password":"securepass123"}'
```

Response:
```json
{
  "user": {
    "id": "uuid-here",
    "username": "john_doe",
    "email": "john@example.com",
    "full_name": "John Doe"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Example: Create a Blog Post (Protected)

```sh
curl -X POST http://localhost:8080/blogposts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"ID":"1", "Title":"Test Title", "Content":"Test Content"}'
```

### Example: Get All Blog Posts

```sh
curl http://localhost:8080/blogposts
```

### Running Tests

To run all tests:

```sh
go test ./...
```

To run tests for a specific package:

```sh
go test ./tests/entities
go test ./tests/usecases
go test ./tests/interfaces
go test ./tests/frameworks/db
go test ./tests/frameworks/logger
```

To run a specific test function:

```sh
go test ./tests/entities -run TestBlogPostCreation
go test ./tests/usecases -run TestCreateBlogPost
```

To run tests with verbose output:

```sh
go test -v ./tests/...
```

To run tests and see coverage:

```sh
go test ./tests/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Logging

Logs are written to server.log for debugging purposes. Ensure the log file is writable and check it for any server-related errors.
License

## Error Handling

This project uses custom error types defined in the `errors` package. All errors are logged and appropriate HTTP status codes are returned to the client.

## Dependencies

- [Gorilla Mux](https://github.com/gorilla/mux) for HTTP routing
- [Viper](https://github.com/spf13/viper) for configuration management
- [JWT-Go](https://github.com/golang-jwt/jwt) for JWT token generation and validation
- [bcrypt](https://golang.org/x/crypto/bcrypt) for secure password hashing
- [godotenv](https://github.com/joho/godotenv) for loading environment variables
- [Google UUID](https://github.com/google/uuid) for generating unique identifiers
- [Supabase Go Client](https://github.com/supabase-community/supabase-go) for Supabase integration
  
## Configuration

The application can be configured using environment variables. Copy `.env.example` to `.env` and configure:

### General Settings
- `SERVER_PORT`: The port on which the server will run (default: ":8080")
- `LOG_LEVEL`: The logging level (default: "info")
- `LOG_FILE`: The path to the log file (default: "server.log")

### Database Configuration
- `DB_TYPE`: Database type - "sqlite", "supabase", or "inmemory" (default: "sqlite")
- `DB_PATH`: The path to the SQLite database file (default: "./blog.db")

### Supabase Configuration (when DB_TYPE=supabase)
- `SUPABASE_URL`: Your Supabase project URL
- `SUPABASE_KEY`: Your Supabase anon/public API key
- `SUPABASE_USER`: Database username (optional, for direct access)
- `SUPABASE_PASS`: Database password (optional, for direct access)

### JWT Authentication
- `JWT_SECRET`: Secret key for JWT token signing (change in production!)
- `JWT_TOKEN_DURATION_HOURS`: Token expiration time in hours (default: 24)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Future Improvements

- Add role-based access control (RBAC) for admin/user permissions
- Implement password reset via email
- Add social login (OAuth2)
- Implement blog post comments and likes
- Add pagination for blog post listings
- Implement full-text search for blog posts
- Add user follow/unfollow functionality
- Implement rate limiting for API endpoints

### Licensing

This project is licensed under the MIT License. See the LICENSE file for details.

## Acknowledgments

- [The Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) by Robert C. Martin
