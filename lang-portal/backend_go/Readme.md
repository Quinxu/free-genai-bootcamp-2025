# Language Portal Backend

This is the backend service for the Language Portal application, built with Go, SQLite3, and Gin.

## Prerequisites

- Go 1.20 or later
- Mage (optional; task runner)

Note: The server uses the pure-Go SQLite driver `modernc.org/sqlite` (no CGO required).

## Setup

1. Install Go from https://golang.org/dl/
2. Install Mage (optional):
   ```
   go install github.com/magefile/mage@latest
   ```
3. Download modules:
   ```
   go mod download
   ```

## Running the API

Recommended entrypoint (includes migrations and optional seeding):
```bash
# from backend_go
go run ./cmd/api -port 8090 -db words.db -seed
```
Flags:
- `-port`: server port (default 8080)
- `-db`: SQLite DB path (default words.db)
- `-seed`: when present, seeds initial groups/words from `internal/database/seeds`

## Quick smoke tests
With the server running on port 8090:
```bash
# Health
curl http://localhost:8090/health

# Words list (flattened counts)
curl http://localhost:8090/api/words

# Word show (includes groups)
curl http://localhost:8090/api/words/1

# Groups list
curl http://localhost:8090/api/groups

# Group words (paginated wrapper)
curl http://localhost:8090/api/groups/1/words

# Start a study session
curl -X POST -H "Content-Type: application/json" \
  -d '{"group_id":1,"study_activity_id":1}' \
  http://localhost:8090/api/study_activities

# Record a review (correct/false)
curl -X POST -H "Content-Type: application/json" \
  -d '{"correct":true}' \
  http://localhost:8090/api/study_sessions/1/words/1/review

# Activity sessions (spec shape)
curl http://localhost:8090/api/study_activities/1/study_sessions

# Dashboard
curl http://localhost:8090/api/dashboard/last_study_session
curl http://localhost:8090/api/dashboard/quick_stats
```

## Project Structure

- `cmd/api`: Main application entry point with migrations and seeding
- `internal/api`: API handlers and routes
- `internal/models`: Database models
- `internal/database`: DB connection, migrations and seeds
- `internal/service`: Business logic
- `pkg`: Public packages

## API Documentation

See the Technical Specs document for detailed API endpoint documentation.
