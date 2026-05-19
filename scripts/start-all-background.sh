#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
LOG_DIR="$ROOT_DIR/logs/dev"
PID_DIR="$ROOT_DIR/tmp/pids"
ENV_FILE="$ROOT_DIR/backend/.env"
START_DELAY_SECONDS="${START_DELAY_SECONDS:-0.5}"
VERIFY_DELAY_SECONDS="${VERIFY_DELAY_SECONDS:-2}"
WAIT_TIMEOUT_SECONDS="${WAIT_TIMEOUT_SECONDS:-90}"

mkdir -p "$LOG_DIR" "$PID_DIR"

# shellcheck disable=SC1091
source "$ROOT_DIR/scripts/dev-process-utils.sh"

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
  [[ -f "$pid_file" ]] && pid_is_alive "$(cat "$pid_file")"
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

  if [[ -f "$pid_file" ]]; then
    echo "[STALE] removing stale pid file for $name"
    rm -f "$pid_file"
  fi

  echo "[START] $name"
  (
    cd "$workdir"
    exec "$@"
  ) >"$log_file" 2>&1 &
  echo "$!" >"$pid_file"
  echo "        pid $(cat "$pid_file"), log $log_file"

  if [[ "$START_DELAY_SECONDS" != "0" ]]; then
    sleep "$START_DELAY_SECONDS"
  fi
}

verify_cmd() {
  local name="$1"
  local pid_file="$PID_DIR/$name.pid"
  local log_file="$LOG_DIR/$name.log"

  if is_running "$pid_file"; then
    echo "[READY] $name (pid $(cat "$pid_file"))"
    return 0
  fi

  echo "[FAIL] $name exited early, check log: $log_file"
  return 1
}

wait_http() {
  local name="$1"
  local url="$2"
  local deadline=$((SECONDS + WAIT_TIMEOUT_SECONDS))

  printf "[WAIT] %s %s" "$name" "$url"
  while (( SECONDS < deadline )); do
    if curl -fsS --max-time 2 "$url" >/dev/null 2>&1; then
      echo " ready"
      return 0
    fi

    printf "."
    sleep 2
  done

  echo " timeout"
  return 1
}

echo "[INFO] Starting local development services..."
echo "[INFO] This script does not start Docker Compose. Use 'docker compose up -d --build' for the Nginx/Docker deployment."

if [[ "${SKIP_PORT_CLEANUP:-0}" != "1" ]]; then
  tracked_pids=()
  while IFS= read -r pid; do
    [[ -n "$pid" ]] && tracked_pids+=("$pid")
  done < <(tracked_pid_tree_from_dir "$PID_DIR" | sort -u)

  if (( ${#tracked_pids[@]} > 0 )); then
    stop_project_port_listeners "startup cleanup" "${tracked_pids[@]}"
  else
    stop_project_port_listeners "startup cleanup"
  fi
fi

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
start_cmd "frontend" "$ROOT_DIR/frontend" npm run dev -- --host 0.0.0.0 --strictPort

echo
echo "[INFO] Verifying started processes..."
sleep "$VERIFY_DELAY_SECONDS"

failed=0
verify_cmd "user-service" || failed=1
verify_cmd "job-service" || failed=1
verify_cmd "talent-service" || failed=1
verify_cmd "message-service" || failed=1
verify_cmd "interview-service" || failed=1
verify_cmd "resume-service" || failed=1
verify_cmd "recommendation-service" || failed=1
verify_cmd "log-service" || failed=1
verify_cmd "evaluator-service" || failed=1
verify_cmd "gateway" || failed=1
verify_cmd "frontend" || failed=1

wait_http "gateway" "http://localhost:8080/health" || failed=1
wait_http "frontend" "http://localhost:5173/" || failed=1

if (( failed != 0 )); then
  echo
  echo "[ERROR] Some services failed to stay running. See logs in $LOG_DIR"
  exit 1
fi

echo
echo "[OK] Local services are starting in background."
echo "     Admin:    http://localhost:5173/login"
echo "     Portal:   http://localhost:5173/portal"
echo "     Frontend: http://localhost:5173"
echo "     Gateway:  http://localhost:8080/health"
echo "     Swagger:  http://localhost:8080/swagger/index.html"
echo "     Logs:     $LOG_DIR"
echo
echo "Stop with: ./stop-all.sh"
