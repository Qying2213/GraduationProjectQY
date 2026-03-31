#!/usr/bin/env bash

ROOT_DIR="${ROOT_DIR:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}"
PID_DIR="${PID_DIR:-$ROOT_DIR/.pids}"
DEV_LOG_DIR="${DEV_LOG_DIR:-$ROOT_DIR/logs/dev}"
BACKEND_ENV="${BACKEND_ENV:-$ROOT_DIR/backend/.env}"
GOCACHE_VALUE="${GOCACHE:-/tmp/go-build-cache}"
DEV_DASHBOARD_PORT="${DEV_DASHBOARD_PORT:-8091}"
DEV_DASHBOARD_MANIFEST="${DEV_DASHBOARD_MANIFEST:-$DEV_LOG_DIR/services.json}"

CORE_SERVICE_SPECS=(
  "user-service|$ROOT_DIR/backend/user-service|go run main.go|8081|true|$DEV_LOG_DIR/user-service.log|http://localhost:8081"
  "job-service|$ROOT_DIR/backend/job-service|go run main.go|8082|true|$DEV_LOG_DIR/job-service.log|http://localhost:8082"
  "interview-service|$ROOT_DIR/backend/interview-service|go run main.go|8083|true|$DEV_LOG_DIR/interview-service.log|http://localhost:8083"
  "resume-service|$ROOT_DIR/backend/resume-service|go run main.go|8084|true|$DEV_LOG_DIR/resume-service.log|http://localhost:8084"
  "message-service|$ROOT_DIR/backend/message-service|go run main.go|8085|true|$DEV_LOG_DIR/message-service.log|http://localhost:8085"
  "talent-service|$ROOT_DIR/backend/talent-service|go run main.go|8086|true|$DEV_LOG_DIR/talent-service.log|http://localhost:8086"
  "recommendation-service|$ROOT_DIR/backend/recommendation-service|go run main.go|8087|true|$DEV_LOG_DIR/recommendation-service.log|http://localhost:8087"
  "log-service|$ROOT_DIR/backend/log-service|go run main.go|8088|true|$DEV_LOG_DIR/log-service.log|http://localhost:8088"
  "evaluator-service|$ROOT_DIR/backend/evaluator-service|go run ./cmd/server|8090|true|$DEV_LOG_DIR/evaluator-service.log|http://localhost:8090"
  "gateway|$ROOT_DIR/backend/gateway|go run main.go|8080|true|$DEV_LOG_DIR/gateway.log|http://localhost:8080/api/v1"
  "frontend|$ROOT_DIR/frontend|npm run dev -- --host 0.0.0.0 --port 5173|5173|false|$DEV_LOG_DIR/frontend.log|http://localhost:5173/login"
)

DEV_DASHBOARD_SPEC="dev-dashboard|$ROOT_DIR|node scripts/dev-dashboard/server.js|$DEV_DASHBOARD_PORT|false|$DEV_LOG_DIR/dev-dashboard.log|http://localhost:$DEV_DASHBOARD_PORT"

BACKGROUND_SERVICE_SPECS=(
  "$DEV_DASHBOARD_SPEC"
  "${CORE_SERVICE_SPECS[@]}"
)

TERMINAL_SERVICE_SPECS=(
  "${CORE_SERVICE_SPECS[@]}"
)

STOPPABLE_SERVICE_SPECS=(
  "${CORE_SERVICE_SPECS[@]}"
  "$DEV_DASHBOARD_SPEC"
)

parse_service_spec() {
  local spec="$1"
  IFS='|' read -r SERVICE_NAME SERVICE_WORKDIR SERVICE_RUN_CMD SERVICE_PORT SERVICE_LOAD_BACKEND_ENV SERVICE_LOG_FILE SERVICE_URL <<< "$spec"
}

service_prereq_reason() {
  local service_name="$1"
  local service_workdir="${2:-}"

  case "$service_name" in
    frontend)
      local frontend_dir="${service_workdir:-$ROOT_DIR/frontend}"
      if [[ ! -d "$frontend_dir/node_modules" ]]; then
        printf '%s' "missing frontend/node_modules. Run: cd frontend && npm install"
        return 0
      fi
      if [[ ! -e "$frontend_dir/node_modules/.bin/vite" ]]; then
        printf '%s' "frontend dependencies are incomplete (vite missing). Run: cd frontend && npm install"
        return 0
      fi
      ;;
  esac

  return 1
}

json_escape() {
  local value="$1"
  value="${value//\\/\\\\}"
  value="${value//\"/\\\"}"
  value="${value//$'\n'/\\n}"
  value="${value//$'\r'/\\r}"
  value="${value//$'\t'/\\t}"
  printf '%s' "$value"
}

write_dev_dashboard_manifest() {
  mkdir -p "$DEV_LOG_DIR" "$PID_DIR"

  {
    printf '[\n'
    local first="true"
    local spec
    for spec in "${BACKGROUND_SERVICE_SPECS[@]}"; do
      parse_service_spec "$spec"
      if [[ "$first" == "true" ]]; then
        first="false"
      else
        printf ',\n'
      fi

      printf '  {"name":"%s","port":%s,"workdir":"%s","logFile":"%s","pidFile":"%s","url":"%s"}' \
        "$(json_escape "$SERVICE_NAME")" \
        "$SERVICE_PORT" \
        "$(json_escape "$SERVICE_WORKDIR")" \
        "$(json_escape "$SERVICE_LOG_FILE")" \
        "$(json_escape "$PID_DIR/$SERVICE_NAME.pid")" \
        "$(json_escape "$SERVICE_URL")"
    done
    printf '\n]\n'
  } > "$DEV_DASHBOARD_MANIFEST"
}
