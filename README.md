# High Stakes REST API

A clean architecture REST API built with Go, using standard library packages for HTTP and SQL, with PostgreSQL as the database.

## Project Structure

```
.
├── bin/                    # Compiled executables
├── cmd/
│   └── http/              # Main HTTP server entry point
├── docs/                   # Project documentation
└── internal/
    ├── adapter/            # External adapters
    │   ├── cache/redis/    # Redis cache adapter
    │   ├── handler/http/   # HTTP request/response handlers
    │   ├── repository/postgres/  # PostgreSQL database adapter
    │   │   └── migrations/ # Database migration files
    │   └── token/paseto/   # Paseto token adapter
    └── core/               # Core business logic
        ├── domain/         # Domain models/entities
        ├── port/           # Interface definitions
        ├── service/        # Business logic services
        └── util/           # Utility functions
```

## Features

- **User Management**: Full CRUD operations for users
- **Password Hashing**: Secure password storage using bcrypt
- **KSUID Generation**: Unique identifier generation for users
- **Clean Architecture**: Separation of concerns with ports and adapters
- **PostgreSQL**: Robust database support with migrations
- **Standard Library**: Uses only Go standard library for HTTP and SQL

## Prerequisites

- Go 1.25.4 or higher
- PostgreSQL database
- [Air](https://github.com/air-verse/air) (for hot reloading during development)
- [golang-migrate](https://github.com/golang-migrate/migrate) (for database migrations)
- Docker and Docker Compose (optional, for local development)

## Setup

1. **Install dependencies:**
   ```bash
   go mod download
   ```

2. **Install development tools:**
   ```bash
   # Install Air for hot reloading
   go install github.com/air-verse/air@latest
   
   # Install golang-migrate for database migrations
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   ```

3. **Set up environment variables:**
   Create a `.env` file in the root directory. If using Docker Compose (recommended for local development), use these values:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=high_stakes
   DB_SSL_MODE=disable
   ```
   
   These values match the `docker-compose.local.yml` configuration. If you're using a different PostgreSQL setup, adjust accordingly.

4. **Start local services with Docker Compose (recommended):**
   ```bash
   make local
   ```
   
   This will start PostgreSQL and Redis containers. The services will be available at:
   - PostgreSQL: `localhost:5432`
   - Redis: `localhost:6379`
   
   To stop the services:
   ```bash
   make local-down
   ```
   
   To view logs:
   ```bash
   make local-logs
   ```

   **Alternative:** If you prefer to set up PostgreSQL manually:
   ```bash
   createdb high_stakes
   ```

5. **Run database migrations:**
   ```bash
   make migrate_up
   ```
   
   **Note:** Make sure your `.env` file is configured and the database is running before running migrations.

## Development

### Using Air (Hot Reload)

Air provides automatic code reloading during development. Simply run:

```bash
air
```

Air will:
- Watch for changes in `.go` files
- Automatically rebuild and restart the server
- Display colored logs for different components
- Store build artifacts in the `tmp/` directory

The Air configuration is in `.air.toml`. You can customize:
- Build command
- Environment variables
- File watching patterns
- Log output

### Using Makefile

The project includes a Makefile with common development tasks:

#### Database Migrations

```bash
# Create a new migration
make create_migration MIGRATION_NAME=add_users_table

# Run all pending migrations
make migrate_up

# Rollback the last migration
make migrate_down

# Show migration help
make help
```

#### Running the Application

```bash
# Run the application directly
make run

# Build the binary
make build

# Run tests with coverage
make test
```

#### Docker Development

```bash
# Start local development environment with Docker Compose
# This starts PostgreSQL and Redis containers
make local

# Stop the containers
make local-down

# View container logs
make local-logs
```

The `docker-compose.local.yml` file configures:
- **PostgreSQL** on port `5432` with:
  - User: `postgres`
  - Password: `postgres`
  - Database: `high_stakes`
- **Redis** on port `6379`

Make sure your `.env` file matches these credentials for migrations and the application to work correctly.

#### Cleanup

```bash
# Clean Docker system
make clean
```

**Note:** The Makefile references `./cmd/api/main.go`, but the actual entry point is `./cmd/http/main.go`. Update the Makefile if needed, or use Air for development.

### Manual Server Execution

Run the server manually:

```bash
go run cmd/http/main.go
```

Or with custom configuration:
```bash
go run cmd/http/main.go \
  -port=8080 \
  -db-host=localhost \
  -db-port=5432 \
  -db-user=postgres \
  -db-password=yourpassword \
  -db-name=high_stakes \
  -db-sslmode=disable
```

### Building the Binary

```bash
go build -o bin/http cmd/http/main.go
```

## API Endpoints

### Health Check
- **GET** `/health`
  - Returns server status

### Users

- **POST** `/users`
  - Create a new user
  - Request body:
    ```json
    {
      "name": "John Doe",
      "email": "john@example.com",
      "phone": "+1234567890",
      "password": "securepassword"
    }
    ```

- **GET** `/users`
  - List all users with pagination
  - Query parameters:
    - `limit` (default: 10): Number of users to return
    - `offset` (default: 0): Number of users to skip
  - Example: `/users?limit=20&offset=0`

- **GET** `/users/:id`
  - Get a user by ID

- **PUT** `/users/:id`
  - Update a user
  - Request body:
    ```json
    {
      "name": "Jane Doe",
      "email": "jane@example.com",
      "phone": "+0987654321"
    }
    ```

- **DELETE** `/users/:id`
  - Delete a user

## Database Schema

### Users Table

```sql
CREATE TABLE users (
    id VARCHAR(27) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(50) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

## Example Usage

### Create a User
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+1234567890",
    "password": "mypassword"
  }'
```

### Get User by ID
```bash
curl http://localhost:8080/users/{user_id}
```

### List Users
```bash
curl http://localhost:8080/users?limit=10&offset=0
```

### Update User
```bash
curl -X PUT http://localhost:8080/users/{user_id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe",
    "email": "jane@example.com"
  }'
```

### Delete User
```bash
curl -X DELETE http://localhost:8080/users/{user_id}
```

## Architecture

This project follows **Clean Architecture** / **Hexagonal Architecture** principles:

- **Domain Layer**: Core business entities (User)
- **Port Layer**: Interfaces defining contracts (UserRepository, UserService)
- **Service Layer**: Business logic implementation
- **Adapter Layer**: External implementations (HTTP handlers, PostgreSQL repository)

This structure ensures:
- Business logic is independent of external frameworks
- Easy to test and maintain
- Easy to swap implementations (e.g., change database or HTTP framework)

## Development Tools

### Air Configuration

Air is configured via `.air.toml` with the following settings:
- **Build output**: `./tmp/main`
- **Watched extensions**: `.go`, `.tpl`, `.tmpl`, `.html`
- **Excluded directories**: `assets`, `tmp`, `vendor`, `frontend/node_modules`
- **Environment variables**: `APP_ENV=dev`, `APP_USER=air`
- **Auto-cleanup**: Removes `tmp/` directory on exit

### Makefile Commands

| Command | Description |
|---------|-------------|
| `make create_migration MIGRATION_NAME=name` | Create a new database migration |
| `make migrate_up` | Run all pending migrations |
| `make migrate_down` | Rollback the last migration |
| `make local` | Start Docker Compose environment (PostgreSQL + Redis) |
| `make local-down` | Stop Docker Compose containers |
| `make local-logs` | View Docker Compose logs |
| `make run` | Run the application |
| `make build` | Build the binary |
| `make test` | Run tests with coverage |
| `make clean` | Clean Docker system |
| `make help` | Show migration help |

**Important**: The Makefile requires a `.env` file with database configuration variables.

## Dependencies

### Runtime Dependencies
- `github.com/lib/pq`: PostgreSQL driver
- `golang.org/x/crypto/bcrypt`: Password hashing

### Development Dependencies
- `github.com/air-verse/air`: Hot reloading tool
- `github.com/golang-migrate/migrate`: Database migration tool

## License

MIT


