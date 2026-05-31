# RepJot

RepJot is a lightweight REST API written in Go for tracking workouts and exercise routines. It provides user authentication, workout history logs, and exercise entry tracking (sets, reps, weights, durations). 

---

## Technical Stack

* **Go 1.26**
* **Routing**: [go-chi/chi](https://github.com/go-chi/chi) v5 (for lightweight, radix-tree based routing)
* **Database**: PostgreSQL (driver: [jackc/pgx](https://github.com/jackc/pgx) v5 with `database/sql` compatibility)
* **Migrations**: [goose](https://github.com/pressly/goose) (SQL-based migrations)
* **Dev Setup**: Docker Compose (multi-db layout for separate development and testing environments)
* **Quality Assurance**: Husky and commitlint for conventional commit standards

---

## Directory Layout

* `main.go`: Application entrypoint, CLI flag parsing, and HTTP server configurations.
* `internal/app/`: Core application container holding database pools, logging utilities, and route handler mappings.
* `internal/api/`: HTTP handlers and request decoding/validation logic.
* `internal/store/`: Data Access Object (DAO) layer interface and Postgres-backed CRUD implementations.
* `internal/middleware/`: CORS, authentication token verification, and handler protection filters.
* `internal/tokens/`: Cryptographic utilities for session and authentication token generation.
* `migrations/`: Schema definition files managed by Goose.

---

## Quickstart

### 1. Requirements
Ensure you have Go, Docker, and Goose installed on your machine.

### 2. Run the Databases
Spin up both the local development database (`port 5432`) and the test database (`port 5433`):
```bash
docker compose up -d
```

### 3. Environment Setup
Configure your environment variables in a `.env` file at the root:
```env
DATABASE_URL="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
TEST_DATABASE_URL="postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable"
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://postgres:postgres@localhost:5432/postgres
GOOSE_MIGRATION_DIR=./migrations
```

### 4. Run Schema Migrations
Apply the migrations to your local development database:
```bash
goose up
```

### 5. Start the Application
By default, the server listens on port `8080`. Run the binary or source directly:
```bash
go run main.go --port=8080
```

---

## API Endpoints

All protected endpoints require a Bearer token supplied in the `Authorization` header.

### Authentication

#### Register a New User
```bash
curl -X POST http://localhost:8080/users \
     -H "Content-Type: application/json" \
     -d '{
           "username": "john_doe",
           "email": "john.doe@example.com",
           "password": "securepassword123",
           "bio": "Developer and runner"
         }'
```

#### Generate Auth Token
```bash
curl -X POST http://localhost:8080/tokens/authentication \
     -H "Content-Type: application/json" \
     -d '{
           "username": "john_doe",
           "password": "securepassword123"
         }'
```
*Response returns a raw token string to pass in the `Authorization: Bearer <token>` header.*

---

### Workouts (Auth Required)

#### Log a New Workout
```bash
curl -X POST http://localhost:8080/workouts \
     -H "Authorization: Bearer YOUR_TOKEN_HERE" \
     -H "Content-Type: application/json" \
     -d '{
           "title": "Leg Day",
           "description": "Lower body strength training",
           "duration_minutes": 60,
           "calories_burned": 400,
           "entries": [
             {
               "exercise_name": "Squats",
               "sets": 3,
               "reps": 12,
               "weight": 100.5,
               "notes": "Focused on depth",
               "order_index": 1
             }
           ]
         }'
```

#### Retrieve a Workout by ID
```bash
curl -H "Authorization: Bearer YOUR_TOKEN_HERE" \
     http://localhost:8080/workouts/1
```

#### Update a Workout
```bash
curl -X PUT http://localhost:8080/workouts/1 \
     -H "Authorization: Bearer YOUR_TOKEN_HERE" \
     -H "Content-Type: application/json" \
     -d '{
           "title": "Updated Leg Day",
           "description": "Lower body strength training session",
           "duration_minutes": 70,
           "calories_burned": 450,
           "entries": [
             {
               "exercise_name": "Squats",
               "sets": 4,
               "reps": 10,
               "weight": 110.0,
               "notes": "Added an extra set",
               "order_index": 1
             }
           ]
         }'
```

#### Delete a Workout
```bash
curl -X DELETE http://localhost:8080/workouts/1 \
     -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

#### Health Status Check
```bash
curl http://localhost:8080/health
```

---

## Running Tests

To execute the unit and database integration tests:
```bash
go test ./...
```

---

## Code Quality & Contributions

This project uses Husky to run pre-commit git hooks. Commits must follow the conventional commit structure (e.g., `feat: ...`, `fix: ...`, `docs: ...`) enforced by `@commitlint`. 

To install the git hooks locally:
```bash
bun install
```