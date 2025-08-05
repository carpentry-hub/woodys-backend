# Woody's Backend API

A clean, layered architecture backend API for a woodworking project sharing platform built with Go, GORM, and PostgreSQL.

## 🏗️ Architecture

This project follows **Clean Architecture** principles with clear separation of concerns:

```
woodys-backend/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── database/        # Database connection & migrations
│   ├── domain/          # Domain models & business entities
│   ├── handlers/        # HTTP handlers (presentation layer)
│   ├── middleware/      # HTTP middleware
│   ├── repositories/    # Data access layer
│   └── services/        # Business logic layer
├── .env.example         # Environment variables template
├── go.mod              # Go module dependencies
├── go.sum              # Go module checksums
└── README.md           # This file
```

### Layers

- **Presentation Layer** (`handlers/`): HTTP request/response handling
- **Business Logic Layer** (`services/`): Core business logic and validation
- **Data Access Layer** (`repositories/`): Database operations and queries
- **Domain Layer** (`domain/`): Business entities and domain models

## 🚀 Features

- **Clean Architecture**: Proper separation of concerns
- **Environment Configuration**: Environment variables for all settings
- **Database Migrations**: Automatic database schema management
- **Middleware Support**: CORS, logging, error handling, rate limiting
- **Structured Logging**: Request tracking with unique IDs
- **Input Validation**: Comprehensive request validation
- **Error Handling**: Standardized error responses
- **Graceful Shutdown**: Proper resource cleanup on shutdown

## 📋 Prerequisites

- Go 1.24.4 or higher
- PostgreSQL database
- Environment variables configured

## 🛠️ Installation & Setup

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd woodys-backend
   ```

2. **Install dependencies**

   ```bash
   go mod download
   ```

3. **Configure environment variables**

   ```bash
   cp .env.example .env
   # Edit .env with your actual values
   ```

4. **Set environment variables**

   ```bash
   export DB_HOST=your-db-host
   export DB_USER=your-db-user
   export DB_PASSWORD=your-db-password
   export DB_NAME=your-db-name
   export DB_PORT=5432
   export DB_SSL_MODE=require
   export SERVER_PORT=8080
   ```

5. **Run the application**
   ```bash
   go run cmd/server/main.go
   ```

## 🌍 Environment Variables

| Variable      | Description                | Default   | Required |
| ------------- | -------------------------- | --------- | -------- |
| `DB_HOST`     | Database host              | localhost | Yes      |
| `DB_USER`     | Database username          | postgres  | Yes      |
| `DB_PASSWORD` | Database password          |           | Yes      |
| `DB_NAME`     | Database name              | postgres  | Yes      |
| `DB_PORT`     | Database port              | 5432      | No       |
| `DB_SSL_MODE` | SSL mode (require/disable) | disable   | No       |
| `SERVER_PORT` | Server port                | 8080      | No       |

## 🔗 API Endpoints

### Health & Info

- `GET /` - API information
- `GET /health` - Health check

### Users

- `POST /api/v1/users` - Create user
- `GET /api/v1/users/{id}` - Get user by ID
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user
- `GET /api/v1/users/{id}/projects` - Get user's projects
- `GET /api/v1/users/uid/{firebase_uid}` - Get user by Firebase UID

### Projects

- `POST /api/v1/projects` - Create project
- `GET /api/v1/projects/{id}` - Get project by ID
- `PUT /api/v1/projects/{id}` - Update project
- `DELETE /api/v1/projects/{id}` - Delete project
- `GET /api/v1/projects/search` - Search projects

### Comments

- `GET /api/v1/projects/{project_id}/comments` - Get project comments
- `POST /api/v1/projects/{project_id}/comments` - Create comment
- `DELETE /api/v1/comments/{id}` - Delete comment
- `GET /api/v1/comments/{id}/replies` - Get comment replies
- `POST /api/v1/comments/{id}/reply` - Create reply

### Ratings

- `POST /api/v1/projects/{project_id}/ratings` - Create/update rating
- `PUT /api/v1/projects/{project_id}/ratings` - Update rating
- `GET /api/v1/projects/{project_id}/ratings` - Get project ratings

### Project Lists

- `POST /api/v1/project-lists` - Create project list
- `GET /api/v1/project-lists/{id}` - Get project list
- `PUT /api/v1/project-lists/{id}` - Update project list
- `DELETE /api/v1/project-lists/{id}` - Delete project list
- `POST /api/v1/project-lists/{id}/projects` - Add project to list
- `DELETE /api/v1/project-lists/{list_id}/projects/{project_id}` - Remove project from list
- `GET /api/v1/users/{user_id}/project-lists` - Get user's project lists

## 📊 Database Schema

The application uses PostgreSQL with the following main entities:

- **Users**: User accounts with Firebase authentication
- **Projects**: Woodworking projects with materials, tools, styles
- **Comments**: Hierarchical comments and replies on projects
- **Ratings**: User ratings for projects (1-5 stars)
- **ProjectLists**: User-created collections of projects
- **ProjectListItems**: Join table for projects in lists
