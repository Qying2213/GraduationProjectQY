#!/usr/bin/env bash
set -euo pipefail

if [[ "$#" -ne 8 ]]; then
  echo "[ERROR] usage: $0 <name> <workdir> <run_cmd> <port> <load_backend_env> <pid_dir> <backend_env> <gocache_value>" >&2
  exit 1
fi

name="$1"
workdir="$2"
run_cmd="$3"
port="$4"
load_backend_env="$5"
pid_dir="$6"
backend_env="$7"
gocache_value="$8"
pid_file="$pid_dir/$name.pid"

cd "$workdir"
if [[ -t 1 ]]; then
  clear || true
fi
echo "=================================================="
echo "[$name] starting at $(date '+%Y-%m-%d %H:%M:%S')"
echo "workdir: $workdir"
echo "port: $port"
echo "=================================================="

mkdir -p "$pid_dir"

existing_pid="$(lsof -t -nP -iTCP:"$port" -sTCP:LISTEN 2>/dev/null | head -n 1 || true)"
if [[ -n "$existing_pid" ]]; then
  echo "[WARN] Port $port is already in use by PID $existing_pid, skip duplicate launch."
  ps -p "$existing_pid" -o pid,ppid,stat,command
  echo "$existing_pid" > "$pid_file"
  echo "[INFO] PID file updated: $pid_file ($existing_pid)"
  exit 0
fi

if [[ "$load_backend_env" == "true" && -f "$backend_env" ]]; then
  set -a
  source "$backend_env"
  set +a
fi
export GOCACHE="$gocache_value"

bash -lc "$run_cmd" &
service_pid=$!
echo "$service_pid" > "$pid_file"
echo "[INFO] Launch PID: $service_pid"

startup_ok=false
startup_attempts="${LAUNCHER_STARTUP_ATTEMPTS:-120}"
startup_interval="${LAUNCHER_STARTUP_INTERVAL_SECONDS:-0.5}"

for _ in $(seq 1 "$startup_attempts"); do
  if ! kill -0 "$service_pid" >/dev/null 2>&1; then
    break
  fi
  listen_pid="$(lsof -t -nP -iTCP:"$port" -sTCP:LISTEN 2>/dev/null | head -n 1 || true)"
  if [[ -n "$listen_pid" ]]; then
    echo "$listen_pid" > "$pid_file"
    echo "[INFO] [$name] listening on :$port (PID: $listen_pid)"
    startup_ok=true
    break
  fi
  sleep "$startup_interval"
done

if [[ "$startup_ok" != "true" ]]; then
  if kill -0 "$service_pid" >/dev/null 2>&1; then
    echo "[WARN] [$name] process is alive but port :$port was not confirmed within timeout."
  else
    echo "[ERROR] [$name] exited before opening port :$port."
    rm -f "$pid_file"
  fi
fi

trap 'rm -f "$pid_file"' EXIT
set +e
wait "$service_pid"
exit_code=$?
set -e
echo
echo "[$name] exited with code $exit_code at $(date '+%Y-%m-%d %H:%M:%S')"
