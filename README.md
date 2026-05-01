# Chirpy

A small Twitter-like HTTP API written in Go. Users can register, log in, post short messages ("chirps"), and read or filter chirps from other users. Built as a learning project working through the Boot.dev backend course.

## Features

- User registration and login with hashed passwords (argon2id)
- JWT access tokens + database-backed refresh tokens with revocation
- Authenticated chirp creation, deletion (only by the author), and update of own account
- Filtering chirps by author and sorting by date
- "Chirpy Red" premium membership upgraded via incoming Polka webhooks (API-key authenticated)
- Admin endpoints (metrics page, dev-only data reset) gated by `PLATFORM=dev`

## Why it exists

It's a complete backend in idiomatic Go's standard library - no web framework. Useful as a reference for:
- HTTP routing with Go 1.22+ method-prefixed `ServeMux` patterns (e.g. `GET /api/chirps/{chirpID}`)
- JWT issue / verify / refresh / revoke flows
- sqlc + goose for type-safe SQL and versioned migrations
- Webhook authentication via API keys

## Tech stack

- **Go** (standard `net/http`)
- **PostgreSQL** with `lib/pq`
- **goose** for migrations, **sqlc** for query code generation
- **golang-jwt/jwt v5** for JWTs, **argon2id** for password hashing

## Setup

### Prerequisites

- Go 1.22+
- PostgreSQL running locally
- [`goose`](https://github.com/pressly/goose) — `go install github.com/pressly/goose/v3/cmd/goose@latest`
- [`sqlc`](https://docs.sqlc.dev/en/latest/overview/install.html) (only needed if you change SQL files)

### 1. Create the database

```bash
psql postgres -c "CREATE DATABASE chirpy;"
```

### 2. Configure environment

Create a `.env` file in the project root:

```env
DB_URL="postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable"
PLATFORM="dev"
JWT_SECRET="<output of: openssl rand -base64 64>"
POLKA_KEY="<your polka api key>"
```

### 3. Apply migrations

```bash
goose -dir sql/schema postgres "$DB_URL" up
```

### 4. Run

```bash
go run .
```

Server listens on `:8080`.

## Quick try

```bash
# Register
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"you@example.com","password":"password123"}'

# Log in (returns access + refresh tokens)
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"you@example.com","password":"password123"}'

# Post a chirp
curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"body":"hello world"}'

# List all chirps
curl http://localhost:8080/api/chirps

# Filter and sort
curl "http://localhost:8080/api/chirps?author_id=<uuid>&sort=desc"
```
