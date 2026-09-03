# ZZHAN Blog System

A full-stack blog system based on Go + Vue 3 with frontend-backend separation, supporting article management, categories & tags, comments & likes, user authentication, and more.

## 🚀 Tech Stack

### Backend (ZZHAN)

| Component | Technology |
|-----------|------------|
| Web Framework | Gin |
| ORM | GORM |
| Database | MySQL |
| Cache | Redis |
| Authentication | JWT |
| Logging | Zap + Lumberjack |
| Configuration | Viper |
| Captcha | base64Captcha |

### Frontend (ZZHAN_web)

- **User Site** (user): Vue 3 + TypeScript + Vite
- **Admin Panel** (admin): Vue 3 + TypeScript + Vite

## 📁 Project Structure

```
ZZHAN/
├── cmd/server/           # Application entry point
├── internal/
│   ├── api/              # Routes and controllers
│   │   ├── admin/        # Admin API endpoints
│   │   └── web/          # User-facing API endpoints
│   ├── app/              # Application initialization
│   ├── middleware/       # Middleware
│   ├── model/            # Data models
│   │   ├── dto/          # Data Transfer Objects
│   │   └── entity/       # Entity definitions
│   ├── repository/       # Data access layer
│   └── service/          # Business logic layer
├── pkg/                  # Public utility packages
│   ├── b64c/             # Captcha
│   ├── config/           # Configuration loading
│   ├── database/         # Database connections
│   ├── errors/           # Error handling
│   ├── jwt/              # JWT authentication
│   ├── logger/           # Logging utilities
│   ├── response/         # Response wrapper
│   ├── storage/          # Storage service
│   └── utils/            # Utility functions
├── config.yaml           # Configuration file
├── Makefile              # Build scripts
└── go.mod                # Go dependencies
```

## ✨ Features

### User Site
- 📝 Article browsing (Markdown rendering)
- 🔍 Article search
- 📂 Category filtering
- 🏷️ Tag filtering
- 👍 Like/Unlike
- 💬 Comments/Replies
- 📊 Site statistics
- 📅 Article archive
- 🔐 User registration/Login
- 👤 User profile

### Admin Panel
- 📊 Dashboard
- 📝 Article management (Publish/Edit/Unpublish)
- 📂 Category management
- 🏷️ Tag management
- 💬 Comment management
- 👥 User management
- ⚙️ Site settings
- 📋 Operation logs

## 🛠️ Quick Start

### Prerequisites

- Go 1.24+
- Node.js 18+
- MySQL 8.0+
- Redis 7.0+

### 1. Clone the Repository

```bash
git clone <repository-url>
cd ZZHAN
```

### 2. Configure Backend

```bash
# Copy example config
cp .env.example .env

# Edit configuration file with your database and Redis connection info
vim config.yaml
```

### 3. Initialize Database

```bash
# Import database schema
mysql -u root -p < zzhan_dump.sql
```

### 4. Start Backend

```bash
# Using Makefile
make run

# Or run directly
go run cmd/server/main.go

# With custom config file
go run cmd/server/main.go -c config.yaml
```

### 5. Start Frontend

```bash
# User site
cd ZZHAN_web/user
npm install
npm run dev

# Admin panel
cd ZZHAN_web/admin
npm install
npm run dev
```

## 📦 Docker Deployment

```bash
# Build and start all services
docker-compose up -d

# View logs
docker-compose logs -f
```

## 🔧 Makefile Commands

| Command | Description |
|---------|-------------|
| `make build` | Build the project |
| `make run` | Run the project |
| `make test` | Run tests |
| `make test-coverage` | Generate test coverage report |
| `make clean` | Clean build artifacts |
| `make deps` | Download dependencies |
| `make fmt` | Format code |
| `make vet` | Static analysis |
| `make lint` | Code linting (requires golangci-lint) |
| `make help` | Show help information |

## 📡 API Endpoints

### User API

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/articles` | Article list |
| GET | `/api/v1/articles/:slug` | Article detail |
| GET | `/api/v1/categories` | Category list |
| GET | `/api/v1/tags` | Tag list |
| GET | `/api/v1/archives` | Archive list |
| GET | `/api/v1/stats` | Site statistics |
| POST | `/api/v1/auth/register` | User registration |
| POST | `/api/v1/auth/login` | User login |
| POST | `/api/v1/articles/:id/like` | Like article |
| POST | `/api/v1/articles/:id/comments` | Post comment |

### Admin API

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/admin/auth/login` | Admin login |
| GET | `/api/admin/dashboard` | Dashboard data |
| GET/POST/PUT/DELETE | `/api/admin/articles` | Article management |
| GET/POST/PUT/DELETE | `/api/admin/categories` | Category management |
| GET/POST/PUT/DELETE | `/api/admin/tags` | Tag management |
| GET/PUT/DELETE | `/api/admin/comments` | Comment management |
| GET/PUT | `/api/admin/users` | User management |
| GET/PUT | `/api/admin/settings` | Site settings |

## 🔐 Authentication

- JWT-based authentication
- Access Token + Refresh Token dual-token mechanism
- Single device login (only one active session per user)
- Token blacklist (invalidated on logout)

## 📊 View Count Statistics

Uses **IP-based deduplication**:
- Same IP counts only once within 10 minutes
- Access records stored in Redis
- Falls back to direct counting when Redis is unavailable

## 📝 License

MIT License

## 👥 Author

- 23107

## 🙏 Acknowledgments

Thanks to everyone who has contributed to this project.
