#!/usr/bin/env bash

# ============================================================
# 智能人才运营平台 - 一键运行脚本
# ============================================================
# 用法:
#   ./run-project.sh            # 启动前后端（默认）
#   ./run-project.sh start      # 启动前后端
#   ./run-project.sh stop       # 停止并移除容器
#   ./run-project.sh restart    # 重启
#   ./run-project.sh status     # 查看运行状态
#   ./run-project.sh logs       # 查看日志
# ============================================================

set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

info() { echo -e "${BLUE}[INFO]${NC} $*"; }
ok() { echo -e "${GREEN}[OK]${NC} $*"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $*"; }
err() { echo -e "${RED}[ERROR]${NC} $*"; }

compose_cmd() {
  if docker compose version >/dev/null 2>&1; then
    docker compose "$@"
    return
  fi

  if command -v docker-compose >/dev/null 2>&1; then
    docker-compose "$@"
    return
  fi

  err "未检测到 docker compose（或 docker-compose）。请先安装 Docker Desktop / OrbStack。"
  exit 1
}

wait_for_http() {
  local name="$1"
  local url="$2"
  local timeout="${3:-120}"
  local i=0

  info "等待 $name 就绪: $url"
  while [ "$i" -lt "$timeout" ]; do
    if curl -fsS -m 2 "$url" >/dev/null 2>&1; then
      ok "$name 已就绪"
      return 0
    fi
    sleep 1
    i=$((i + 1))
  done

  err "$name 启动超时（>${timeout}s）"
  return 1
}

start() {
  info "启动项目（前端 + 后端 + 依赖）..."
  cd "$PROJECT_ROOT"
  compose_cmd up -d --build

  wait_for_http "网关" "http://localhost:8080/health" 180
  wait_for_http "前端" "http://localhost:3000" 120

  echo ""
  ok "项目已可用"
  echo "前端地址: http://localhost:3000"
  echo "网关地址: http://localhost:8080"
  echo "日志页面: http://localhost:8088/logs"
}

stop() {
  info "停止项目..."
  cd "$PROJECT_ROOT"
  compose_cmd down
  ok "项目已停止"
}

status() {
  cd "$PROJECT_ROOT"
  compose_cmd ps
  echo ""

  if curl -fsS -m 2 "http://localhost:8080/health" >/dev/null 2>&1; then
    ok "网关健康检查: OK"
  else
    warn "网关健康检查: FAIL"
  fi

  if curl -fsS -m 2 "http://localhost:3000" >/dev/null 2>&1; then
    ok "前端访问检查: OK"
  else
    warn "前端访问检查: FAIL"
  fi
}

logs() {
  cd "$PROJECT_ROOT"
  compose_cmd logs -f --tail=120
}

restart() {
  stop
  start
}

usage() {
  cat <<'EOF'
用法:
  ./run-project.sh [start|stop|restart|status|logs]
EOF
}

cmd="${1:-start}"
case "$cmd" in
  start) start ;;
  stop) stop ;;
  restart) restart ;;
  status) status ;;
  logs) logs ;;
  -h|--help|help) usage ;;
  *)
    err "未知命令: $cmd"
    usage
    exit 1
    ;;
esac

