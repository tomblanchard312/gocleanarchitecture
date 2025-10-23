# Go Clean Architecture Blog - Complete Feature Set

This document provides a comprehensive overview of all implemented features in the Go Clean Architecture Blog application.

## 🎯 Implemented Features (100% Complete)

### ✅ 1. Role-Based Access Control (RBAC)

**Description**: Multi-level user permission system with admin and user roles.

**Features**:
- User and Admin roles with distinct permissions
- Admin middleware for protected routes
- Admin-only endpoints for user management
- Role-based authorization checks
- Database-level role enforcement (SQLite, Supabase, In-Memory)

**Admin Endpoints**:
- `GET /admin/users` - List all users
- `GET /admin/users/{id}` - Get user details
- `PUT /admin/users/{id}/role` - Update user role
- `DELETE /admin/users/{id}` - Delete user (cannot delete self)

**Database Schema**:
- `users.role` field with CHECK constraint (`user` or `admin`)
- Row-level security policies in Supabase
- Default role: `user`

---

### ✅ 2. Blog Post Comments System

**Description**: Nested comment system with support for replies, author tracking, and ownership validation.

**Features**:
- Create comments on blog posts
- Reply to existing comments (nested comments)
- Update own comments
- Delete own comments (or admin can delete any)
- Author tracking and ownership validation
- Automatic timestamps (created_at, updated_at)

**Endpoints**:
- `GET /blogposts/{id}/comments` - Get all comments for a blog post
- `POST /blogposts/{id}/comments` - Create a comment (protected)
- `GET /comments/{id}/replies` - Get replies to a comment
- `PUT /comments/{id}` - Update a comment (protected, author only)
- `DELETE /comments/{id}` - Delete a comment (protected, author or admin)

**Database Schema**:
- `comments` table with foreign keys to `blog_posts` and `users`
- `parent_id` for nested comments/replies
- Cascade delete when blog post or user is deleted
- Full support for SQLite, Supabase, and In-Memory

---

### ✅ 3. WebSocket Real-Time Updates

**Description**: Real-time notifications for new blog posts and comments using WebSocket connections.

**Features**:
- WebSocket hub for connection management
- Broadcast new blog posts to all connected clients
- Broadcast new comments to all connected clients
- Connection management (register/unregister)
- Ping/pong keepalive mechanism
- Message queuing and buffering
- Support for authenticated and anonymous connections

**Endpoint**:
- `GET /ws` - WebSocket connection endpoint

**Message Types**:
- `connection` - Connection confirmation
- `new_blog_post` - New blog post created
- `new_comment` - New comment created
- `error` - Error message

**Usage Example** (JavaScript):
```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  
  switch(message.type) {
    case 'connection':
      console.log('Connected to WebSocket');
      break;
    case 'new_blog_post':
      console.log('New blog post:', message.data);
      // Update UI with new blog post
      break;
    case 'new_comment':
      console.log('New comment:', message.data);
      // Update UI with new comment
      break;
  }
};
```

---

### ✅ 4. OAuth2 Social Login (Google & GitHub)

**Description**: Social login integration with Google and GitHub OAuth2 providers.

**Features**:
- Google OAuth2 login
- GitHub OAuth2 login
- Automatic user creation or linking
- Profile information import (name, email, avatar)
- Secure token exchange
- Configurable redirect URLs

**Configuration** (`.env`):
```env
# Google OAuth2
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback

# GitHub OAuth2
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
GITHUB_REDIRECT_URL=http://localhost:8080/auth/github/callback

BASE_URL=http://localhost:8080
```

**OAuth2 Providers**:
- **Google**: Uses `https://www.googleapis.com/oauth2/v2/userinfo`
- **GitHub**: Uses `https://api.github.com/user`

**Scopes**:
- Google: `userinfo.email`, `userinfo.profile`
- GitHub: `user:email`, `read:user`

---

## 🔐 Security Features

### Authentication & Authorization
- JWT-based authentication
- Bcrypt password hashing
- Token expiration and validation
- Protected routes with auth middleware
- Role-based access control
- Owner-only operations (update/delete own content)

### Database Security
- Row-level security (RLS) in Supabase
- Foreign key constraints
- Cascade deletes
- Input validation
- SQL injection prevention (parameterized queries)

---

## 📊 Database Support

The application supports three database backends:

### 1. SQLite (Default)
- Local file-based database
- Perfect for development and small deployments
- Auto-creates schema on startup
- Includes all tables: `users`, `blog_posts`, `comments`

### 2. Supabase (PostgreSQL)
- Cloud-hosted PostgreSQL
- Row-level security (RLS)
- Real-time capabilities
- Scalable for production
- Requires SQL schema execution

### 3. In-Memory
- Fast, non-persistent storage
- Great for testing
- Data lost on restart
- No setup required

---

## 🌐 Complete API Reference

### Public Endpoints

#### Authentication
- `POST /auth/register` - Register new user
- `POST /auth/login` - Login and get JWT token
- `GET /auth/users/{username}` - Get public user profile

#### Blog Posts (Read)
- `GET /blogposts` - Get all blog posts
- `GET /blogposts/{id}` - Get specific blog post

#### Comments (Read)
- `GET /blogposts/{id}/comments` - Get comments for a blog post
- `GET /comments/{id}/replies` - Get replies to a comment

#### WebSocket
- `GET /ws` - WebSocket connection (supports auth)

### Protected Endpoints (Require JWT)

#### User Profile
- `GET /auth/profile` - Get current user profile
- `PUT /auth/profile` - Update profile
- `POST /auth/change-password` - Change password

#### Blog Posts (Write)
- `POST /blogposts` - Create blog post
- `PUT /blogposts/{id}` - Update own blog post
- `DELETE /blogposts/{id}` - Delete own blog post

#### Comments (Write)
- `POST /blogposts/{id}/comments` - Create comment
- `PUT /comments/{id}` - Update own comment
- `DELETE /comments/{id}` - Delete own comment

### Admin-Only Endpoints (Require Admin Role)

- `GET /admin/users` - List all users
- `GET /admin/users/{id}` - Get user details
- `PUT /admin/users/{id}/role` - Change user role
- `DELETE /admin/users/{id}` - Delete user

---

## 🛠️ Technology Stack

### Core
- **Go 1.21+** - Programming language
- **Gorilla Mux** - HTTP router
- **Gorilla WebSocket** - WebSocket support

### Authentication & Security
- **JWT-Go** - JWT token generation/validation
- **Bcrypt** - Password hashing
- **OAuth2** - Social login

### Database
- **SQLite** - Embedded database
- **Supabase Go Client** - PostgreSQL cloud
- **In-Memory** - Testing

### Configuration
- **Viper** - Configuration management
- **GoDotEnv** - Environment variables

---

## 🧪 Testing

The project includes comprehensive tests for:
- Entity validation
- Use case business logic
- Repository operations
- Error handling

Run tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

---

## 📦 Deployment

### Option 1: Binary Deployment
```bash
go build -o blog-server ./cmd/main.go
./blog-server
```

### Option 2: Docker
```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o blog-server ./cmd/main.go
CMD ["./blog-server"]
```

### Option 3: Cloud Platform
- Deploy to AWS, Google Cloud, Azure, or any Go-supporting platform
- Set environment variables for configuration
- Use Supabase for database (no local DB management)

---

## 🎓 Architecture Highlights

### Clean Architecture Layers
1. **Entities** (`entities/`) - Core business logic
2. **Use Cases** (`usecases/`) - Application business rules
3. **Interface Adapters** (`interfaces/`) - Controllers, gateways
4. **Frameworks** (`frameworks/`) - External tools (DB, web, auth)

### Design Patterns
- Repository Pattern
- Dependency Injection
- Adapter Pattern
- Observer Pattern (WebSocket)
- Strategy Pattern (multiple DB backends)

### SOLID Principles
- Single Responsibility
- Open/Closed
- Liskov Substitution
- Interface Segregation
- Dependency Inversion

---

## 📝 Configuration Summary

### Required Environment Variables
```env
# Server
SERVER_PORT=:8080

# Database (choose one)
DB_TYPE=sqlite  # or "supabase" or "inmemory"
DB_PATH=./blog.db  # for SQLite

# Supabase (if DB_TYPE=supabase)
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_KEY=your-anon-key

# JWT
JWT_SECRET=your-secret-key
JWT_TOKEN_DURATION_HOURS=24

# OAuth2 (optional)
GOOGLE_CLIENT_ID=...
GOOGLE_CLIENT_SECRET=...
GITHUB_CLIENT_ID=...
GITHUB_CLIENT_SECRET=...
```

---

## 🚀 Quick Start

1. Clone repository
2. Copy `.env.example` to `.env` and configure
3. Run `go mod tidy`
4. Run `go run cmd/main.go`
5. Access `http://localhost:8080`

**Default credentials** (if using in-memory):
- Register first user → becomes admin
- Subsequent users → regular users

---

## 📈 Future Enhancements

Potential additions:
- File upload for images
- Full-text search
- Email notifications
- Rate limiting
- API versioning
- GraphQL support
- Caching layer (Redis)
- Metrics and monitoring

---

## 📄 License

MIT License - See LICENSE file for details.

---

## 👥 Contributing

Contributions are welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Make changes with tests
4. Submit a pull request

---

**Made with ❤️ using Go and Clean Architecture principles**

