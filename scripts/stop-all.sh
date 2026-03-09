#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PID_DIR="$ROOT_DIR/.pids"

SERVICES=(
  "frontend"
  "gateway"
  "user-service"
  "job-service"
  "interview-service"
  "resume-service"
  "message-service"
  "talent-service"
  "recommendation-service"
  "log-service"
  "evaluator-service"
)

timestamp() {
  date "+%Y-%m-%d %H:%M:%S"
}

info() {
  echo "[$(timestamp)] [INFO] $*"
}

warn() {
  echo "[$(timestamp)] [WARN] $*"
}

is_pid_running() {
  local pid="$1"
  kill -0 "$pid" >/dev/null 2>&1
}

get_service_port() {
  local name="$1"
  case "$name" in
    frontend) echo "5173" ;;
    gateway) echo "8080" ;;
    user-service) echo "8081" ;;
    job-service) echo "8082" ;;
    interview-service) echo "8083" ;;
    resume-service) echo "8084" ;;
    message-service) echo "8085" ;;
    talent-service) echo "8086" ;;
    recommendation-service) echo "8087" ;;
    log-service) echo "8088" ;;
    evaluator-service) echo "8090" ;;
    *) echo "" ;;
  esac
}

stop_by_port() {
  local name="$1"
  local port="$2"

  if [[ -z "$port" ]]; then
    return 0
  fi

  local pids
  pids="$(lsof -t -nP -iTCP:"$port" -sTCP:LISTEN 2>/dev/null || true)"
  if [[ -z "$pids" ]]; then
    return 0
  fi

  info "Stopping $name by port $port (PID: $(echo "$pids" | tr '\n' ' ')) ..."
  while IFS= read -r p; do
    [[ -z "$p" ]] && continue
    kill "$p" >/dev/null 2>&1 || true
  done <<< "$pids"

  sleep 1
  pids="$(lsof -t -nP -iTCP:"$port" -sTCP:LISTEN 2>/dev/null || true)"
  if [[ -n "$pids" ]]; then
    warn "$name still listening on port $port, forcing stop."
    while IFS= read -r p; do
      [[ -z "$p" ]] && continue
      kill -9 "$p" >/dev/null 2>&1 || true
    done <<< "$pids"
  fi
}

stop_pid_file() {
  local name="$1"
  local pid_file="$PID_DIR/$name.pid"
  local port
  port="$(get_service_port "$name")"

  if [[ ! -f "$pid_file" ]]; then
    stop_by_port "$name" "$port"
    info "$name pid file not found."
    return 0
  fi

  local pid
  pid="$(cat "$pid_file" 2>/dev/null || true)"
  if [[ -z "$pid" ]]; then
    warn "$name pid file is empty, removing."
    rm -f "$pid_file"
    stop_by_port "$name" "$port"
    return 0
  fi

  if ! is_pid_running "$pid"; then
    warn "$name PID $pid is not running, removing stale pid."
    rm -f "$pid_file"
    stop_by_port "$name" "$port"
    return 0
  fi

  info "Stopping $name (PID: $pid) ..."
  kill "$pid" >/dev/null 2>&1 || true

  for _ in {1..10}; do
    if ! is_pid_running "$pid"; then
      break
    fi
    sleep 1
  done

  if is_pid_running "$pid"; then
    warn "$name did not exit in time, force killing PID $pid."
    kill -9 "$pid" >/dev/null 2>&1 || true
  fi

  rm -f "$pid_file"
  stop_by_port "$name" "$port"
  info "$name stopped."
}

if [[ ! -d "$PID_DIR" ]]; then
  warn "PID directory not found: $PID_DIR"
  exit 0
fi

for svc in "${SERVICES[@]}"; do
  stop_pid_file "$svc"
done

info "All stop tasks finished."
