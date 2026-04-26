# 🔐 Users Microservice

> Authentication and user management service for the Hotel Reservation Platform.

## Overview

The Users Microservice is the **identity provider** for the entire platform. It handles user registration, authentication via JWT (RSA-256), profile management, and role-based access control. All other services trust the JWT tokens issued by this service.

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.25 |
| Router | [go-chi/chi](https://github.com/go-chi/chi) v5 |
| Database | PostgreSQL 16 |
| DB Driver | [pgx](https://github.com/jackc/pgx) v5 |
| Auth | JWT with RSA-256 (private/public key pair) |
| Password | bcrypt hashing |
| UUID | Google UUID v7 (time-sortable) |
| Container | Docker (multi-stage Alpine build) |

## Architecture

```
app/
├── cmd/api/          # Application entrypoint
│   └── main.go
├── internal/
│   ├── config/       # YAML config loader with env var expansion
│   ├── database/     # PostgreSQL connection pool (pgxpool)
│   ├── handler/      # HTTP handlers, routing, JWT middleware
│   ├── helper/       # Validators, error types, response helpers
│   ├── logging/      # Structured slog logger
│   ├── models/       # Domain entities (User, Login, DTOs)
│   ├── repo/         # Repository interface + PostgreSQL implementation
│   └── service/      # Business logic layer
├── sql/
│   └── migrations/   # golang-migrate compatible SQL migrations
├── config.yaml       # Service configuration
├── Dockerfile        # Multi-stage Docker build
└── go.mod
```

## API Endpoints

### Public Routes (No Authentication)

| Method | Path | Description |
|---|---|---|
| `GET` | `/health` | Liveness probe |
| `GET` | `/ready` | Readiness probe (checks DB) |
| `POST` | `/register` | Create a new user account |
| `POST` | `/login` | Authenticate and receive JWT token |

### Protected Routes (JWT Required)

| Method | Path | Description |
|---|---|---|
| `GET` | `/profile` | Get current user's profile |
| `PUT` | `/profile` | Update current user's display name |

## Data Model

### `users` Table

| Column | Type | Description |
|---|---|---|
| `id` | UUID v7 | Primary key |
| `email` | VARCHAR | Unique, indexed |
| `display_name` | VARCHAR(50) | User's display name |
| `user_type` | VARCHAR | `admin`, `user`, or `receptionist` |
| `is_active` | BOOLEAN | Account active flag |
| `created_at` | TIMESTAMP | Record creation time |
| `updated_at` | TIMESTAMP | Last update time |

### `logins` Table

| Column | Type | Description |
|---|---|---|
| `id` | UUID v7 | Primary key |
| `password_hash` | TEXT | bcrypt hash |
| `last_login` | TIMESTAMP | Last successful login |
| `failed_attempts` | INT | Failed login counter |
| `is_locked` | BOOLEAN | Account lock flag |
| `user_id` | UUID | FK → `users.id` |

## Flow Diagram

```mermaid
flowchart TD
    A["Client Request"] --> B{"Route Type?"}
    B -->|Public| C{"Endpoint?"}
    B -->|Protected| D["JWT Middleware"]
    
    C -->|POST /register| E["Decode UserCreation"]
    C -->|POST /login| F["Decode LoginInput"]
    
    E --> E1["Validate Input"]
    E1 --> E2["Generate UUID v7"]
    E2 --> E3["bcrypt Hash Password"]
    E3 --> E4["Insert User + Login (Transaction)"]
    E4 --> E5["201 Created"]
    
    F --> F1["Validate Credentials"]
    F1 --> F2["Query User by Email"]
    F2 --> F3{"Account Locked?"}
    F3 -->|Yes| F4["403 Forbidden"]
    F3 -->|No| F5["bcrypt Compare"]
    F5 -->|Match| F6["Generate JWT (RSA-256)"]
    F6 --> F7["Update Last Login"]
    F7 --> F8["Return access_token"]
    F5 -->|No Match| F9["401 Unauthorized"]
    
    D --> D1{"Token Valid?"}
    D1 -->|Yes| D2["Extract Claims → Context"]
    D1 -->|No| D3["401 Unauthorized"]
    
    D2 --> G{"Endpoint?"}
    G -->|GET /profile| H["Get User by ID"]
    G -->|PUT /profile| I["Update Display Name"]
    
    H --> H1["Return UserIO"]
    I --> I1["Validate Input"]
    I1 --> I2["Update in DB"]
    I2 --> I3["200 OK"]
```

## Use Case Diagram

```mermaid
graph LR
    subgraph Actors
        Guest["🧑 Guest"]
        User["👤 Authenticated User"]
        Admin["🔑 Admin"]
    end
    
    subgraph "Users Microservice"
        UC1["Register Account"]
        UC2["Login"]
        UC3["View Profile"]
        UC4["Update Profile"]
    end
    
    Guest --> UC1
    Guest --> UC2
    User --> UC3
    User --> UC4
    Admin --> UC3
    Admin --> UC4
```

## State Diagram

```mermaid
stateDiagram-v2
    [*] --> Unregistered
    Unregistered --> Active : POST /register
    Active --> Authenticated : POST /login (success)
    Authenticated --> Active : Token expires
    Active --> Locked : Too many failed attempts
    Locked --> Active : Admin unlocks
    Active --> [*] : Account deleted
    
    state Active {
        [*] --> Idle
        Idle --> ProfileViewing : GET /profile
        Idle --> ProfileEditing : PUT /profile
        ProfileViewing --> Idle
        ProfileEditing --> Idle
    }
```

## Package Diagram

```mermaid
graph TB
    subgraph "cmd/api"
        Main["main.go"]
    end
    
    subgraph "internal"
        subgraph "handler"
            Handlers["handlers.go"]
            Routing["routing.go"]
            Middleware["middleware.go (JWT)"]
        end
        
        subgraph "service"
            SVC["service.go"]
        end
        
        subgraph "repo"
            RepoIF["repo.go (interface)"]
            DBRepo["database_repo.go"]
        end
        
        subgraph "models"
            Models["models.go"]
        end
        
        subgraph "helper"
            Helper["validators, errors, responses"]
        end
        
        subgraph "config"
            Config["config.go"]
        end
        
        subgraph "database"
            DB["connection.go"]
        end
        
        subgraph "logging"
            Logger["logger.go"]
        end
    end
    
    Main --> Config
    Main --> Logger
    Main --> DB
    Main --> Handlers
    Main --> SVC
    
    Handlers --> SVC
    Handlers --> Helper
    Handlers --> Models
    Routing --> Handlers
    Routing --> Middleware
    
    SVC --> RepoIF
    SVC --> Models
    SVC --> Helper
    
    DBRepo -.->|implements| RepoIF
    DBRepo --> DB
    DBRepo --> Models
```

## Configuration

```yaml
server:
  host: "0.0.0.0"
  port: 8080

jwt:
  issuer: "users-service"
  expiration: "24h"

logging:
  level: "info"
  format: "json"
```

### Environment Variables

| Variable | Description |
|---|---|
| `DATABASE_URL` | PostgreSQL connection string |

### Volume Mounts (Docker)

| Host Path | Container Path | Description |
|---|---|---|
| `./keys/private.pem` | `/app/keys/private.pem` | JWT signing key (RSA private) |
| `./keys/public.pem` | `/app/keys/public.pem` | JWT verification key (RSA public) |

## Running Locally

```bash
# With Docker Compose (from Infra/)
docker compose up user-service user-db

# Direct (requires PostgreSQL)
DATABASE_URL=postgres://user:pass@localhost:5432/userdb go run ./app/cmd/api
```

## Port Mapping

| Context | Port |
|---|---|
| Internal (container) | `8080` |
| External (host) | `8081` |
| Database (host) | `5433` → `5432` |
