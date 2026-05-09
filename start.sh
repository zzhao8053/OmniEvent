#!/bin/bash

# OmniEvent 启动脚本

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# 检查依赖
check_deps() {
    log_info "检查依赖..."

    # 检查 Go
    if ! command -v go &> /dev/null; then
        log_error "Go 未安装"
        exit 1
    fi

    # 检查 Node
    if ! command -v node &> /dev/null; then
        log_error "Node.js 未安装"
        exit 1
    fi

    # 检查前端 node_modules
    if [ ! -d "omnievent-frontend/node_modules" ]; then
        log_warn "前端依赖未安装，正在安装..."
        cd omnievent-frontend
        npm install
        cd ..
    fi

    log_info "依赖检查完成"
}

# 启动后端
start_backend() {
    log_info "启动后端服务..."
    cd omnievent-backend
    go run cmd/server/main.go &
    BACKEND_PID=$!
    echo $BACKEND_PID > /tmp/omnievent_backend.pid
    log_info "后端服务已启动 (PID: $BACKEND_PID)"
}

# 启动前端
start_frontend() {
    log_info "启动前端服务..."
    cd omnievent-frontend
    npm run dev &
    FRONTEND_PID=$!
    echo $FRONTEND_PID > /tmp/omnievent_frontend.pid
    log_info "前端服务已启动 (PID: $FRONTEND_PID)"
}

# 停止服务
stop_services() {
    log_info "停止服务..."

    if [ -f /tmp/omnievent_backend.pid ]; then
        BACKEND_PID=$(cat /tmp/omnievent_backend.pid)
        kill $BACKEND_PID 2>/dev/null || true
        rm /tmp/omnievent_backend.pid
        log_info "后端服务已停止"
    fi

    if [ -f /tmp/omnievent_frontend.pid ]; then
        FRONTEND_PID=$(cat /tmp/omnievent_frontend.pid)
        kill $FRONTEND_PID 2>/dev/null || true
        rm /tmp/omnievent_frontend.pid
        log_info "前端服务已停止"
    fi
}

# 编译检查
build_check() {
    log_info "编译检查..."

    log_info "后端编译..."
    cd omnievent-backend
    go build ./...
    if [ $? -eq 0 ]; then
        log_info "后端编译通过"
    else
        log_error "后端编译失败"
        exit 1
    fi

    log_info "前端类型检查..."
    cd ../omnievent-frontend
    npm run type-check
    if [ $? -eq 0 ]; then
        log_info "前端类型检查通过"
    else
        log_error "前端类型检查失败"
        exit 1
    fi
}

# 使用说明
usage() {
    echo "OmniEvent 启动脚本"
    echo ""
    echo "用法: $0 [命令]"
    echo ""
    echo "命令:"
    echo "  all         启动后端和前端（默认）"
    echo "  backend     仅启动后端"
    echo "  frontend    仅启动前端"
    echo "  stop        停止所有服务"
    echo "  build       编译检查"
    echo "  help        显示帮助"
}

# 主逻辑
case "${1:-all}" in
    all)
        check_deps
        start_backend
        start_frontend
        log_info "OmniEvent 已启动"
        log_info "前端: http://localhost:5173"
        log_info "后端: http://localhost:8080"
        ;;
    backend)
        check_deps
        start_backend
        log_info "后端已启动 http://localhost:8080"
        ;;
    frontend)
        check_deps
        start_frontend
        log_info "前端已启动 http://localhost:5173"
        ;;
    stop)
        stop_services
        ;;
    build)
        build_check
        ;;
    help)
        usage
        ;;
    *)
        log_error "未知命令: $1"
        usage
        exit 1
        ;;
esac
