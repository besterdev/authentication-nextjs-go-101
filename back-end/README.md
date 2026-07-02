# Authentication API (Go + Fiber)

Backend authentication API built with JWT Access + Refresh tokens, PostgreSQL, and bcrypt password hashing.

## Architecture

```
Client
  │
  ├─ POST /auth/register  → bcrypt hash → save User
  ├─ POST /auth/login     → verify password → issue JWT pair → store refresh hash
  ├─ POST /auth/refresh   → validate refresh → rotate tokens
  │
  └─ Protected (Bearer access token)
       ├─ GET  /auth/me
       └─ POST /auth/logout → revoke refresh tokens
```

## Prerequisites

- Go 1.22+
- Docker (for PostgreSQL)

## Setup

1. **Start PostgreSQL**

```bash
docker compose up -d
```

2. **Configure environment**

```bash
cp .env.example .env
# Edit .env if needed — secrets must be at least 32 characters
```

3. **Install dependencies**

```bash
GOPROXY=https://proxy.golang.org,direct go mod tidy
```

4. **Run the server**

```bash
go run ./cmd/server
```

Server starts at `http://localhost:8080`.

## API Endpoints

| Method | Path             | Auth | Description              |
| ------ | ---------------- | ---- | ------------------------ |
| GET    | `/health`        | No   | Health check             |
| POST   | `/auth/register` | No   | Register a new user      |
| POST   | `/auth/login`    | No   | Login and get tokens     |
| POST   | `/auth/refresh`  | No   | Rotate access + refresh  |
| POST   | `/auth/logout`   | Yes  | Revoke refresh tokens    |
| GET    | `/auth/me`       | Yes  | Get current user profile |

## Testing with Postman

Import these files from the `postman/` folder:

1. `postman/Auth-API.postman_collection.json` — all API requests with auto-tests
2. `postman/Local.postman_environment.json` — local variables (`baseUrl`, `email`, `password`)

**Suggested run order:** Health Check → Register → Login → Get Me → Refresh Token → Logout

Login and Refresh automatically save `access_token` and `refresh_token` to collection variables.

## Testing with curl

```bash
# Health check
curl http://localhost:8080/health

# Register
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Get profile (replace ACCESS_TOKEN)
curl http://localhost:8080/auth/me \
  -H "Authorization: Bearer ACCESS_TOKEN"

# Refresh tokens (replace REFRESH_TOKEN)
curl -X POST http://localhost:8080/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"REFRESH_TOKEN"}'

# Logout
curl -X POST http://localhost:8080/auth/logout \
  -H "Authorization: Bearer ACCESS_TOKEN"
```

## Key Concepts

### Password hashing (bcrypt)

Passwords are never stored in plain text. `bcrypt.GenerateFromPassword` creates a one-way hash with cost factor 12.

### JWT Access Token

Short-lived (default 15m). Sent in `Authorization: Bearer <token>` header. Contains `user_id` and `email` claims.

### JWT Refresh Token

Long-lived (default 7 days). Only the SHA-256 hash is stored in the database — if the DB leaks, tokens cannot be reused.

### Token rotation

Each `/auth/refresh` call deletes the old refresh token and issues a new pair. This prevents replay attacks.

## Project Structure

```
cmd/server/main.go          Entry point, route wiring
internal/config/            Environment configuration
internal/database/          PostgreSQL connection + migrations
internal/models/            User, RefreshToken GORM models
internal/repository/        Database access layer
internal/service/           Auth business logic
internal/handlers/          HTTP handlers
internal/middleware/        Auth middleware, error handler
internal/dto/               Request/response types
```

## Environment Variables

| Variable             | Description                        | Default        |
| -------------------- | ---------------------------------- | -------------- |
| `PORT`               | HTTP port                          | `8080`         |
| `DATABASE_URL`       | PostgreSQL connection string       | (required)     |
| `JWT_ACCESS_SECRET`  | Access token signing key (≥32 ch)  | (required)     |
| `JWT_REFRESH_SECRET` | Refresh token signing key (≥32 ch) | (required)     |
| `JWT_ACCESS_EXPIRY`  | Access token lifetime              | `15m`          |
| `JWT_REFRESH_EXPIRY` | Refresh token lifetime             | `168h` (7 days)|
| `CORS_ORIGIN`        | Allowed frontend origin            | `http://localhost:3000` |
