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

- **CRUD operations for blog posts** with author tracking and ownership validation
- **User Authentication & Authorization**: JWT-based authentication with secure password hashing
- **Role-Based Access Control (RBAC)**: Admin and user roles with permission management
- **Comments System**: Hierarchical comments with replies on blog posts
- **Real-time Updates**: WebSocket support for live notifications of new posts and comments
- **OAuth2 Social Login**: Integration with Google and GitHub (framework ready)
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

- `POST /blogposts`: Create a new blog post (automatically assigned to authenticated user)
- `PUT /blogposts/{id}`: Update a blog post (only if you're the author)
- `DELETE /blogposts/{id}`: Delete a blog post (only if you're the author)

### Comment Endpoints (Public - Read Only)

- `GET /blogposts/{blogPostId}/comments`: Get all comments for a blog post
- `GET /comments/{commentId}/replies`: Get all replies to a specific comment

### Comment Endpoints (Protected - Requires JWT Token)

- `POST /blogposts/{blogPostId}/comments`: Create a new comment on a blog post
- `PUT /comments/{commentId}`: Update your own comment
- `DELETE /comments/{commentId}`: Delete your own comment

### Admin Endpoints (Protected - Requires Admin Role)

- `GET /admin/users`: List all users
- `GET /admin/users/{id}`: Get detailed user information
- `PUT /admin/users/{id}/role`: Update a user's role (admin/user)
- `DELETE /admin/users/{id}`: Delete a user account

### WebSocket Endpoint

- `GET /ws`: WebSocket connection for real-time updates
  - Broadcasts new blog posts when created
  - Broadcasts new comments when posted

### OAuth2 Endpoints (Coming Soon)

- `GET /auth/google`: Initiate Google OAuth login
- `GET /auth/google/callback`: Google OAuth callback
- `GET /auth/github`: Initiate GitHub OAuth login
- `GET /auth/github/callback`: GitHub OAuth callback

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
- [Gorilla WebSocket](https://github.com/gorilla/websocket) for WebSocket connections
- [Viper](https://github.com/spf13/viper) for configuration management
- [JWT-Go](https://github.com/golang-jwt/jwt) for JWT token generation and validation
- [bcrypt](https://golang.org/x/crypto/bcrypt) for secure password hashing
- [OAuth2](https://golang.org/x/oauth2) for social login integration
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

### OAuth2 Configuration (Optional - for social login)
- `BASE_URL`: Base URL for OAuth callbacks (default: "http://localhost:8080")
- `GOOGLE_CLIENT_ID`: Google OAuth2 Client ID
- `GOOGLE_CLIENT_SECRET`: Google OAuth2 Client Secret
- `GOOGLE_REDIRECT_URL`: Google OAuth callback URL (default: "http://localhost:8080/auth/google/callback")
- `GITHUB_CLIENT_ID`: GitHub OAuth2 Client ID
- `GITHUB_CLIENT_SECRET`: GitHub OAuth2 Client Secret
- `GITHUB_REDIRECT_URL`: GitHub OAuth callback URL (default: "http://localhost:8080/auth/github/callback")

## API Documentation

### Swagger/OpenAPI
The API is documented using OpenAPI 3.0 specification. You can find the complete API documentation in `api-documentation.yaml`.

To view the documentation:
1. **Using Swagger UI**: Upload the `api-documentation.yaml` file to [Swagger Editor](https://editor.swagger.io/)
2. **Using Swagger UI locally**:
   ```bash
   docker run -p 80:8080 -e SWAGGER_JSON=/api-documentation.yaml -v $(pwd):/usr/share/nginx/html swaggerapi/swagger-ui
   ```

### Postman Collection
Import the `Go-Clean-Architecture-Blog.postman_collection.json` file into Postman for easy API testing.

**Setup**:
1. Import the collection into Postman
2. The collection includes environment variables:
   - `base_url`: API server URL (default: http://localhost:8080)
   - `jwt_token`: Auto-saved after login/register
   - `user_id`: Auto-saved after authentication
   - `post_id`: Auto-saved after creating a post

### Feature Documentation
For detailed documentation on advanced features, see [`FEATURES.md`](FEATURES.md):
- **Role-Based Access Control (RBAC)**: Admin and user roles
- **Comments System**: Hierarchical comments with replies
- **WebSocket Real-time Updates**: Live notifications
- **OAuth2 Social Login**: Google and GitHub integration (framework)

## Deployment

### Prerequisites for Production

- Go 1.17 or higher
- Database (SQLite, Supabase, or other PostgreSQL-compatible)
- SSL/TLS certificates (recommended)
- Reverse proxy (nginx or similar)

### Deployment Options

#### Option 1: Deploy with SQLite (Simple)

1. **Build the application**:
   ```bash
   CGO_ENABLED=1 go build -o blog-api ./cmd/main.go
   ```

2. **Create a production `.env` file**:
   ```bash
   cp .env.example .env
   ```

3. **Configure production settings**:
   ```env
   DB_TYPE=sqlite
   DB_PATH=/var/lib/blog-api/blog.db
   SERVER_PORT=:8080
   LOG_LEVEL=info
   LOG_FILE=/var/log/blog-api/server.log
   JWT_SECRET=your-super-secret-key-minimum-32-characters-long
   JWT_TOKEN_DURATION_HOURS=24
   ```

4. **Run the application**:
   ```bash
   ./blog-api
   ```

5. **Set up as a systemd service** (Linux):
   
   Create `/etc/systemd/system/blog-api.service`:
   ```ini
   [Unit]
   Description=Go Clean Architecture Blog API
   After=network.target

   [Service]
   Type=simple
   User=blog-api
   WorkingDirectory=/opt/blog-api
   ExecStart=/opt/blog-api/blog-api
   Restart=on-failure
   RestartSec=10

   Environment="PATH=/usr/bin:/usr/local/bin"

   [Install]
   WantedBy=multi-user.target
   ```

   Enable and start the service:
   ```bash
   sudo systemctl enable blog-api
   sudo systemctl start blog-api
   sudo systemctl status blog-api
   ```

#### Option 2: Deploy with Supabase (Recommended for Production)

1. **Set up Supabase project**:
   - Create a project at [supabase.com](https://supabase.com)
   - Note your project URL and anon key
   - Run the SQL scripts in Supabase SQL Editor (in order):
     - `supabase_users_schema.sql` (users table with RBAC)
     - `supabase_schema.sql` (blog posts table)
     - `supabase_comments_schema.sql` (comments table)

2. **Build the application**:
   ```bash
   go build -o blog-api ./cmd/main.go
   ```

3. **Configure for Supabase**:
   ```env
   DB_TYPE=supabase
   SUPABASE_URL=https://your-project.supabase.co
   SUPABASE_KEY=your-anon-public-key
   SERVER_PORT=:8080
   LOG_LEVEL=info
   LOG_FILE=/var/log/blog-api/server.log
   JWT_SECRET=your-super-secret-key-minimum-32-characters-long
   JWT_TOKEN_DURATION_HOURS=24
   ```

4. **Deploy to a VPS or cloud platform**

#### Option 3: Docker Deployment

1. **Create a `Dockerfile`**:
   ```dockerfile
   FROM golang:1.21-alpine AS builder
   
   WORKDIR /app
   COPY go.mod go.sum ./
   RUN go mod download
   
   COPY . .
   RUN go build -o blog-api ./cmd/main.go
   
   FROM alpine:latest
   RUN apk --no-cache add ca-certificates
   
   WORKDIR /root/
   COPY --from=builder /app/blog-api .
   COPY --from=builder /app/.env.example .env
   
   EXPOSE 8080
   CMD ["./blog-api"]
   ```

2. **Create `docker-compose.yml`** (with SQLite):
   ```yaml
   version: '3.8'
   
   services:
     api:
       build: .
       ports:
         - "8080:8080"
       environment:
         - DB_TYPE=sqlite
         - DB_PATH=/data/blog.db
         - JWT_SECRET=your-secret-key
       volumes:
         - ./data:/data
         - ./logs:/var/log/blog-api
       restart: unless-stopped
   ```

3. **Build and run**:
   ```bash
   docker-compose up -d
   ```

### Nginx Reverse Proxy Configuration

Create `/etc/nginx/sites-available/blog-api`:

```nginx
server {
    listen 80;
    server_name yourdomain.com;
    
    # Redirect to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name yourdomain.com;
    
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
    
    # WebSocket support
    location /ws {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "Upgrade";
        proxy_set_header Host $host;
        proxy_read_timeout 86400;
    }
}
```

Enable the site:
```bash
sudo ln -s /etc/nginx/sites-available/blog-api /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### Security Checklist for Production

- [ ] Change `JWT_SECRET` to a strong, randomly generated secret
- [ ] Use HTTPS/TLS encryption
- [ ] Set appropriate CORS policies
- [ ] Enable rate limiting
- [ ] Use environment variables for all secrets
- [ ] Keep the `.env` file secure (never commit it)
- [ ] Regularly update dependencies
- [ ] Set up monitoring and logging
- [ ] Configure firewall rules
- [ ] Use a strong database password (if applicable)
- [ ] Enable Supabase Row Level Security (RLS)
- [ ] Implement backup strategy
- [ ] Set up health check endpoints

### Environment-Specific Configuration

**Development**:
```env
DB_TYPE=inmemory
LOG_LEVEL=debug
```

**Staging**:
```env
DB_TYPE=sqlite
DB_PATH=./staging.db
LOG_LEVEL=info
```

**Production**:
```env
DB_TYPE=supabase
LOG_LEVEL=warn
JWT_TOKEN_DURATION_HOURS=24
```

### Monitoring and Health Checks

Consider adding health check endpoints and integrating with:
- **Prometheus** for metrics
- **Grafana** for dashboards
- **ELK Stack** for log aggregation
- **Sentry** for error tracking

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Future Improvements

- ✅ ~~Add role-based access control (RBAC) for admin/user permissions~~ **DONE**
- ✅ ~~Implement blog post comments~~ **DONE**
- ✅ ~~Add WebSocket real-time updates~~ **DONE**
- ✅ ~~Add OAuth2 framework for social login~~ **DONE** (Google/GitHub ready)
- Implement password reset via email
- Add blog post likes and reactions
- Add pagination for blog post and comment listings
- Implement full-text search for blog posts
- Add user follow/unfollow functionality
- Implement rate limiting for API endpoints
- Add email notifications for comments and replies
- Add image upload support for blog posts and user avatars
- Implement comment moderation for admins
- Add analytics and metrics tracking

### Licensing

This project is licensed under the MIT License. See the LICENSE file for details.

## Acknowledgments

- [The Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) by Robert C. Martin
