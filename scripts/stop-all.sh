#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PID_DIR="$ROOT_DIR/tmp/pids"

# shellcheck disable=SC1091
source "$ROOT_DIR/scripts/dev-process-utils.sh"

shopt -s nullglob
if [[ -d "$PID_DIR" ]]; then
  pid_files=("$PID_DIR"/*.pid)
else
  echo "[INFO] No local service pid directory found."
  pid_files=()
fi

if (( ${#pid_files[@]} == 0 )); then
  echo "[INFO] No local service pid files found."
else
  for pid_file in "${pid_files[@]}"; do
    name="$(basename "$pid_file" .pid)"
    pid="$(cat "$pid_file" 2>/dev/null || true)"

    if pid_is_alive "$pid"; then
      stop_pid_tree "$pid" "$name"
    else
      echo "[SKIP] $name is not running"
    fi

    rm -f "$pid_file"
  done
fi

if [[ "${SKIP_PORT_CLEANUP:-0}" != "1" ]]; then
  stop_project_port_listeners "cleanup"
fi

echo "[OK] Local services stopped."
