# Task Management Web Service

This project is a Task Management Web Service built in Go, providing a RESTful API for managing tasks. It supports CRUD operations, allowing users to create, retrieve, update, and delete tasks with user-specific access control.

## Features

- Create, Read, Update, and Delete tasks with user isolation
- Middleware for authentication (using User-Id header) and logging
- SQLite as a local database for data storage
- User-specific data handling and access control
- Structured project layout following best practices

## Project Structure

```
task-manager
├── cmd
│   └── server
│       └── main.go          # Entry point of the application
├── internal
│   ├── api
│   │   ├── handler
│   │   │   └── task_handler.go  # HTTP request handlers for tasks
│   │   ├── middleware
│   │   │   ├── auth.go         # Authentication middleware
│   │   │   └── logging.go      # Logging middleware
│   │   └── router.go          # HTTP routes setup
│   ├── config
│   │   └── config.go          # Configuration settings
│   ├── model
│   │   └── task.go            # Task entity definition
│   ├── repository
│   │   ├── sqlite
│   │   │   ├── db.go          # SQLite database connection
│   │   │   └── task_repository.go  # SQLite task repository
│   │   └── repository.go       # Task repository interface
│   └── service
│       └── task_service.go     # Business logic for task management
├── migrations
│   ├── 001_create_tasks_table.down.sql  # Rollback migration for tasks table
│   └── 001_create_tasks_table.up.sql    # Migration to create tasks table
├── pkg
│   └── utils
│       └── logger.go            # Logger utility
├── go.mod                        # Module definition
├── go.sum                        # Module dependency checksums
├── run.sh                        # Script to build and run the server
├── test_user_api.sh              # Script to test the API endpoints
└── README.md                     # Project documentation
```

## Setup Instructions

1. Make sure you have Go installed (version 1.16 or newer)

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run the application:
   ```bash
   ./run.sh
   ```
   
   This script will:
   - Create the database schema using SQLite
   - Build the application
   - Start the server

4. Test the API endpoints:
   ```bash
   ./test_user_api.sh
   ```
   
   This script sends requests with the required `User-Id` header for authentication and tests all API endpoints.

## Configuration

The application can be configured using environment variables:

| Variable        | Description                            | Default           |
|----------------|----------------------------------------|-------------------|
| PORT           | The port the server listens on          | 8080              |
| DATABASE_PATH  | Path to the SQLite database file        | data/tasks.db     |
| LOG_LEVEL      | Logging level (info, debug, error)      | info              |

Example usage:
```bash
PORT=9090 LOG_LEVEL=debug ./run.sh
```

## API Endpoints

| Method | URL          | Description                            | Required Headers         |
|--------|--------------|----------------------------------------|--------------------------|
| POST   | /tasks       | Create a new task                      | User-Id: [user_id]        |
| GET    | /tasks       | Get all tasks for the authenticated user| User-Id: [user_id]        |
| GET    | /tasks/:id   | Get a specific task by ID              | User-Id: [user_id]        |
| PUT    | /tasks/:id   | Update a task by ID                    | User-Id: [user_id]        |
| DELETE | /tasks/:id   | Delete a task by ID                    | User-Id: [user_id]        |

**Note:** The `User-Id` header is required for all requests to authenticate the user and ensure they can only access their own tasks.

## Task Model

```json
{
  "id": 1,
  "title": "Task Title",
  "description": "Task description",
  "due_date": "2025-07-01T00:00:00Z",
  "status": "Todo",
  "user_id": "user123",
  "created_at": "2025-06-23T10:00:00Z",
  "updated_at": "2025-06-23T10:00:00Z"
}
```

## Status Values

The following status values are supported for tasks:

- `Todo` - Task is not started
- `InProgress` - Task is being worked on
- `Done` - Task is completed

**Note:** Status values must be used exactly as shown (no spaces). The API will validate these values when creating or updating tasks.

## Usage Examples

Here are some examples using curl to interact with the API:

### Create a new task
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -H "User-Id: user123" \
  -d '{"title": "Complete project", "description": "Finish the task manager API", "due_date": "2025-07-01T00:00:00Z", "status": "Todo"}'
```

### Get all tasks for a user
```bash
curl -X GET http://localhost:8080/tasks \
  -H "User-Id: user123"
```

### Get a specific task by ID
```bash
curl -X GET http://localhost:8080/tasks/1 \
  -H "User-Id: user123"
```

### Update a task
```bash
curl -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -H "User-Id: user123" \
  -d '{"title": "Complete project", "description": "Finish the task manager API", "due_date": "2025-07-01T00:00:00Z", "status": "InProgress"}'
```

### Delete a task
```bash
curl -X DELETE http://localhost:8080/tasks/1 \
  -H "User-Id: user123"
```

## Database Schema

The application uses SQLite for data storage. The database schema is defined in the migration files:

```sql
CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,   
    user_id INTEGER NOT NULL,
    due_date DATETIME,
    status TEXT DEFAULT 'Todo',
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_tasks_due_date ON tasks(due_date);
CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks(user_id);
```

The database is automatically created when you run the application using the `run.sh` script, which executes the migration file.

## Development

### Project Structure

The project follows a layered architecture:

1. **API Layer** (`internal/api`) - Handles HTTP requests and responses
   - `handler` - Contains HTTP handlers for different resources
   - `middleware` - Contains middleware functions (auth, logging)
   - `router.go` - Sets up the routes and middleware

2. **Service Layer** (`internal/service`) - Contains business logic
   - `task_service.go` - Implements task management operations

3. **Repository Layer** (`internal/repository`) - Data access layer
   - `repository.go` - Defines interfaces for data access
   - `sqlite` - Contains SQLite implementation of the repository interfaces

4. **Model Layer** (`internal/model`) - Contains domain models
   - `task.go` - Defines the Task model and related constants