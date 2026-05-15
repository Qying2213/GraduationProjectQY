#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PID_DIR="$ROOT_DIR/tmp/pids"

if [[ ! -d "$PID_DIR" ]]; then
  echo "[INFO] No local service pid directory found."
  exit 0
fi

shopt -s nullglob
pid_files=("$PID_DIR"/*.pid)

if (( ${#pid_files[@]} == 0 )); then
  echo "[INFO] No local service pid files found."
  exit 0
fi

for pid_file in "${pid_files[@]}"; do
  name="$(basename "$pid_file" .pid)"
  pid="$(cat "$pid_file" 2>/dev/null || true)"

  if [[ -n "$pid" ]] && kill -0 "$pid" 2>/dev/null; then
    echo "[STOP] $name (pid $pid)"
    kill "$pid" 2>/dev/null || true
  else
    echo "[SKIP] $name is not running"
  fi

  rm -f "$pid_file"
done

echo "[OK] Local services stopped."
