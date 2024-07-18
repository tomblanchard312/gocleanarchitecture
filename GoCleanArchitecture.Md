# Go Clean Architecture Project Structure (Sample)

This document outlines the clean architecture structure used in our Go project. The architecture is designed to separate concerns, maintain independence from frameworks, and enhance testability.

## Project Structure

```
.
├── .vscode/
│   └── settings.json
├── cmd/
│   └── main.go
├── config/
│   └── config.go
├── entities/
│   ├── blog_post.go
│   └── blog_post_test.go
├── errors/
│   ├── errors.go
│   └── errors_test.go
├── frameworks/
│   ├── db/
│   │   ├── in_memory_blog_post_repository.go
│   │   └── in_memory_blog_post_repository_test.go
│   ├── logger/
│   │   ├── logger.go
│   │   └── logger_test.go
│   └── web/
│       ├── middleware/
│       │   ├── logging.go
│       │   └── recovery.go
│       └── router.go
├── interfaces/
│   ├── blog_post_controller.go
│   ├── blog_post_controller_test.go
│   └── blog_post_repository.go
├── usecases/
│   ├── blog_post_usecase.go
│   └── blog_post_usecase_test.go
├── gocleanarchitecture.Md
├── go.mod
├── go.sum
├── LICENSE
├── README.Md
└── server.log
```

## Directory Explanations

### .vscode

- `settings.json`: VS Code specific settings for the project.

### config

- `config.go`: Defines configuration structures and functions for loading application settings.

### entities

- `blog_post.go`: Defines the BlogPost entity and its properties.
- `blog_post_test.go`: Contains unit tests for the BlogPost entity.

### errors

- `errors.go`: Defines custom error types and error handling utilities.
- `errors_test.go`: Contains unit tests for the custom error types and utilities.

### frameworks

The `frameworks` directory contains implementations of external tools and frameworks.

#### db

- `in_memory_blog_post_repository.go`: Implements the blog post repository interface using in-memory storage.
- `in_memory_blog_post_repository_test.go`: Contains unit tests for the in-memory repository.

#### logger

- `logger.go`: Implements a custom logging system.
- `logger_test.go`: Contains unit tests for the logger implementation.

#### web

- `router.go`: Sets up the HTTP router and defines the routes.
- `middleware/`:
  - `logging.go`: Implements logging middleware for HTTP requests.
  - `recovery.go`: Implements recovery middleware to handle panics.

### interfaces

- `blog_post_controller.go`: Handles HTTP requests and responses for blog post operations.
- `blog_post_controller_test.go`: Contains unit tests for the blog post controller.
- `blog_post_repository.go`: Defines the interface for data persistence operations.

### usecases

- `blog_post_usecase.go`: Defines the use cases for blog post operations (e.g., Create, Read, Update, Delete).
- `blog_post_usecase_test.go`: Contains unit tests for the blog post use cases.

### Root Directory Files

- `gocleanarchitecture.Md`: Documentation about the clean architecture implementation.
- `go.mod` and `go.sum`: Go module definition and checksum files.
- `LICENSE`: The license file for the project.
- `main.go`: The entry point of the application.
- `main.exe`: The compiled executable of the application.
- `README.Md`: Project readme file.
- `server.log`: Log file for the server.

## Clean Architecture Principles

This structure adheres to the following clean architecture principles:

1. **Independence of Frameworks**: The core business logic (entities and use cases) does not depend on external frameworks or tools.
2. **Testability**: Each component has its own test file, allowing for comprehensive unit testing.
3. **Independence of UI**: The business logic can be tested without the UI, database, web server, or any external element.
4. **Independence of Database**: The core entities can be used by different databases easily, as the repository is just an interface.
5. **Separation of Concerns**: Each directory has a specific responsibility, making the codebase more modular and maintainable.

## Dependency Flow

The dependency rule in clean architecture dictates that source code dependencies can only point inwards:

Entities ← Use Cases ← Interfaces ← Frameworks

This means that:

- Entities are independent and have no dependencies.
- Use Cases depend only on Entities.
- Interfaces depend on Use Cases and Entities.
- Frameworks depend on Interfaces, Use Cases, and Entities.

## Benefits

- **Maintainability**: The separation of concerns makes it easier to locate and fix issues.
- **Flexibility**: Outer layers like databases or web frameworks can be easily swapped without affecting the core business logic.
- **Testability**: Each component can be tested in isolation, leading to more robust and reliable code.
- **Independence**: The core business logic is not tied to any specific framework or tool, making it more portable and future-proof.

## Conclusion

This clean architecture structure provides a solid foundation for building scalable, maintainable, and testable Go applications. It separates concerns, promotes independence from external frameworks, and enhances overall code quality.
