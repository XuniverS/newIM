# 🚀 IM 即时通讯系统
一个基于 Go + React 的现代化即时通讯系统，支持端到端加密、离线消息队列、实时 WebSocket 通信。

## 📋 项目特性

### 后端特性
- ✅ **RSA 端到端加密** - 所有消息都通过 RSA 加密传输
- ✅ **WebSocket 实时通信** - 低延迟的双向通信
- ✅ **离线消息队列** - 使用 Kafka 存储离线消息
- ✅ **JWT 身份认证** - 安全的用户认证机制
- ✅ **密钥自动管理** - 用户登录时自动生成和管理 RSA 密钥对
- ✅ **PostgreSQL 数据库** - 持久化存储用户和消息数据

### 前端特性
- ✅ **现代化 UI** - 使用 React 构建的美观界面
- ✅ **实时消息** - 即时接收和发送消息
- ✅ **用户在线状态** - 显示在线用户列表
- ✅ **响应式设计** - 支持各种屏幕尺寸
- ✅ **本地密钥存储** - 私钥安全存储在本地

## 🏗️ 系统架构

```
┌─────────────────────────────────────────────────────────┐
│                    前端 (React)                          │
│  ┌──────────────────────────────────────────────────┐   │
│  │  登录/注册 → 密钥生成 → 聊天界面 → 消息收发      │   │
│  └──────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
                          ↕ WebSocket + HTTP
┌─────────────────────────────────────────────────────────┐
│                  后端 (Go + Gin)                         │
│  ┌──────────────────────────────────────────────────┐   │
│  │  认证 → 密钥管理 → 消息路由 → 在线状态管理      │   │
│  └──────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
         ↕                    ↕                    ↕
    PostgreSQL            Kafka              Redis
    (用户/消息)        (离线消息队列)      (会话缓存)
```

## 📦 项目结构

```
newIM/
├── main.go                          # 主程序入口
├── go.mod                           # Go 模块定义
├── .env.example                     # 环境变量示例
├── internal/
│   ├── config/
│   │   └── config.go               # 配置管理
│   ├── db/
│   │   ├── db.go                   # 数据库初始化
│   │   ├── user.go                 # 用户数据操作
│   │   ├── public_key.go           # 公钥数据操作
│   │   └── message.go              # 消息数据操作
│   ├── crypto/
│   │   └── rsa.go                  # RSA 加密/解密
│   ├── kafka/
│   │   └── kafka.go                # Kafka 消息队列
│   └── server/
│       ├── server.go               # 服务器主体
│       ├── auth.go                 # 认证处理
│       ├── websocket.go            # WebSocket 处理
│       ├── keys.go                 # 密钥管理
│       └── messages.go             # 消息处理
└── web/
    ├── package.json                # 前端依赖
    ├── vite.config.js              # Vite 配置
    ├── index.html                  # HTML 入口
    └── src/
        ├── main.jsx                # React 入口
        ├── App.jsx                 # 主应用组件
        ├── index.css               # 全局样式
        ├── pages/
        │   ├── AuthPage.jsx        # 登录/注册页面
        │   ├── AuthPage.css
        │   ├── ChatPage.jsx        # 聊天页面
        │   └── ChatPage.css
        └── components/
            ├── UserList.jsx        # 用户列表组件
            ├── UserList.css
            ├── ChatWindow.jsx      # 聊天窗口组件
            └── ChatWindow.css
```

## 🚀 快速开始

### 前置要求

- Go 1.21+
- Node.js 16+
- PostgreSQL 12+
- Kafka 3.0+
- Redis 6.0+

### 1. 安装依赖

#### 后端
```bash
cd /Users/xunivers/newIM
go mod tidy
```

#### 前端
```bash
cd web
npm install
```

### 2. 配置环境变量

```bash
# 复制示例配置
cp .env.example .env

# 编辑 .env 文件，配置数据库和其他服务
```

### 3. 启动数据库和消息队列

```bash
# 启动 PostgreSQL
# macOS 使用 Homebrew
brew services start postgresql

# 启动 Kafka
# 假设 Kafka 已安装在 /usr/local/kafka
/usr/local/kafka/bin/kafka-server-start.sh /usr/local/kafka/config/server.properties

# 启动 Redis
redis-server
```

### 4. 启动后端服务

```bash
cd /Users/xunivers/newIM
go run main.go
```

服务器将在 `http://localhost:8080` 启动

### 5. 启动前端开发服务器

```bash
cd web
npm run dev
```

前端将在 `http://localhost:3000` 启动

## 📡 API 文档

### 认证 API

#### 注册
```
POST /api/auth/register
Content-Type: application/json

{
  "username": "user1",
  "password": "password123"
}

Response:
{
  "token": "eyJhbGc...",
  "user_id": 1,
  "username": "user1"
}
```

#### 登录
```
POST /api/auth/login
Content-Type: application/json

{
  "username": "user1",
  "password": "password123"
}

Response:
{
  "token": "eyJhbGc...",
  "user_id": 1,
  "username": "user1"
}
```

### 密钥管理 API

#### 生成密钥对
```
POST /api/keys/generate
Authorization: Bearer {token}

Response:
{
  "public_key": "-----BEGIN PUBLIC KEY-----\n...",
  "private_key": "-----BEGIN RSA PRIVATE KEY-----\n..."
}
```

#### 获取用户公钥
```
GET /api/keys/{userID}
Authorization: Bearer {token}

Response:
{
  "public_key": "-----BEGIN PUBLIC KEY-----\n..."
}
```

#### 上传公钥
```
POST /api/keys/upload
Authorization: Bearer {token}
Content-Type: application/json

{
  "public_key": "-----BEGIN PUBLIC KEY-----\n..."
}

Response:
{
  "message": "Public key uploaded successfully"
}
```

### 消息 API

#### 发送消息
```
POST /api/messages/send
Authorization: Bearer {token}
Content-Type: application/json

{
  "receiver_id": 2,
  "content": "encrypted_message_content"
}

Response:
{
  "message_id": 1,
  "status": "sent"
}
```

#### 获取未读消息
```
GET /api/messages/unread
Authorization: Bearer {token}

Response:
{
  "messages": [
    {
      "id": 1,
      "sender_id": 2,
      "receiver_id": 1,
      "encrypted_content": "...",
      "is_read": false,
      "created_at": "2024-01-01T12:00:00Z"
    }
  ]
}
```

#### 标记消息为已读
```
POST /api/messages/{messageID}/read
Authorization: Bearer {token}

Response:
{
  "message": "Message marked as read"
}
```

### WebSocket API

#### 连接
```
WS /api/ws?token={token}
```

#### 消息格式

发送消息：
```json
{
  "type": "message",
  "receiver_id": 2,
  "content": "encrypted_message_content"
}
```

接收消息：
```json
{
  "type": "message",
  "content": "encrypted_message_content",
  "sender_id": 2,
  "timestamp": "2024-01-01T12:00:00Z"
}
```

## 🔐 安全特性

### 1. 端到端加密
- 使用 RSA-2048 加密算法
- 每个用户拥有唯一的公私钥对
- 消息在客户端加密，服务器无法解密

### 2. 身份认证
- JWT token 认证
- Token 有效期 24 小时
- 支持 Bearer token 认证

### 3. 密钥管理
- 用户注册时自动生成 RSA 密钥对
- 私钥存储在客户端本地
- 公钥上传到服务器供其他用户使用

### 4. 数据库安全
- 密码使用 bcrypt 加密存储
- 支持 SQL 参数化查询防止 SQL 注入

## 📊 数据库模式

### users 表
```sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### public_keys 表
```sql
CREATE TABLE public_keys (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  public_key TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(user_id)
);
```

### messages 表
```sql
CREATE TABLE messages (
  id SERIAL PRIMARY KEY,
  sender_id INTEGER NOT NULL REFERENCES users(id),
  receiver_id INTEGER NOT NULL REFERENCES users(id),
  encrypted_content TEXT NOT NULL,
  is_read BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### sessions 表
```sql
CREATE TABLE sessions (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token VARCHAR(255) UNIQUE NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  expires_at TIMESTAMP NOT NULL
);
```

## 🔄 消息流程

### 在线消息流程
1. 用户 A 获取用户 B 的公钥
2. 用户 A 使用用户 B 的公钥加密消息
3. 用户 A 通过 WebSocket 发送加密消息到服务器
4. 服务器检查用户 B 是否在线
5. 如果在线，服务器直接转发消息给用户 B
6. 用户 B 使用自己的私钥解密消息

### 离线消息流程
1. 用户 A 发送消息时，用户 B 不在线
2. 服务器将消息放入 Kafka 消息队列
3. 用户 B 上线时，从 Kafka 读取离线消息
4. 用户 B 接收并解密消息

## 🛠️ 开发指南

### 添加新的 API 端点

1. 在 `internal/server/` 中创建新的处理函数
2. 在 `server.go` 中注册路由
3. 添加必要的中间件（如认证）

### 修改数据库模式

1. 编辑 `internal/db/db.go` 中的 `createTables` 函数
2. 添加新的数据操作函数
3. 重新运行服务器以应用更改

### 前端开发

1. 在 `web/src/components/` 中创建新组件
2. 在 `web/src/pages/` 中创建新页面
3. 使用 `npm run dev` 启动开发服务器

## 📝 常见问题

### Q: 如何重置数据库？
A: 删除 PostgreSQL 中的数据库，重新启动服务器会自动创建新的表。

### Q: 如何更改 JWT 密钥？
A: 编辑 `.env` 文件中的 `JWT_SECRET` 变量。

### Q: 如何扩展消息加密？
A: 修改 `internal/crypto/rsa.go` 中的加密算法。

### Q: 前端如何存储私钥？
A: 私钥存储在浏览器的 localStorage 中（生产环境建议使用更安全的存储方式）。

## 🚀 生产部署

### 后端部署

1. 构建二进制文件
```bash
go build -o im-server main.go
```

2. 使用 systemd 或 Docker 运行
```bash
./im-server
```

### 前端部署

1. 构建生产版本
```bash
cd web
npm run build
```

2. 将 `dist` 目录部署到 Web 服务器

### Docker 部署

```dockerfile
# 后端 Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o im-server main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/im-server .
EXPOSE 8080
CMD ["./im-server"]
```

## 📄 许可证

MIT License

## 👥 贡献

欢迎提交 Issue 和 Pull Request！

## 📞 联系方式

如有问题或建议，请通过以下方式联系：
- 提交 GitHub Issue
- 发送邮件至 support@example.com

---

**祝你使用愉快！** 🎉
