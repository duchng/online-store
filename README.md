# Store Management System

#### Following Clean Architecture principles for separating concerns and ensuring maintainability.

## Features

- **User Management**
  - User authentication with JWT
  - Role-based access control (User/Admin roles)
  - User profile management
  - Password management

- **Product Management**
  - Product CRUD operations
  - Category management
  - Product reviews system
  - Wishlist functionality

## Technology Stack

- **Backend**
  - Go 1.24
  - Echo for HTTP routing
  - JWT for authentication, ed25519 for signing
  - Argon2 for password hashing
  - WebSocket for real-time communication
  - Swagger for API documentation
  - Test-containers for integration testing
  - Logging with slog - go's standard library structured logging

- **Database & Caching**
  - PostgreSQL 14 for persistent storage
  - Redis 6.2 for caching and pub/sub
  - Bun for ORM
  - golang-migrate for database migrations

- **Infrastructure**
  - Docker and Docker Compose for containerization
  - Clean Architecture
  - Dependency Injection

## Project Structure

```
.
├── cmd/ecommerce/          # Application entry point
├── docs/                   # Swagger documentation
├── internal/
│   ├── adapter/           # External interfaces implementation
│   │   ├── middlewares/   # HTTP middlewares
│   │   ├── product/      # Product adapters (HTTP, PostgreSQL)
│   │   ├── user/         # User adapters (HTTP, PostgreSQL, Redis)
│   │   └── server/       # HTTP server setup
│   ├── core/             # Business logic and entities
│   │   ├── product/      # Product domain
│   │   └── user/         # User domain
│   └── config/           # Application configuration
└── test/                 # Integration tests
```

## DB schema 
[db.png]

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.22 or later (for local development)

### Configuration

1. Configure the application:
   - Review `internal/assets/config.yaml` for application settings, this file is intended for local development, for other environments, use environment variables.
   - Environment variables in `docker-compose.yaml` for container setup

### Running with Docker

1. Start all services:
   ```bash
   docker-compose up -d
   ```

   This will start:
   - API server on port 8080
   - PostgreSQL on port 5432
   - Redis on port 6379

2. Access the API:
   - Swagger documentation: http://localhost:8080/swagger/index.html or `docs/swagger.json`
   - API endpoint: http://localhost:8080/

### Database Migrations

Migrations are automatically applied on startup. Migration files are located in:
```
internal/assets/migrations/
```

## API Documentation

Access the Swagger documentation at `http://localhost:8080/swagger/index.html` or `docs/swagger.json` for:
- Detailed API endpoints
- Request/Response schemas
- Authentication requirements
- Example requests

## Development

### Local Setup

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Run the application:
   ```bash
   go run cmd/ecommerce/main.go
   ```

### Testing

Run integration tests:
```bash
go test ./test/...
```

## Requirements (Senior Level):

- **Real-time Data Processing (Focus on real-time problem solving)**
    - I have implemented a websocket server to for receiving real-time data of user wishlist and product reviews.
    - The approach is to use Redis Pub/Sub for real-time notifications. Keep the stat data in Redis to avoid recalculating it every time.
    - Even when the system has a large number of writings, most read operations are performed on the cache, so the system can handle a large number of requests.
    - Considerations: For serious systems, we should consider the following:
        - Use specialized data structures like hyperloglog for counting metrics.
        - Use a more robust message broker (e.g., Kafka, RabbitMQ) for real-time data processing.
        - Consider change data capture (CDC) to decouple event processing from the main application.
- **Product Search Optimization (Focus on data search performance):**
    - I did not implement this feature, but I have some ideas:
        - Instead of auto-incrementing the ID, we can use UUIDs to avoid bottlenecks inserting millions of records. (or sortable ones like ulid, xid)
        - For the scope of this project, a dedicated search engine like Elasticsearch or Meili is not necessary. I'm thinking about an embedded indexing engine like Bleve.
- **Log Aggregation Optimization (Focus on efficient log handling):**
    - The error handling and logging is straightforward: only log at the top level handler like http handler.
    - Use error wrapping to provide context for errors (fmt.Errorf with %w for wrapping).
    - All errors are logged with a structured logger (slog), in middleware internal/adapter/middlewares/http.go:69 (CustomHTTPErrorHandler).