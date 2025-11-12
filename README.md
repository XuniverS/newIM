# 🚀 IM 即时通讯系统

一个基于 Go + React 的现代化即时通讯系统，支持端到端 RSA 加密、离线消息队列、实时 WebSocket 通信。

## 📋 项目特性

- ✅ **RSA 端到端加密** - 所有消息都通过 RSA 加密传输
- ✅ **WebSocket 实时通信** - 低延迟的双向通信
- ✅ **离线消息队列** - 使用 Kafka 存储离线消息
- ✅ **JWT 身份认证** - 安全的用户认证机制
- ✅ **密钥自动管理** - 用户登录时自动生成和管理 RSA 密钥对
- ✅ **PostgreSQL 数据库** - 持久化存储用户和消息数据

## 🏗️ 系统架构

```
客户端 (Web Browser)
    ↕ HTTP + WebSocket (加密)
服务端 (Go + Gin)
    ↕
PostgreSQL + Kafka + Redis
```

## 📂 项目结构

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
│       ├── websocket_auth.go       # WebSocket 认证
│       ├── keys.go                 # 密钥管理
│       └── messages.go             # 消息处理
└── web/
    ├── package.json                # 前端依赖
    ├── vite.config.js              # Vite 配置
    ├── index.html                  # HTML 入口
    └── src/
        ├── main.jsx                # React 入口
        ├── App.jsx                 # 主应用组件
        ├── pages/
        │   ├── AuthPage.jsx        # 登录/注册页面
        │   └── ChatPage.jsx        # 聊天页面
        └── components/
            ├── UserList.jsx        # 用户列表组件
            └── ChatWindow.jsx      # 聊天窗口组件
```

## 🔄 核心流程

### 用户注册流程
1. 用户输入用户名和密码
2. 前端发送注册请求到服务端
3. 服务端验证并创建用户（密码使用 bcrypt 加密）
4. 服务端生成 JWT token
5. 前端接收 token 后自动请求生成 RSA 密钥对
6. 服务端生成密钥对，保存公钥到数据库
7. 前端接收私钥并保存到 localStorage

### 消息发送流程（在线）
1. 用户 A 获取用户 B 的公钥
2. 用户 A 使用用户 B 的公钥加密消息
3. 用户 A 通过 WebSocket 发送加密消息
4. 服务端检查用户 B 是否在线
5. 如果在线，直接转发消息给用户 B
6. 用户 B 使用自己的私钥解密消息

### 消息发送流程（离线）
1. 服务端检测到用户 B 离线
2. 将消息保存到数据库
3. 将消息发布到 Kafka 队列
4. 用户 B 上线时，从 Kafka 读取离线消息
5. 用户 B 接收并解密消息

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
```

#### 获取未读消息
```
GET /api/messages/unread
Authorization: Bearer {token}
```

#### 标记消息为已读
```
POST /api/messages/{messageID}/read
Authorization: Bearer {token}
```

### 用户 API

#### 获取在线用户
```
GET /api/users/online
Authorization: Bearer {token}
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

---

## 🚀 初始化和启动指南

### 前置要求

- Go 1.21+
- Node.js 16+
- PostgreSQL 12+
- Kafka 3.0+ (需要先安装 Zookeeper)
- Redis 6.0+ (可选，用于缓存)

### 1. 初始化数据库

#### 安装 PostgreSQL (macOS)
```bash
# 使用 Homebrew 安装
brew install postgresql@15

# 启动 PostgreSQL
brew services start postgresql@15
```

#### 创建数据库
```bash
# 连接到 PostgreSQL
psql postgres

# 创建数据库
CREATE DATABASE im_db;

# 创建用户（可选）
CREATE USER im_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE im_db TO im_user;

# 退出
\q
```

数据库表会在首次运行服务端时自动创建。

### 2. 初始化 Kafka

#### 安装 Kafka (macOS)
```bash
# 使用 Homebrew 安装
brew install kafka

# Kafka 会自动安装 Zookeeper 作为依赖
```

#### 启动 Zookeeper
```bash
# 启动 Zookeeper（Kafka 的依赖）
zookeeper-server-start /usr/local/etc/kafka/zookeeper.properties
```

#### 启动 Kafka（新终端窗口）
```bash
# 启动 Kafka
kafka-server-start /usr/local/etc/kafka/server.properties
```

#### 创建主题（可选，服务端会自动创建）
```bash
# 创建消息主题
kafka-topics --create \
  --bootstrap-server localhost:9092 \
  --topic messages \
  --partitions 3 \
  --replication-factor 1
```

### 3. 配置环境变量

```bash
# 复制示例配置文件
cp .env.example .env

# 编辑 .env 文件
# 根据你的实际配置修改以下内容：
# DB_HOST=localhost
# DB_PORT=5432
# DB_USER=postgres
# DB_PASSWORD=postgres
# DB_NAME=im_db
# KAFKA_HOST=localhost
# KAFKA_PORT=9092
# JWT_SECRET=your-secret-key-change-in-production
```

### 4. 安装依赖

#### 后端依赖
```bash
cd /Users/xunivers/newIM
go mod tidy
```

#### 前端依赖
```bash
cd web
npm install
```

### 5. 启动服务端

```bash
# 在项目根目录
cd /Users/xunivers/newIM
go run main.go
```

成功启动后会看到：
```
🚀 IM Server starting on port 8080
```

### 6. 启动前端开发服务器

```bash
# 在新的终端窗口
cd /Users/xunivers/newIM/web
npm run dev
```

成功启动后会看到：
```
  VITE v5.0.0  ready in 123 ms
  ➜  Local:   http://localhost:3000/
```

### 7. 访问应用

打开浏览器访问：`http://localhost:3000`

### 8. 测试功能

1. **注册账户**：在登录页面点击"注册"，输入用户名和密码
2. **登录**：使用注册的账户登录
3. **查看在线用户**：左侧会显示在线用户列表
4. **发送消息**：选择一个用户，在底部输入框输入消息并发送

### 9. 验证服务状态

#### 检查数据库连接
```bash
psql -h localhost -U postgres -d im_db -c "SELECT * FROM users;"
```

#### 检查 Kafka 主题
```bash
kafka-topics --list --bootstrap-server localhost:9092
```

#### 测试 API
```bash
# 注册用户
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'

# 登录
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

### 常见问题

#### Q: 数据库连接失败
**A:** 检查 PostgreSQL 是否运行：
```bash
brew services list | grep postgresql
```
如果未运行，启动它：
```bash
brew services start postgresql@15
```

#### Q: Kafka 连接失败
**A:** 确保 Zookeeper 和 Kafka 都在运行。检查端口是否被占用：
```bash
lsof -i :2181  # Zookeeper
lsof -i :9092  # Kafka
```

#### Q: 前端无法连接后端
**A:** 检查后端是否在 8080 端口运行：
```bash
lsof -i :8080
```

#### Q: WebSocket 连接失败
**A:** 检查浏览器控制台错误信息，确保 token 正确传递。

### 生产环境部署

#### 构建前端
```bash
cd web
npm run build
```
构建产物在 `web/dist` 目录，可以部署到任何静态文件服务器。

#### 构建后端
```bash
go build -o im-server main.go
```
生成的 `im-server` 可执行文件可以直接运行。

#### 环境变量配置
生产环境务必修改以下配置：
- `JWT_SECRET`: 使用强随机字符串
- `DB_PASSWORD`: 使用强密码
- 配置 HTTPS/WSS 加密传输

---

## 📄 许可证

MIT License

## 👥 贡献

欢迎提交 Issue 和 Pull Request！
