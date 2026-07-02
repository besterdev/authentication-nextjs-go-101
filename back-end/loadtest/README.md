# Load testing

Requires [k6](https://k6.io/docs/get-started/installation/).

## Prerequisites

1. Start infrastructure: `docker compose up -d`
2. Run migrations: `go run ./cmd/migrate`
3. Start API: `go run ./cmd/server`

## Run

```bash
k6 run loadtest/k6.js
```

## Rate limit awareness

The API limits **`/auth/*` to 10 requests/min per IP** (shared bucket). The k6 script respects this:

| Scenario | Default rate | Notes |
|----------|--------------|-------|
| `health` | 100 req/s | Not rate-limited |
| `me` | ~5 req/min | JWT verify path |
| `login` | ~3 req/min | bcrypt-bound |

`setup()` uses 2 auth requests (register + login); scenarios stay within the remaining budget.

## Environment variables

| Variable | Default | Description |
|----------|---------|-------------|
| `BASE_URL` | `http://localhost:8080` | API base URL |
| `TEST_EMAIL` | `loadtest@example.com` | Test user email |
| `TEST_PASSWORD` | `password123` | Test user password |
| `TEST_DURATION` | `30s` | Scenario duration |
| `HEALTH_RATE` | `100` | Health checks per second |
| `AUTH_RATE_LIMIT_PER_MIN` | `10` | Match server rate limit |
| `ME_RATE_PER_MIN` | auto | Override `/auth/me` rate |
| `LOGIN_RATE_PER_MIN` | auto | Override `/auth/login` rate |

Example with custom duration:

```bash
TEST_DURATION=1m k6 run loadtest/k6.js
```

## SLO thresholds (built into script)

- Error rate < 1%
- `/health` p99 < 50ms
- `/auth/me` p99 < 50ms
- `/auth/login` p99 < 500ms

## Metrics endpoint

Prometheus metrics are exposed at `GET /metrics`.
