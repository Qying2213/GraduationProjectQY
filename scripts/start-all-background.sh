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

SKIPPED_SERVICE_NAMES=()
SKIPPED_SERVICE_REASONS=()
FAILED_SERVICE_NAMES=()
FAILED_SERVICE_REASONS=()
SMOKE_CHECK_FAILURES=()
HAS_BLOCKERS="false"

STARTUP_WAIT_SECONDS="${STARTUP_WAIT_SECONDS:-120}"
STARTUP_POLL_SECONDS="${STARTUP_POLL_SECONDS:-0.5}"
STARTUP_LOG_TAIL_LINES="${STARTUP_LOG_TAIL_LINES:-40}"
SMOKE_LOGIN_USER="${SMOKE_LOGIN_USER:-admin}"
SMOKE_LOGIN_PASSWORD="${SMOKE_LOGIN_PASSWORD:-admin123}"

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

record_skipped_service() {
  local service_name="$1"
  local reason="$2"
  SKIPPED_SERVICE_NAMES+=("$service_name")
  SKIPPED_SERVICE_REASONS+=("$reason")
}

record_failed_service() {
  local service_name="$1"
  local reason="$2"
  FAILED_SERVICE_NAMES+=("$service_name")
  FAILED_SERVICE_REASONS+=("$reason")
}

skipped_service_reason() {
  local service_name="$1"
  local i=0

  for ((i = 0; i < ${#SKIPPED_SERVICE_NAMES[@]}; i++)); do
    if [[ "${SKIPPED_SERVICE_NAMES[$i]}" == "$service_name" ]]; then
      printf '%s' "${SKIPPED_SERVICE_REASONS[$i]}"
      return 0
    fi
  done

  return 1
}

failed_service_reason() {
  local service_name="$1"
  local i=0

  for ((i = 0; i < ${#FAILED_SERVICE_NAMES[@]}; i++)); do
    if [[ "${FAILED_SERVICE_NAMES[$i]}" == "$service_name" ]]; then
      printf '%s' "${FAILED_SERVICE_REASONS[$i]}"
      return 0
    fi
  done

  return 1
}

write_skip_log() {
  local service_name="$1"
  local log_file="$2"
  local reason="$3"

  {
    echo "=================================================="
    echo "[$service_name] skipped at $(timestamp)"
    echo "reason: $reason"
    echo "=================================================="
  } > "$log_file"
}

read_pid_file() {
  local pid_file="$1"
  local pid=""

  if [[ -f "$pid_file" ]]; then
    read -r pid < "$pid_file" || true
  fi

  printf '%s' "$pid"
}

pid_is_running() {
  local pid="$1"
  if [[ -z "$pid" ]]; then
    return 1
  fi

  kill -0 "$pid" >/dev/null 2>&1
}

service_runtime_state() {
  local spec="$1"
  local pid_file pid reason

  parse_service_spec "$spec"
  reason="$(failed_service_reason "$SERVICE_NAME" || true)"
  if [[ -n "$reason" ]]; then
    printf 'failed|%s' "$reason"
    return 0
  fi

  reason="$(skipped_service_reason "$SERVICE_NAME" || true)"
  if [[ -n "$reason" ]]; then
    printf 'skipped|%s' "$reason"
    return 0
  fi

  if lsof -t -nP -iTCP:"$SERVICE_PORT" -sTCP:LISTEN >/dev/null 2>&1; then
    printf 'ready|%s' "$SERVICE_URL"
    return 0
  fi

  pid_file="$PID_DIR/$SERVICE_NAME.pid"
  pid="$(read_pid_file "$pid_file")"
  if pid_is_running "$pid"; then
    printf 'starting|%s' "$SERVICE_URL"
    return 0
  fi

  printf 'failed|%s' "$SERVICE_LOG_FILE"
}

wait_for_service_ready() {
  local spec="$1"
  local reason pid_file pid elapsed

  parse_service_spec "$spec"
  reason="$(skipped_service_reason "$SERVICE_NAME" || true)"
  if [[ -n "$reason" ]]; then
    return 0
  fi

  elapsed=0
  pid_file="$PID_DIR/$SERVICE_NAME.pid"

  while awk "BEGIN {exit !($elapsed <= $STARTUP_WAIT_SECONDS)}"; do
    if lsof -t -nP -iTCP:"$SERVICE_PORT" -sTCP:LISTEN >/dev/null 2>&1; then
      info "$SERVICE_NAME is ready on :$SERVICE_PORT"
      return 0
    fi

    pid="$(read_pid_file "$pid_file")"
    if [[ -n "$pid" ]]; then
      if ! pid_is_running "$pid"; then
        local crash_reason="$SERVICE_NAME process (PID: $pid) exited before opening port :$SERVICE_PORT"
        record_failed_service "$SERVICE_NAME" "$crash_reason"
        warn "$crash_reason"
        if [[ -f "$SERVICE_LOG_FILE" ]]; then
          warn "Last $STARTUP_LOG_TAIL_LINES lines of $SERVICE_LOG_FILE:"
          tail -n "$STARTUP_LOG_TAIL_LINES" "$SERVICE_LOG_FILE" || true
        fi
        return 1
      fi
    fi

    sleep "$STARTUP_POLL_SECONDS"
    elapsed="$(awk "BEGIN {print $elapsed + $STARTUP_POLL_SECONDS}")"
  done

  local timeout_reason="$SERVICE_NAME did not become ready on :$SERVICE_PORT within ${STARTUP_WAIT_SECONDS}s"
  record_failed_service "$SERVICE_NAME" "$timeout_reason"
  warn "$timeout_reason"
  if [[ -f "$SERVICE_LOG_FILE" ]]; then
    warn "Last $STARTUP_LOG_TAIL_LINES lines of $SERVICE_LOG_FILE:"
    tail -n "$STARTUP_LOG_TAIL_LINES" "$SERVICE_LOG_FILE" || true
  fi
  return 1
}

run_smoke_checks() {
  local token login_resp code endpoint
  local -a endpoints=(
    "jobs?page=1&page_size=1"
    "talents?page=1&page_size=1"
    "resumes?page=1&page_size=1"
    "evaluations?page=1&page_size=1"
    "ai/current-task"
  )

  if ! lsof -t -nP -iTCP:8080 -sTCP:LISTEN >/dev/null 2>&1; then
    return 0
  fi

  info "Running startup smoke checks through gateway..."
  login_resp="$(curl -sS -m 15 -X POST "http://localhost:8080/api/v1/login" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$SMOKE_LOGIN_USER\",\"password\":\"$SMOKE_LOGIN_PASSWORD\"}" || true)"
  token="$(printf '%s' "$login_resp" | sed -n 's/.*"token":"\([^"]*\)".*/\1/p')"

  if [[ -z "$token" ]]; then
    SMOKE_CHECK_FAILURES+=("gateway auth login failed")
    warn "Gateway auth smoke check failed. Response: ${login_resp:0:180}"
    return 1
  fi

  for endpoint in "${endpoints[@]}"; do
    code="$(curl -sS -m 20 -o /tmp/start-all-smoke.json -w '%{http_code}' \
      -H "Authorization: Bearer $token" \
      "http://localhost:8080/api/v1/$endpoint" || true)"
    if [[ "$code" != "200" ]]; then
      SMOKE_CHECK_FAILURES+=("GET /api/v1/$endpoint returned HTTP $code")
      warn "Gateway API smoke check failed: /api/v1/$endpoint => $code"
      head -c 180 /tmp/start-all-smoke.json 2>/dev/null || true
      echo
      return 1
    fi
  done

  code="$(curl -sS -m 20 -o /tmp/start-all-smoke.json -w '%{http_code}' \
    -H "Authorization: Bearer $token" \
    "http://localhost:8080/api/v1/messages/unread-count" || true)"
  if [[ "$code" != "200" ]]; then
    SMOKE_CHECK_FAILURES+=("GET /api/v1/messages/unread-count returned HTTP $code (possible JWT secret mismatch)")
    warn "JWT smoke check failed: /api/v1/messages/unread-count => $code"
    head -c 180 /tmp/start-all-smoke.json 2>/dev/null || true
    echo
    return 1
  fi

  info "Startup smoke checks passed."
}

print_service_summary() {
  local spec status_data state detail

  HAS_BLOCKERS="false"
  echo "Service status:"

  for spec in "${BACKGROUND_SERVICE_SPECS[@]}"; do
    parse_service_spec "$spec"
    status_data="$(service_runtime_state "$spec")"
    state="${status_data%%|*}"
    detail="${status_data#*|}"

    case "$state" in
      ready)
        printf '  [ready]    %-22s %s\n' "$SERVICE_NAME" "$detail"
        ;;
      starting)
        printf '  [starting] %-22s %s\n' "$SERVICE_NAME" "$detail"
        ;;
      skipped)
        printf '  [skipped]  %-22s %s\n' "$SERVICE_NAME" "$detail"
        HAS_BLOCKERS="true"
        ;;
      failed)
        printf '  [failed]   %-22s see %s\n' "$SERVICE_NAME" "$detail"
        HAS_BLOCKERS="true"
        ;;
    esac
  done
}

start_service() {
  local spec="$1"
  parse_service_spec "$spec"

  local pid_file="$PID_DIR/$SERVICE_NAME.pid"
  local existing_pid prereq_reason

  mkdir -p "$(dirname "$SERVICE_LOG_FILE")"
  existing_pid="$(lsof -t -nP -iTCP:"$SERVICE_PORT" -sTCP:LISTEN 2>/dev/null | head -n 1 || true)"

  if [[ -n "$existing_pid" ]]; then
    echo "$existing_pid" > "$pid_file"
    info "$SERVICE_NAME is already listening on :$SERVICE_PORT (PID: $existing_pid)"
    return 0
  fi

  prereq_reason="$(service_prereq_reason "$SERVICE_NAME" "$SERVICE_WORKDIR" || true)"
  if [[ -n "$prereq_reason" ]]; then
    rm -f "$pid_file"
    write_skip_log "$SERVICE_NAME" "$SERVICE_LOG_FILE" "$prereq_reason"
    record_skipped_service "$SERVICE_NAME" "$prereq_reason"
    warn "Skipping $SERVICE_NAME: $prereq_reason"
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
require_cmd curl
require_cmd awk

if [[ ! -x "$LAUNCHER" ]]; then
  echo "[ERROR] launcher script is not executable: $LAUNCHER" >&2
  exit 1
fi

mkdir -p "$PID_DIR" "$DEV_LOG_DIR"
write_dev_dashboard_manifest

if [[ ! -f "$BACKEND_ENV" ]]; then
  warn "backend/.env not found, backend services will use defaults."
fi

info "Launching all services in background..."

local_spec=""
for local_spec in "${BACKGROUND_SERVICE_SPECS[@]}"; do
  start_service "$local_spec"
  wait_for_service_ready "$local_spec" || true
  sleep 0.2
done

if wait_for_port "$DEV_DASHBOARD_PORT" 20; then
  info "Development dashboard is ready: http://localhost:$DEV_DASHBOARD_PORT"
else
  warn "Development dashboard did not confirm port :$DEV_DASHBOARD_PORT within timeout."
fi

run_smoke_checks || true

echo
echo "=================================================="
echo "Development services launch summary"
print_service_summary
if [[ "${#SMOKE_CHECK_FAILURES[@]}" -gt 0 ]]; then
  HAS_BLOCKERS="true"
  echo "Smoke check failures:"
  for reason in "${SMOKE_CHECK_FAILURES[@]}"; do
    printf '  - %s\n' "$reason"
  done
fi
echo "Logs dir:  $DEV_LOG_DIR"
echo "Stop all:  ./stop-all.sh"
echo "=================================================="

if [[ "$HAS_BLOCKERS" == "true" ]]; then
  exit 1
fi
