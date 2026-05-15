#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
LOG_DIR="$ROOT_DIR/logs/dev"
PID_DIR="$ROOT_DIR/tmp/pids"
ENV_FILE="$ROOT_DIR/backend/.env"

mkdir -p "$LOG_DIR" "$PID_DIR"

if [[ -f "$ENV_FILE" ]]; then
  set -a
  # shellcheck disable=SC1090
  source "$ENV_FILE"
  set +a
fi

export DB_HOST="${DB_HOST:-localhost}"
export DB_PORT="${DB_PORT:-5432}"
export DB_USER="${DB_USER:-qinyang}"
export DB_PASSWORD="${DB_PASSWORD:-}"
export DB_NAME="${DB_NAME:-talent_platform}"
export DB_SSLMODE="${DB_SSLMODE:-disable}"
export REDIS_HOST="${REDIS_HOST:-localhost}"
export REDIS_PORT="${REDIS_PORT:-6379}"
export ES_URL="${ES_URL:-http://127.0.0.1:19200}"
export ES_HOST="${ES_HOST:-localhost}"
export ES_PORT="${ES_PORT:-19200}"
export JWT_SECRET="${JWT_SECRET:-qinyang_talent_platform_2024_secret}"

is_running() {
  local pid_file="$1"
  [[ -f "$pid_file" ]] && kill -0 "$(cat "$pid_file")" 2>/dev/null
}

start_cmd() {
  local name="$1"
  local workdir="$2"
  shift 2
  local pid_file="$PID_DIR/$name.pid"
  local log_file="$LOG_DIR/$name.log"

  if is_running "$pid_file"; then
    echo "[SKIP] $name already running (pid $(cat "$pid_file"))"
    return
  fi

  echo "[START] $name"
  (
    cd "$workdir"
    "$@"
  ) >"$log_file" 2>&1 &
  echo "$!" >"$pid_file"
  echo "        pid $(cat "$pid_file"), log $log_file"
}

echo "[INFO] Starting local development services..."
echo "[INFO] This script does not start Docker Compose. Use 'docker compose up -d --build' for the Nginx/Docker deployment."

start_cmd "user-service" "$ROOT_DIR/backend/user-service" go run .
start_cmd "job-service" "$ROOT_DIR/backend/job-service" go run .
start_cmd "talent-service" "$ROOT_DIR/backend/talent-service" go run .
start_cmd "message-service" "$ROOT_DIR/backend/message-service" go run .
start_cmd "interview-service" "$ROOT_DIR/backend/interview-service" go run .
start_cmd "resume-service" "$ROOT_DIR/backend/resume-service" go run .
start_cmd "recommendation-service" "$ROOT_DIR/backend/recommendation-service" go run .
start_cmd "log-service" "$ROOT_DIR/backend/log-service" go run .
start_cmd "evaluator-service" "$ROOT_DIR/backend/evaluator-service" go run ./cmd/server
start_cmd "gateway" "$ROOT_DIR/backend/gateway" go run .
start_cmd "frontend" "$ROOT_DIR/frontend" npm run dev -- --host 0.0.0.0

echo
echo "[OK] Local services are starting in background."
echo "     Frontend: http://localhost:5173"
echo "     Gateway:  http://localhost:8080"
echo "     Logs:     $LOG_DIR"
echo
echo "Stop with: ./stop-all.sh"
