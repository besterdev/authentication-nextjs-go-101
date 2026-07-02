# PgBouncer

This project includes PgBouncer in `docker-compose.yml` for production-style connection pooling.

## Local usage

```bash
docker compose up -d pgbouncer
```

Connect through the pooler:

```
DATABASE_URL=postgres://auth_user:auth_pass@localhost:6432/auth_db?sslmode=disable
```

## Settings (via compose environment)

| Variable | Value | Purpose |
|----------|-------|---------|
| `POOL_MODE` | `transaction` | Best for stateless API handlers |
| `MAX_CLIENT_CONN` | `200` | Max client connections to PgBouncer |
| `DEFAULT_POOL_SIZE` | `25` | Server connections per database/user |

Tune `DB_MAX_OPEN_CONNS` in the app to stay below PgBouncer pool limits.
