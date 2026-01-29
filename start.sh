#!/bin/bash

# ============================================
# 智能人才运营平台 - 一键启动脚本
# ============================================
# 使用方法: ./start.sh [选项]
#   ./start.sh          - 启动所有服务（后端+前端）
#   ./start.sh backend  - 只启动后端服务
#   ./start.sh frontend - 只启动前端
#   ./start.sh stop     - 停止所有服务
#   ./start.sh status   - 查看服务状态
# ============================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"
ENV_FILE="$PROJECT_ROOT/backend/.env"
LOG_DIR="$PROJECT_ROOT/logs"
PID_DIR="$PROJECT_ROOT/.pids"

# 服务定义: 目录名:端口:服务名
SERVICES=(
    "gateway:8080:Gateway网关"
    "user-service:8081:用户服务"
    "job-service:8082:职位服务"
    "interview-service:8083:面试服务"
    "resume-service:8084:简历服务"
    "message-service:8085:消息服务"
    "talent-service:8086:人才服务"
    "recommendation-service:8087:推荐服务"
)

print_banner() {
    echo ""
    echo -e "${CYAN}╔════════════════════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║${NC}       ${GREEN}智能人才运营平台 - 启动管理脚本${NC}              ${CYAN}║${NC}"
    echo -e "${CYAN}╚════════════════════════════════════════════════════════╝${NC}"
    echo ""
}

print_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
print_success() { echo -e "${GREEN}[✓]${NC} $1"; }
print_error() { echo -e "${RED}[✗]${NC} $1"; }
print_warn() { echo -e "${YELLOW}[!]${NC} $1"; }

# 创建必要目录
init_dirs() {
    mkdir -p "$LOG_DIR"
    mkdir -p "$PID_DIR"
}

# 加载环境变量
load_env() {
    if [ -f "$ENV_FILE" ]; then
        print_info "加载环境变量: $ENV_FILE"
        set -a
        source "$ENV_FILE"
        set +a
    else
        print_warn "未找到 .env 文件，使用默认配置"
        print_info "可复制 backend/.env.example 为 backend/.env"
    fi
}

# 检查依赖
check_dependencies() {
    print_info "检查依赖..."
    
    # 检查 Go
    if ! command -v go &> /dev/null; then
        print_error "未安装 Go，请先安装 Go 1.21+"
        exit 1
    fi
    print_success "Go $(go version | awk '{print $3}')"
    
    # 检查 Node.js
    if ! command -v node &> /dev/null; then
        print_warn "未安装 Node.js，前端将无法启动"
    else
        print_success "Node.js $(node -v)"
    fi
    
    # 检查 PostgreSQL
    if command -v psql &> /dev/null; then
        print_success "PostgreSQL 已安装"
    else
        print_warn "未检测到 psql 命令，请确保 PostgreSQL 服务已运行"
    fi
    
    echo ""
}

# 检查端口是否被占用
check_port() {
    local port=$1
    if lsof -i :$port &> /dev/null; then
        return 0  # 端口被占用
    else
        return 1  # 端口空闲
    fi
}

# 启动单个后端服务
start_service() {
    local dir=$1
    local port=$2
    local name=$3
    local service_path="$PROJECT_ROOT/backend/$dir"
    local log_file="$LOG_DIR/$dir.log"
    local pid_file="$PID_DIR/$dir.pid"
    
    # 检查目录是否存在
    if [ ! -d "$service_path" ]; then
        print_error "$name 目录不存在: $service_path"
        return 1
    fi
    
    # 检查端口
    if check_port $port; then
        print_warn "$name (端口 $port) 已在运行"
        return 0
    fi
    
    print_info "启动 $name (端口: $port)..."
    
    # 启动服务
    cd "$service_path"
    nohup go run main.go > "$log_file" 2>&1 &
    local pid=$!
    echo $pid > "$pid_file"
    
    # 等待服务启动
    sleep 2
    
    if check_port $port; then
        print_success "$name 启动成功 (PID: $pid)"
    else
        print_error "$name 启动失败，查看日志: $log_file"
        return 1
    fi
}

# 启动所有后端服务
start_backend() {
    print_info "启动后端服务..."
    echo ""
    
    load_env
    
    for service in "${SERVICES[@]}"; do
        IFS=':' read -r dir port name <<< "$service"
        start_service "$dir" "$port" "$name"
    done
    
    echo ""
    print_success "后端服务启动完成！"
}

# 启动前端
start_frontend() {
    print_info "启动前端服务..."
    
    local frontend_path="$PROJECT_ROOT/frontend"
    local log_file="$LOG_DIR/frontend.log"
    local pid_file="$PID_DIR/frontend.pid"
    
    if [ ! -d "$frontend_path" ]; then
        print_error "前端目录不存在"
        return 1
    fi
    
    # 检查端口 5173
    if check_port 5173; then
        print_warn "前端 (端口 5173) 已在运行"
        return 0
    fi
    
    cd "$frontend_path"
    
    # 检查 node_modules
    if [ ! -d "node_modules" ]; then
        print_info "安装前端依赖..."
        npm install
    fi
    
    # 启动前端
    nohup npm run dev > "$log_file" 2>&1 &
    local pid=$!
    echo $pid > "$pid_file"
    
    sleep 3
    
    if check_port 5173; then
        print_success "前端启动成功 (PID: $pid)"
    else
        print_error "前端启动失败，查看日志: $log_file"
    fi
}

# 停止所有服务
stop_all() {
    print_info "停止所有服务..."
    
    # 停止后端服务
    for service in "${SERVICES[@]}"; do
        IFS=':' read -r dir port name <<< "$service"
        local pid_file="$PID_DIR/$dir.pid"
        
        if [ -f "$pid_file" ]; then
            local pid=$(cat "$pid_file")
            if kill -0 $pid 2>/dev/null; then
                kill $pid 2>/dev/null
                print_success "停止 $name (PID: $pid)"
            fi
            rm -f "$pid_file"
        fi
        
        # 强制杀死占用端口的进程
        local pids=$(lsof -t -i :$port 2>/dev/null)
        if [ -n "$pids" ]; then
            echo "$pids" | xargs kill -9 2>/dev/null
            print_success "强制停止 $name (端口: $port)"
        fi
    done
    
    # 停止前端
    print_info "停止前端服务..."
    local frontend_pid_file="$PID_DIR/frontend.pid"
    if [ -f "$frontend_pid_file" ]; then
        local pid=$(cat "$frontend_pid_file")
        if kill -0 $pid 2>/dev/null; then
            kill $pid 2>/dev/null
            print_success "停止前端 (PID: $pid)"
        fi
        rm -f "$frontend_pid_file"
    fi
    
    # 强制杀死前端端口（包括 node 进程）
    local frontend_pids=$(lsof -t -i :5173 2>/dev/null)
    if [ -n "$frontend_pids" ]; then
        echo "$frontend_pids" | xargs kill -9 2>/dev/null
        print_success "强制停止前端 (端口: 5173)"
    fi
    
    # 杀死所有 vite 相关进程
    pkill -f "vite" 2>/dev/null && print_success "停止 Vite 进程" || true
    
    # 等待进程完全退出
    sleep 1
    
    print_success "所有服务已停止"
}

# 查看服务状态
show_status() {
    echo ""
    echo -e "${CYAN}服务状态:${NC}"
    echo "─────────────────────────────────────────────"
    printf "%-20s %-8s %-10s\n" "服务名称" "端口" "状态"
    echo "─────────────────────────────────────────────"
    
    for service in "${SERVICES[@]}"; do
        IFS=':' read -r dir port name <<< "$service"
        if check_port $port; then
            printf "%-20s %-8s ${GREEN}%-10s${NC}\n" "$name" "$port" "运行中"
        else
            printf "%-20s %-8s ${RED}%-10s${NC}\n" "$name" "$port" "未运行"
        fi
    done
    
    # 前端状态
    if check_port 5173; then
        printf "%-20s %-8s ${GREEN}%-10s${NC}\n" "前端" "5173" "运行中"
    else
        printf "%-20s %-8s ${RED}%-10s${NC}\n" "前端" "5173" "未运行"
    fi
    
    echo "─────────────────────────────────────────────"
    echo ""
}

# 显示访问地址
show_urls() {
    echo ""
    echo -e "${CYAN}访问地址:${NC}"
    echo "─────────────────────────────────────────────"
    echo -e "  前端页面:     ${GREEN}http://localhost:5173${NC}"
    echo -e "  求职者端:     ${GREEN}http://localhost:5173/portal${NC}"
    echo -e "  后台管理:     ${GREEN}http://localhost:5173/login${NC}"
    echo -e "  API网关:      ${GREEN}http://localhost:8080${NC}"
    echo -e "  健康检查:     ${GREEN}http://localhost:8080/health${NC}"
    echo "─────────────────────────────────────────────"
    echo ""
    echo -e "${YELLOW}日志目录: $LOG_DIR${NC}"
    echo ""
}

# 显示帮助
show_help() {
    echo "使用方法: ./start.sh [命令]"
    echo ""
    echo "命令:"
    echo "  (无参数)    启动所有服务（后端+前端）"
    echo "  backend     只启动后端服务"
    echo "  frontend    只启动前端"
    echo "  stop        停止所有服务"
    echo "  restart     重启所有服务"
    echo "  status      查看服务状态"
    echo "  logs        查看日志（tail -f）"
    echo "  help        显示帮助"
    echo ""
}

# 查看日志
show_logs() {
    local service=${1:-gateway}
    local log_file="$LOG_DIR/$service.log"
    
    if [ -f "$log_file" ]; then
        tail -f "$log_file"
    else
        print_error "日志文件不存在: $log_file"
        echo "可用的日志: $(ls $LOG_DIR 2>/dev/null | tr '\n' ' ')"
    fi
}

# 主函数
main() {
    print_banner
    init_dirs
    
    case "${1:-all}" in
        all|start)
            check_dependencies
            start_backend
            echo ""
            start_frontend
            show_urls
            ;;
        backend)
            check_dependencies
            start_backend
            show_urls
            ;;
        frontend)
            start_frontend
            show_urls
            ;;
        stop)
            stop_all
            ;;
        restart)
            stop_all
            sleep 2
            check_dependencies
            start_backend
            start_frontend
            show_urls
            ;;
        status)
            show_status
            ;;
        logs)
            show_logs "${2:-gateway}"
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            print_error "未知命令: $1"
            show_help
            exit 1
            ;;
    esac
}

main "$@"
