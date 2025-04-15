# Gaming Leaderboard System

A high-performance gaming leaderboard system built with Go, featuring real-time score tracking, player rankings, and a RESTful API.

## Features

- Submit game scores for players
- View top players leaderboard
- Check individual player rankings
- Automatic rank updates
- Database migrations support
- PostgreSQL for data persistence

## Prerequisites

- Go 1.24 or higher
- PostgreSQL 15 or higher
- Docker (optional)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/tashif-hoda/gaming-leaderboard.git
cd gaming-leaderboard
```

2. Set up the database:
```bash
createdb leaderboard
```

3. Set environment variables:
```bash
export DB_HOST=localhost
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=leaderboard
export PORT=8080
export GIN_MODE=debug
export MIGRATIONS_PATH=internal/database/migrations
export API_SECRET_KEY=top-secret-api-key
export NEW_RELIC_LICENSE_KEY=<new-relic-licence-key>
```

4. Run database migrations:
```bash
go run cmd/server/main.go
```

## API Documentation

### Rate Limiting
The API implements rate limiting to ensure fair usage and system stability:

- Score Submission: 5 requests per second with burst of 10
- Leaderboard/Rank Queries: 10 requests per second with burst of 20

Rate limits are applied per IP address. When exceeded, the API returns:
```json
{
    "error": "Rate limit exceeded. Please try again later."
}
```
with HTTP status code 429 (Too Many Requests).

### Submit Score
`POST /api/leaderboard/submit`

Submit a new score for a player.

Request body:
```json
{
    "user_id": 1,
    "score": 1000
}
```

### Get Leaderboard
`GET /api/leaderboard/top?limit=10`

Retrieve top players sorted by total score.

Query Parameters:
- `limit` (optional): Number of players to return (default: 10)

### Get Player Rank
`GET /api/leaderboard/rank/{user_id}`

Get a specific player's current rank.

## Database Migrations

The system uses golang-migrate for database schema management.

Available commands:
```bash
# Run all migrations
go run cmd/server/main.go

# Run down migrations
go run cmd/server/main.go -down

# Migrate to specific version
go run cmd/server/main.go -version 1
```

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── database/
│   │   ├── database.go
│   │   ├── migrate.go
│   │   └── migrations/
│   ├── handlers/
│   │   └── handlers.go
│   └── models/
│       └── models.go
└── README.md
```

## Development

To run the server in development mode:
```bash
go run cmd/server/main.go
```

The server will start on http://localhost:8080 by default.

## Docker Support

### Building the Image
```bash
docker build -t gaming-leaderboard .
```

### Running with Docker
1. First, start a PostgreSQL instance (if you don't have one):
```bash
docker run --name postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=leaderboard -p 5432:5432 -d postgres:12
```

2. Run the application:
```bash
docker run -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_USER=postgres \
  -e DB_PASSWORD=postgres \
  -e DB_NAME=leaderboard \
  -e API_SECRET_KEY=top-secret-api-key \
  -e NEW_RELIC_LICENSE_KEY=<new-relic-licence-key>
  gaming-leaderboard
```

### Docker Compose (Alternative)
A docker-compose.yml file is provided for easy deployment. To run the application with Docker Compose:

```bash
docker-compose up -d
```

This will start both the application and PostgreSQL database in detached mode. The services will be available at:
- Application: http://localhost:8080
- PostgreSQL: localhost:5432

To stop the services:
```bash
docker-compose down
```

To stop the services and remove the persistent volume:
```bash
docker-compose down -v
```

## License

MIT License