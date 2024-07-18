# Go Clean Architecture Project Structure (Sample)

This document outlines the clean architecture structure used in our Go project. The architecture is designed to separate concerns, maintain independence from frameworks, and enhance testability.

## Project Structure

```
.
│   .gitignore
│   go.mod
│   go.sum
│   GoCleanArchitecture.Md
│   LICENSE
│   README.Md
│
├───.vscode
│       settings.json
│
├───cmd
│       cmd.exe
│       main.go
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
│   │       in_memory_blog_post_repository.go
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

## Directory Explanations

### Root Directory Files

- `.gitignore`: Specifies intentionally untracked files to ignore.
- `go.mod` and `go.sum`: Go module definition and checksum files.
- `GoCleanArchitecture.Md`: This documentation file.
- `LICENSE`: The license file for the project.
- `README.Md`: Project readme file.

### .vscode

- `settings.json`: VS Code specific settings for the project.

### cmd

- `main.go`: The entry point of the application.
- `cmd.exe`: The compiled executable of the application.

### config

- `config.go`: Defines configuration structures and functions for loading application settings.

### entities

- `blog_post.go`: Defines the BlogPost entity and its properties.

### errors

- `errors.go`: Defines custom error types and error handling utilities.

### frameworks

The `frameworks` directory contains implementations of external tools and frameworks.

#### db

- `in_memory_blog_post_repository.go`: Implements the blog post repository interface using in-memory storage.

#### logger

- `logger.go`: Implements a custom logging system.

#### web

- `router.go`: Sets up the HTTP router and defines the routes.
- `middleware/`:
  - `logging.go`: Implements logging middleware for HTTP requests.
  - `recovery.go`: Implements recovery middleware to handle panics.

### interfaces

- `blog_post_controller.go`: Handles HTTP requests and responses for blog post operations.
- `blog_post_repository.go`: Defines the interface for data persistence operations.

### usecases

- `blog_post_usecase.go`: Defines the use cases for blog post operations (e.g., Create, Read, Update, Delete).

### tests

This directory contains all the test files, mirroring the structure of the main application.

- `entities/blog_post_test.go`: Tests for the BlogPost entity.
- `errors/errors_test.go`: Tests for custom error types and utilities.
- `frameworks/db/in_memory_blog_post_repository_test.go`: Tests for the in-memory repository.
- `frameworks/logger/logger_test.go`: Tests for the logger implementation.
- `interfaces/blog_post_controller_test.go`: Tests for the blog post controller.
- `usecases/blog_post_usecase_test.go`: Tests for the blog post use cases.
  
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
