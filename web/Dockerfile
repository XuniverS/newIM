# 多阶段构建 - 前端
# 第一阶段：构建
FROM node:18-alpine AS builder

WORKDIR /app

# 复制 package.json 和 package-lock.json
COPY package*.json ./

# 安装依赖
RUN npm ci

# 复制源代码
COPY . .

# 构建应用
RUN npm run build

# 第二阶段：运行
FROM node:18-alpine

WORKDIR /app

# 安装 serve 用于提供静态文件
RUN npm install -g serve

# 从构建阶段复制构建结果
COPY --from=builder /app/dist ./dist

# 暴露端口
EXPOSE 3000

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000 || exit 1

# 运行应用
CMD ["serve", "-s", "dist", "-l", "3000"]
