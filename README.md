# Hotel Microservice Blueprint

A lightweight Go microservice built with a clean architecture pattern, featuring PostgreSQL integration, structured logging, and HTTP request handling via `chi` router.

## Architecture

The project follows a layered architecture:

```
cmd/api/main.go        → Entry point, wires dependencies
internal/handler       → HTTP handlers, routing, and middleware
internal/service       → Business logic layer
internal/repo          → Data access layer
internal/database      → Database connection management
internal/logging       → Structured logging setup
internal/models        → Domain models
internal/helper        → Utility functions
```

## Tech Stack

- **Router**: [go-chi/chi/v5](https://github.com/go-chi/chi)
- **Logging**: [go-chi/httplog/v3](https://github.com/go-chi/httplog) + `log/slog`
- **Database**: [jackc/pgx/v5](https://github.com/jackc/pgx) (PostgreSQL connection pool)

## Prerequisites

- Go 1.25.7+
- PostgreSQL database
- Docker & Docker Compose (optional, for local development)

## Getting Started

### 1. Set Environment Variables

```bash
export DATABASE_URL="postgres://user:password@localhost:5432/dbname?sslmode=disable"
```

### 2. Run the Service

```bash
go run app/cmd/api/main.go
```

The server starts on `localhost:8000`.

### 3. Test the Health Endpoint

```bash
curl http://localhost:8000/health
```

Response:
```json
{"status": "healthy"}
```

## Docker Compose

Use `docker-compose.yml` to spin up dependencies (e.g., PostgreSQL):

```bash
docker-compose up -d
```

## Project Structure

| Path | Description |
|------|-------------|
| `app/cmd/api/main.go` | Application entry point. Wires together database, repository, service, and handler layers, then starts the HTTP server. |
| `app/internal/database/` | Database connection pool initialization using `pgx`. |
| `app/internal/handler/` | HTTP handlers, request routing (`chi`), and middleware (security headers, request logging, cache control). |
| `app/internal/service/` | Business logic layer. Defines service interfaces and implements use cases. |
| `app/internal/repo/` | Data access layer. Handles all database queries and transactions. |
| `app/internal/logging/` | Structured JSON logger configuration using `slog` and `httplog`. |
| `app/internal/models/` | Domain models and data structures shared across layers. |
| `app/internal/helper/` | Utility/helper functions used across the application. |
| `app/sql/` | SQL migration files and queries. |
| `app/test/` | Test files. |
| `docker-compose.yml` | Docker Compose configuration for local development services. |

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Health check endpoint. Returns service health status with database connectivity check. |

## Adding New Features

1. **Models**: Define structs in `app/internal/models/models.go`
2. **Repository**: Add data access methods to `app/internal/repo/repo.go`
3. **Service**: Add business logic methods to `app/internal/service/service.go` (update the `Service` interface)
4. **Handler**: Add HTTP handler functions to `app/internal/handler/handlers.go`
5. **Routing**: Register new routes in `app/internal/handler/routing.go`

## Configuration

| Variable | Description |
|----------|-------------|
| `DATABASE_URL` | PostgreSQL connection string (required) |
