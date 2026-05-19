#!/usr/bin/env bash

# Ports owned by the local development stack. Database/cache ports are
# intentionally excluded because they are usually managed by Homebrew/Docker.
PROJECT_PORTS=(
  8080 8081 8082 8083 8084 8085 8086 8087 8088 8090 8091
  5173 5174 5175 5176 5177 5178 5179
)

pid_is_alive() {
  local pid="${1:-}"
  [[ -n "$pid" ]] && kill -0 "$pid" 2>/dev/null
}

collect_pid_tree() {
  local pid="$1"
  local child

  pid_is_alive "$pid" || return 0
  echo "$pid"

  while IFS= read -r child; do
    [[ -n "$child" ]] || continue
    collect_pid_tree "$child"
  done < <(pgrep -P "$pid" 2>/dev/null || true)
}

contains_pid() {
  local wanted="$1"
  shift || true

  local pid
  for pid in "$@"; do
    [[ "$pid" == "$wanted" ]] && return 0
  done

  return 1
}

stop_pid_tree() {
  local root_pid="$1"
  local label="${2:-process}"
  local pid
  local pids=()

  pid_is_alive "$root_pid" || return 0

  while IFS= read -r pid; do
    [[ -n "$pid" ]] && pids+=("$pid")
  done < <(collect_pid_tree "$root_pid")

  if (( ${#pids[@]} == 0 )); then
    return 0
  fi

  echo "[STOP] $label (${pids[*]})"

  for pid in "${pids[@]}"; do
    [[ "$pid" == "$$" ]] && continue
    kill "$pid" 2>/dev/null || true
  done

  sleep 1

  for pid in "${pids[@]}"; do
    [[ "$pid" == "$$" ]] && continue
    if pid_is_alive "$pid"; then
      echo "[KILL] $label force stopping pid $pid"
      kill -9 "$pid" 2>/dev/null || true
    fi
  done
}

tracked_pid_tree_from_dir() {
  local pid_dir="$1"
  local pid_file
  local pid

  shopt -s nullglob
  for pid_file in "$pid_dir"/*.pid; do
    pid="$(cat "$pid_file" 2>/dev/null || true)"
    [[ -n "$pid" ]] || continue
    collect_pid_tree "$pid"
  done
}

stop_project_port_listeners() {
  local label="${1:-project ports}"
  shift || true

  local skip_pids=("$@")
  local port
  local pid
  local pids

  for port in "${PROJECT_PORTS[@]}"; do
    pids="$(lsof -tiTCP:"$port" -sTCP:LISTEN 2>/dev/null || true)"
    [[ -n "$pids" ]] || continue

    while IFS= read -r pid; do
      [[ -n "$pid" ]] || continue
      [[ "$pid" == "$$" ]] && continue

      if contains_pid "$pid" "${skip_pids[@]}"; then
        continue
      fi

      echo "[PORT] $label: port $port is held by pid $pid"
      stop_pid_tree "$pid" "port-$port"
    done <<< "$pids"
  done
}
