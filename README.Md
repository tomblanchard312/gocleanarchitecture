# Go Clean Architecture Example

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
├── .vscode/
│   └── settings.json
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
├── main.go
├── main.exe
├── README.Md
└── server.log
```

## Why Clean Architecture?

Clean Architecture promotes separation of concerns, making your codebase easier to maintain and test. Here are some benefits:

- **Testability**: Each layer can be tested independently.
- **Maintainability**: Changes in one layer have minimal impact on other layers.
- **Scalability**: Easily extend functionality by adding new layers or components.
- **Flexibility**: Swap out frameworks or tools with minimal changes to the overall architecture.

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) 1.17 or higher

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
    go build -o main.exe
    ```

2. Run the application:

    ```sh
    ./main.exe
    ```

3. The server will start on port 8080. You can test the endpoints using curl or any API testing tool.

### API Endpoints

#### Create a Blog Post

```sh
curl -X POST http://localhost:8080/blogposts -H "Content-Type: application/json" -d '{"ID":"1", "Title":"Test Title", "Content":"Test Content"}'
```

### Get All Blog Posts

```sh
curl http://localhost:8080/blogposts
```

### Running Tests

 To run all tests:

```sh
go test ./...
```

### Logging

Logs are written to server.log for debugging purposes. Ensure the log file is writable and check it for any server-related errors.
License


## Error Handling

This project uses custom error types defined in the `errors` package. All errors are logged and appropriate HTTP status codes are returned to the client.

## Dependencies

- [Gorilla Mux](https://github.com/gorilla/mux) for HTTP routing
- [Viper](https://github.com/spf13/viper) for configuration management

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Future Improvements

- Add authentication and authorization
- Implement a persistent database
- Add more CRUD operations for blog posts

### Licensing

This project is licensed under the MIT License. See the LICENSE file for details.

## Acknowledgments

- [The Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) by Robert C. Martin
  
#
