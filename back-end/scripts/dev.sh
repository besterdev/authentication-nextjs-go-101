#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

if [[ ! -f .env ]]; then
  echo "Missing .env — run: cp .env.example .env"
  exit 1
fi

export GOPROXY="${GOPROXY:-https://proxy.golang.org,direct}"

echo "Starting PostgreSQL and Redis..."
docker compose up -d

echo "Waiting for PostgreSQL..."
for _ in $(seq 1 30); do
  if docker compose exec -T postgres pg_isready -U auth_user -d auth_db >/dev/null 2>&1; then
    break
  fi
  sleep 1
done

if ! docker compose exec -T postgres pg_isready -U auth_user -d auth_db >/dev/null 2>&1; then
  echo "PostgreSQL did not become ready in time"
  exit 1
fi

echo "Running migrations..."
go run ./cmd/migrate

echo "Starting Go server..."
exec go run ./cmd/server
