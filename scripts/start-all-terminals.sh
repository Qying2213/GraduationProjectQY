#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "$ROOT_DIR/scripts/service-config.sh"
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

echo "[INFO] Opening one Terminal window per service..."
has_blockers="false"

for svc in "${TERMINAL_SERVICE_SPECS[@]}"; do
  IFS='|' read -r name workdir run_cmd port load_backend_env _log_file _url <<< "$svc"
  prereq_reason="$(service_prereq_reason "$name" "$workdir" || true)"
  if [[ -n "$prereq_reason" ]]; then
    echo "[WARN] Skipping $name: $prereq_reason"
    has_blockers="true"
    continue
  fi
  open_terminal_for_service "$name" "$workdir" "$run_cmd" "$port" "$load_backend_env"
  sleep 0.2
done

echo "[INFO] All service terminals have been opened."

if [[ "$has_blockers" == "true" ]]; then
  exit 1
fi
