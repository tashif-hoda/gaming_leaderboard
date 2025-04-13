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
git clone https://github.com/yourusername/gaming-leaderboard.git
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
```

4. Run database migrations:
```bash
go run cmd/server/main.go
```

## API Documentation

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

## License

MIT License