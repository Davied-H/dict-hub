#!/bin/bash

# Dict-Hub 一键启动脚本
# 同时启动前端和后端服务

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"
BACKEND_DIR="$PROJECT_ROOT/backend"
FRONTEND_DIR="$PROJECT_ROOT/frontend"

# PID 文件
PID_FILE="$PROJECT_ROOT/.pids"

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_banner() {
    echo -e "${CYAN}"
    echo "╔═══════════════════════════════════════╗"
    echo "║         Dict-Hub 启动脚本             ║"
    echo "╚═══════════════════════════════════════╝"
    echo -e "${NC}"
}

# 检查命令是否存在
check_command() {
    if ! command -v "$1" &> /dev/null; then
        print_error "$1 未安装，请先安装 $1"
        exit 1
    fi
}

# 检查依赖
check_dependencies() {
    print_info "检查依赖..."
    check_command "go"
    check_command "node"
    check_command "npm"
    print_success "依赖检查通过"
}

# 安装前端依赖
install_frontend_deps() {
    if [ ! -d "$FRONTEND_DIR/node_modules" ]; then
        print_info "安装前端依赖..."
        cd "$FRONTEND_DIR"
        npm install
        print_success "前端依赖安装完成"
    fi
}

# 启动后端
start_backend() {
    print_info "启动后端服务..."
    cd "$BACKEND_DIR"
    
    # 确保 data 目录存在
    mkdir -p data
    mkdir -p dicts
    
    go run cmd/server/main.go &
    BACKEND_PID=$!
    echo "BACKEND_PID=$BACKEND_PID" > "$PID_FILE"
    
    print_success "后端服务已启动 (PID: $BACKEND_PID)"
    print_info "后端地址: http://localhost:8080"
}

# 启动前端
start_frontend() {
    print_info "启动前端服务..."
    cd "$FRONTEND_DIR"
    
    npm run dev &
    FRONTEND_PID=$!
    echo "FRONTEND_PID=$FRONTEND_PID" >> "$PID_FILE"
    
    print_success "前端服务已启动 (PID: $FRONTEND_PID)"
    print_info "前端地址: http://localhost:3000"
}

# 停止所有服务
stop_services() {
    print_info "停止所有服务..."
    
    if [ -f "$PID_FILE" ]; then
        source "$PID_FILE"
        
        if [ -n "$BACKEND_PID" ] && kill -0 "$BACKEND_PID" 2>/dev/null; then
            kill "$BACKEND_PID" 2>/dev/null || true
            print_success "后端服务已停止"
        fi
        
        if [ -n "$FRONTEND_PID" ] && kill -0 "$FRONTEND_PID" 2>/dev/null; then
            kill "$FRONTEND_PID" 2>/dev/null || true
            print_success "前端服务已停止"
        fi
        
        rm -f "$PID_FILE"
    fi
    
    # 额外清理可能残留的进程
    pkill -f "go run cmd/server/main.go" 2>/dev/null || true
    pkill -f "vite" 2>/dev/null || true
    
    print_success "所有服务已停止"
}

# 显示状态
show_status() {
    echo ""
    echo -e "${CYAN}═══════════════════════════════════════${NC}"
    echo -e "${GREEN}  服务已全部启动！${NC}"
    echo -e "${CYAN}═══════════════════════════════════════${NC}"
    echo ""
    echo -e "  📦 后端 API:  ${YELLOW}http://localhost:8080${NC}"
    echo -e "  🌐 前端页面:  ${YELLOW}http://localhost:3000${NC}"
    echo ""
    echo -e "${CYAN}═══════════════════════════════════════${NC}"
    echo -e "  按 ${RED}Ctrl+C${NC} 停止所有服务"
    echo -e "${CYAN}═══════════════════════════════════════${NC}"
    echo ""
}

# 处理 Ctrl+C
cleanup() {
    echo ""
    print_warn "收到停止信号..."
    stop_services
    exit 0
}

# 主函数
main() {
    print_banner
    
    case "${1:-start}" in
        start)
            check_dependencies
            install_frontend_deps
            start_backend
            sleep 2  # 等待后端启动
            start_frontend
            show_status
            
            # 捕获 Ctrl+C
            trap cleanup SIGINT SIGTERM
            
            # 等待子进程
            wait
            ;;
        stop)
            stop_services
            ;;
        restart)
            stop_services
            sleep 1
            main start
            ;;
        *)
            echo "用法: $0 {start|stop|restart}"
            exit 1
            ;;
    esac
}

main "$@"
