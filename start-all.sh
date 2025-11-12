#!/bin/bash

# IM 系统启动脚本
# 支持启动服务端、客户端或全部服务

# 显示帮助信息
show_help() {
    echo "IM 系统启动脚本"
    echo ""
    echo "用法: ./start-all.sh [选项]"
    echo ""
    echo "选项:"
    echo "  all       启动所有服务（默认）"
    echo "  server    只启动服务端"
    echo "  client    只启动客户端（后端+前端）"
    echo "  help      显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  ./start-all.sh          # 启动所有服务"
    echo "  ./start-all.sh server   # 只启动服务端"
    echo "  ./start-all.sh client   # 只启动客户端"
}

# 检查依赖
check_dependencies() {
    echo "[检查] 检查依赖..."
    
    if ! command -v go &> /dev/null; then
        echo "[错误] Go 未安装"
        echo "请访问 https://golang.org/dl/ 安装 Go"
        exit 1
    fi

    if ! command -v node &> /dev/null; then
        echo "[错误] Node.js 未安装"
        echo "请访问 https://nodejs.org/ 安装 Node.js"
        exit 1
    fi

    if ! command -v psql &> /dev/null; then
        echo "[警告] PostgreSQL 未安装"
        echo "如果需要使用数据库，请安装 PostgreSQL"
    fi

    echo "[完成] 依赖检查通过"
    echo ""
}

# 检查环境变量
check_env() {
    if [ ! -f .env ]; then
        echo "[警告] 未找到 .env 文件，从 .env.example 复制..."
        cp .env.example .env
        echo "[完成] 已创建 .env 文件"
        echo ""
    fi
}

# 启动服务端
start_server() {
    echo "[启动] 启动服务端..."
    cd server
    
    # 安装依赖
    go mod tidy > /dev/null 2>&1
    
    # 启动服务
    go run cmd/server/main.go &
    SERVER_PID=$!
    cd ..
    
    echo "[完成] 服务端已启动 (PID: $SERVER_PID)"
    echo "[地址] 服务端地址: http://localhost:8080"
    echo ""
}

# 启动客户端后端
start_client_backend() {
    echo "[启动] 启动客户端后端..."
    cd client
    
    # 安装依赖
    go mod tidy > /dev/null 2>&1
    
    # 启动服务
    go run cmd/client/main.go &
    CLIENT_BACKEND_PID=$!
    cd ..
    
    echo "[完成] 客户端后端已启动 (PID: $CLIENT_BACKEND_PID)"
    echo "[地址] 客户端后端地址: http://localhost:3001"
    echo ""
}

# 启动客户端前端
start_client_frontend() {
    echo "[启动] 启动客户端前端..."
    cd client/web
    
    # 安装依赖（如果需要）
    if [ ! -d "node_modules" ]; then
        echo "[安装] 安装前端依赖..."
        npm install
    fi
    
    # 启动开发服务器
    npm run dev &
    CLIENT_FRONTEND_PID=$!
    cd ../..
    
    echo "[完成] 客户端前端已启动 (PID: $CLIENT_FRONTEND_PID)"
    echo "[地址] 客户端前端地址: http://localhost:3000"
    echo ""
}

# 主函数
main() {
    MODE=${1:-all}
    
    case $MODE in
        help)
            show_help
            exit 0
            ;;
        server)
            echo "[启动] 启动 IM 服务端..."
            echo ""
            check_dependencies
            check_env
            start_server
            echo "[完成] 服务端启动完成！"
            echo "按 Ctrl+C 停止服务"
            trap "echo ''; echo '[停止] 正在停止服务...'; kill $SERVER_PID 2>/dev/null; exit" INT
            wait
            ;;
        client)
            echo "[启动] 启动 IM 客户端..."
            echo ""
            check_dependencies
            check_env
            start_client_backend
            sleep 2
            start_client_frontend
            echo "[完成] 客户端启动完成！"
            echo "按 Ctrl+C 停止所有服务"
            trap "echo ''; echo '[停止] 正在停止所有服务...'; kill $CLIENT_BACKEND_PID $CLIENT_FRONTEND_PID 2>/dev/null; exit" INT
            wait
            ;;
        all|*)
            echo "[启动] 启动 IM 系统（服务端 + 客户端）..."
            echo ""
            check_dependencies
            check_env
            
            # 启动服务端
            start_server
            sleep 3
            
            # 启动客户端后端
            start_client_backend
            sleep 2
            
            # 启动客户端前端
            start_client_frontend
            
            echo "[完成] IM 系统启动完成！"
            echo ""
            echo "[地址] 访问地址:"
            echo "   服务端: http://localhost:8080"
            echo "   客户端后端: http://localhost:3001"
            echo "   客户端前端: http://localhost:3000"
            echo ""
            echo "按 Ctrl+C 停止所有服务"
            
            # 捕获 Ctrl+C 信号
            trap "echo ''; echo '[停止] 正在停止所有服务...'; kill $SERVER_PID $CLIENT_BACKEND_PID $CLIENT_FRONTEND_PID 2>/dev/null; exit" INT
            
            # 等待
            wait
            ;;
    esac
}

# 运行主函数
main "$@"
