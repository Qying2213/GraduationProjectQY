#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_ENV="$ROOT_DIR/backend/.env"
SCHEMA_SQL="$ROOT_DIR/backend/databaseSQL/schema.sql"
SEED_SQL="$ROOT_DIR/backend/databaseSQL/init_data.sql"

if [[ -f "$BACKEND_ENV" ]]; then
  set -a
  # shellcheck disable=SC1090
  source "$BACKEND_ENV"
  set +a
fi

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "[ERROR] Missing required command: $1" >&2
    exit 1
  fi
}

require_cmd psql

DB_HOST="${INIT_DB_HOST:-${DB_HOST:-localhost}}"
DB_PORT="${INIT_DB_PORT:-${DB_PORT:-5432}}"
DB_USER="${INIT_DB_USER:-${DB_USER:-postgres}}"
DB_PASSWORD="${INIT_DB_PASSWORD:-${DB_PASSWORD:-}}"
DB_NAME="${INIT_DB_NAME:-${DB_NAME:-talent_platform}}"

if [[ "$DB_HOST" == "postgres" ]]; then
  DB_HOST="localhost"
fi

export PGPASSWORD="$DB_PASSWORD"
PSQL_BASE=(psql -v ON_ERROR_STOP=1 -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER")

echo "[INFO] Initializing database '$DB_NAME' on $DB_HOST:$DB_PORT as $DB_USER"

db_exists="$("${PSQL_BASE[@]}" -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'" | tr -d '[:space:]')"
if [[ "$db_exists" != "1" ]]; then
  echo "[INFO] Database '$DB_NAME' does not exist, creating it..."
  "${PSQL_BASE[@]}" -d postgres -c "CREATE DATABASE \"$DB_NAME\""
fi

echo "[INFO] Applying schema: $SCHEMA_SQL"
"${PSQL_BASE[@]}" -d "$DB_NAME" -f "$SCHEMA_SQL"

echo "[INFO] Seeding test data: $SEED_SQL"
"${PSQL_BASE[@]}" -d "$DB_NAME" -f "$SEED_SQL"

echo "[INFO] Database initialization completed successfully."
