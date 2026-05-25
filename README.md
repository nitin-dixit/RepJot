# RepJot
A simple, fast workout tracker API written in Go.

## Features
- Track workouts and individual exercise entries (sets, reps, weight, duration, notes).
- User registration and authentication using secure token-based access.
- Embedded database migrations with Goose.

## Tech Stack
- Go (Chi Router)
- PostgreSQL (pgxpool)
- Docker Compose

## Getting Started
1. Clone the repository.
2. Set up your `.env` file:
   ```env
   DATABASE_URL=postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable
   ```
3. Start the PostgreSQL database:
   ```bash
   docker-compose up -d
   ```
4. Run the API server:
   ```bash
   go run main.go
   ```