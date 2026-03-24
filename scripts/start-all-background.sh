#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "$ROOT_DIR/scripts/service-config.sh"

LAUNCHER="$ROOT_DIR/scripts/service-terminal-launcher.sh"

timestamp() {
  date "+%Y-%m-%d %H:%M:%S"
}

info() {
  echo "[$(timestamp)] [INFO] $*"
}

warn() {
  echo "[$(timestamp)] [WARN] $*"
}

require_cmd() {
  local cmd="$1"
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "[ERROR] Missing required command: $cmd" >&2
    exit 1
  fi
}

wait_for_port() {
  local port="$1"
  local attempts="${2:-20}"

  for _ in $(seq 1 "$attempts"); do
    if lsof -t -nP -iTCP:"$port" -sTCP:LISTEN >/dev/null 2>&1; then
      return 0
    fi
    sleep 0.5
  done

  return 1
}

start_service() {
  local spec="$1"
  parse_service_spec "$spec"

  local pid_file="$PID_DIR/$SERVICE_NAME.pid"
  local existing_pid

  mkdir -p "$(dirname "$SERVICE_LOG_FILE")"
  existing_pid="$(lsof -t -nP -iTCP:"$SERVICE_PORT" -sTCP:LISTEN 2>/dev/null | head -n 1 || true)"

  if [[ -n "$existing_pid" ]]; then
    echo "$existing_pid" > "$pid_file"
    info "$SERVICE_NAME is already listening on :$SERVICE_PORT (PID: $existing_pid)"
    return 0
  fi

  : > "$SERVICE_LOG_FILE"
  nohup "$LAUNCHER" \
    "$SERVICE_NAME" \
    "$SERVICE_WORKDIR" \
    "$SERVICE_RUN_CMD" \
    "$SERVICE_PORT" \
    "$SERVICE_LOAD_BACKEND_ENV" \
    "$PID_DIR" \
    "$BACKEND_ENV" \
    "$GOCACHE_VALUE" >>"$SERVICE_LOG_FILE" 2>&1 &

  local launcher_pid="$!"
  info "Starting $SERVICE_NAME on :$SERVICE_PORT (launcher PID: $launcher_pid)"
  info "Log file: $SERVICE_LOG_FILE"
}

require_cmd lsof
require_cmd go
require_cmd npm
require_cmd node

if [[ ! -x "$LAUNCHER" ]]; then
  echo "[ERROR] launcher script is not executable: $LAUNCHER" >&2
  exit 1
fi

mkdir -p "$PID_DIR" "$DEV_LOG_DIR"
write_dev_dashboard_manifest

if [[ ! -f "$BACKEND_ENV" ]]; then
  warn "backend/.env not found, backend services will use defaults."
fi

if [[ ! -d "$ROOT_DIR/frontend/node_modules" ]]; then
  warn "frontend/node_modules not found. Run: cd frontend && npm install"
fi

info "Launching all services in background..."

local_spec=""
for local_spec in "${BACKGROUND_SERVICE_SPECS[@]}"; do
  start_service "$local_spec"
  sleep 0.2
done

if wait_for_port "$DEV_DASHBOARD_PORT" 20; then
  info "Development dashboard is ready: http://localhost:$DEV_DASHBOARD_PORT"
else
  warn "Development dashboard did not confirm port :$DEV_DASHBOARD_PORT within timeout."
fi

echo
echo "=================================================="
echo "Development services are starting in the background."
echo "Dashboard: http://localhost:$DEV_DASHBOARD_PORT"
echo "Frontend:  http://localhost:5173/login"
echo "Gateway:   http://localhost:8080/api/v1"
echo "Logs dir:  $DEV_LOG_DIR"
echo "Stop all:  ./stop-all.sh"
echo "=================================================="
