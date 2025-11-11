.PHONY: help install dev build run clean docker-up docker-down
help:
	@echo "IM 即时通讯系统 - 开发命令"
	@echo ""
	@echo "可用命令:"
	@echo "  make install      - 安装所有依赖"
	@echo "  make dev          - 启动开发环境"
	@echo "  make build        - 构建项目"
	@echo "  make run          - 运行后端服务"
	@echo "  make clean        - 清理构建文件"
	@echo "  make docker-up    - 启动 Docker 容器"
	@echo "  make docker-down  - 停止 Docker 容器"
	@echo "  make frontend-dev - 启动前端开发服务器"
	@echo "  make backend-dev  - 启动后端开发服务器"

install:
	@echo "安装 Go 依赖..."
	go mod tidy
	@echo "安装前端依赖..."
	cd web && npm install
	@echo "✅ 依赖安装完成"

dev: docker-up
	@echo "启动开发环境..."
	@echo "后端服务: http://localhost:8080"
	@echo "前端服务: http://localhost:3000"
	@echo ""
	@echo "在新的终端窗口中运行:"
	@echo "  make backend-dev"
	@echo "  make frontend-dev"

build:
	@echo "构建后端..."
	go build -o im-server main.go
	@echo "构建前端..."
	cd web && npm run build
	@echo "✅ 构建完成"

run: build
	@echo "运行服务器..."
	./im-server

clean:
	@echo "清理构建文件..."
	rm -f im-server
	rm -rf web/dist
	@echo "✅ 清理完成"

docker-up:
	@echo "启动 Docker 容器..."
	docker-compose up -d
	@echo "等待服务启动..."
	sleep 5
	@echo "✅ Docker 容器已启动"
	@echo "  PostgreSQL: localhost:5432"
	@echo "  Kafka: localhost:9092"
	@echo "  Redis: localhost:6379"

docker-down:
	@echo "停止 Docker 容器..."
	docker-compose down
	@echo "✅ Docker 容器已停止"

frontend-dev:
	@echo "启动前端开发服务器..."
	cd web && npm run dev

backend-dev:
	@echo "启动后端开发服务器..."
	go run main.go

logs:
	docker-compose logs -f

ps:
	docker-compose ps
