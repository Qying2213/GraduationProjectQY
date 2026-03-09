#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BACKEND_ENV="$ROOT_DIR/backend/.env"
GOCACHE_VALUE="${GOCACHE:-/tmp/go-build-cache}"
PID_DIR="$ROOT_DIR/.pids"
LAUNCHER="$ROOT_DIR/scripts/service-terminal-launcher.sh"

require_cmd() {
  local cmd="$1"
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "[ERROR] Missing required command: $cmd" >&2
    exit 1
  fi
}

escape_for_applescript() {
  local s="$1"
  s="${s//\\/\\\\}"
  s="${s//\"/\\\"}"
  printf '%s' "$s"
}

open_terminal_for_service() {
  local name="$1"
  local workdir="$2"
  local run_cmd="$3"
  local port="$4"
  local load_backend_env="$5"
  local shell_cmd
  shell_cmd="$(printf '%q ' \
    "$LAUNCHER" \
    "$name" \
    "$workdir" \
    "$run_cmd" \
    "$port" \
    "$load_backend_env" \
    "$PID_DIR" \
    "$BACKEND_ENV" \
    "$GOCACHE_VALUE")"
  shell_cmd="${shell_cmd% }"

  local escaped_cmd
  escaped_cmd="$(escape_for_applescript "$shell_cmd")"

  osascript \
    -e 'tell application "Terminal" to activate' \
    -e "tell application \"Terminal\" to do script \"$escaped_cmd\""
}

require_cmd osascript
require_cmd lsof
require_cmd go
require_cmd npm

if [[ ! -x "$LAUNCHER" ]]; then
  echo "[ERROR] launcher script is not executable: $LAUNCHER" >&2
  exit 1
fi

if [[ ! -f "$BACKEND_ENV" ]]; then
  echo "[WARN] backend/.env not found, backend services will use defaults."
fi

mkdir -p "$PID_DIR"

if [[ ! -d "$ROOT_DIR/frontend/node_modules" ]]; then
  echo "[WARN] frontend/node_modules not found. Run: cd frontend && npm install"
fi

services=(
  "user-service|$ROOT_DIR/backend/user-service|go run main.go|8081|true"
  "job-service|$ROOT_DIR/backend/job-service|go run main.go|8082|true"
  "interview-service|$ROOT_DIR/backend/interview-service|go run main.go|8083|true"
  "resume-service|$ROOT_DIR/backend/resume-service|go run main.go|8084|true"
  "message-service|$ROOT_DIR/backend/message-service|go run main.go|8085|true"
  "talent-service|$ROOT_DIR/backend/talent-service|go run main.go|8086|true"
  "recommendation-service|$ROOT_DIR/backend/recommendation-service|go run main.go|8087|true"
  "log-service|$ROOT_DIR/backend/log-service|go run main.go|8088|true"
  "evaluator-service|$ROOT_DIR/backend/evaluator-service|go run ./cmd/server|8090|true"
  "gateway|$ROOT_DIR/backend/gateway|go run main.go|8080|true"
  "frontend|$ROOT_DIR/frontend|npm run dev -- --host 0.0.0.0 --port 5173|5173|false"
)

echo "[INFO] Opening one Terminal window per service..."

for svc in "${services[@]}"; do
  IFS='|' read -r name workdir run_cmd port load_backend_env <<< "$svc"
  open_terminal_for_service "$name" "$workdir" "$run_cmd" "$port" "$load_backend_env"
  sleep 0.2
done

echo "[INFO] All service terminals have been opened."
