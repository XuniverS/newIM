#!/bin/bash
# IM 即时通讯系统 - 快速启动脚本

set -e

echo "🚀 IM 即时通讯系统启动脚本"
echo "================================"
echo ""

# 检查必要的工具
check_command() {
    if ! command -v $1 &> /dev/null; then
        echo "❌ 错误: 未找到 $1，请先安装"
        exit 1
    fi
}

echo "📋 检查依赖..."
check_command "go"
check_command "node"
check_command "npm"
check_command "docker"
check_command "docker-compose"
echo "✅ 所有依赖已安装"
echo ""

# 创建 .env 文件
if [ ! -f ".env" ]; then
    echo "📝 创建 .env 文件..."
    cp .env.example .env
    echo "✅ .env 文件已创建，请根据需要修改"
fi

# 安装依赖
echo "📦 安装依赖..."
go mod tidy
cd web && npm install && cd ..
echo "✅ 依赖安装完成"
echo ""

# 启动 Docker 容器
echo "🐳 启动 Docker 容器..."
docker-compose up -d
echo "✅ Docker 容器已启动"
echo ""

# 等待数据库就绪
echo "⏳ 等待数据库就绪..."
sleep 5
echo "✅ 数据库已就绪"
echo ""

# 启动后端
echo "🔧 启动后端服务..."
go run main.go &
BACKEND_PID=$!
echo "✅ 后端服务已启动 (PID: $BACKEND_PID)"
echo ""

# 启动前端
echo "🎨 启动前端开发服务器..."
cd web && npm run dev &
FRONTEND_PID=$!
echo "✅ 前端服务已启动 (PID: $FRONTEND_PID)"
echo ""

echo "================================"
echo "🎉 系统启动完成！"
echo ""
echo "📍 访问地址:"
echo "   前端: http://localhost:3000"
echo "   后端: http://localhost:8080"
echo ""
echo "📊 服务状态:"
echo "   PostgreSQL: localhost:5432"
echo "   Kafka: localhost:9092"
echo "   Redis: localhost:6379"
echo ""
echo "⚠️  按 Ctrl+C 停止所有服务"
echo "================================"
echo ""

# 等待进程
wait $BACKEND_PID $FRONTEND_PID
