#!/bin/sh
set -eu

mkdir -p /app/logs /app/uploads /app/data

PIDS=""

start_service() {
  name="$1"
  binary="./$name"

  if [ ! -x "$binary" ]; then
    echo "[$name] binary not found or not executable, skip"
    return
  fi

  echo "[$name] starting..."
  "$binary" >"/app/logs/$name.log" 2>&1 &
  pid="$!"
  PIDS="$PIDS $pid"
  echo "[$name] started with pid $pid"
}

stop_all() {
  echo "stopping services..."
  for pid in $PIDS; do
    if kill -0 "$pid" 2>/dev/null; then
      kill "$pid" 2>/dev/null || true
    fi
  done
  wait 2>/dev/null || true
}

trap stop_all INT TERM

start_service "user-service"
start_service "job-service"
start_service "talent-service"
start_service "message-service"
start_service "interview-service"
start_service "resume-service"
start_service "recommendation-service"
start_service "log-service"
start_service "evaluator-service"
start_service "gateway"

echo "all backend services started"

while true; do
  for pid in $PIDS; do
    if ! kill -0 "$pid" 2>/dev/null; then
      echo "service process $pid exited, stopping container"
      stop_all
      exit 1
    fi
  done
  sleep 5
done
