# 📱 IM 即时通讯系统 - 项目概览
## 🎯 项目目标

构建一个安全、高效、可扩展的即时通讯系统，支持：
- 用户认证与授权
- 端到端加密消息传输
- 实时 WebSocket 通信
- 离线消息队列存储
- 在线用户状态管理

## 🏗️ 系统架构

### 整体架构图

```
┌─────────────────────────────────────────────────────────────┐
│                     客户端 (Web Browser)                     │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  React 应用                                           │  │
│  │  ├─ 登录/注册页面                                    │  │
│  │  ├─ 聊天页面                                         │  │
│  │  ├─ 用户列表                                         │  │
│  │  └─ 消息窗口                                         │  │
│  └───────────────────────────────────────────────────────┘  │
│                          ↕                                    │
│              HTTP + WebSocket (加密)                         │
└─────────────────────────────────────────────────────────────┘
                          ↕
┌─────────────────────────────────────────────────────────────┐
│                   后端服务器 (Go + Gin)                      │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  API 层                                               │  │
│  │  ├─ 认证 API (登录/注册)                             │  │
│  │  ├─ 密钥管理 API                                     │  │
│  │  ├─ 消息 API                                         │  │
│  │  └─ 用户 API                                         │  │
│  └───────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  业务逻辑层                                           │  │
│  │  ├─ WebSocket 连接管理                               │  │
│  │  ├─ 消息路由                                         │  │
│  │  ├─ 在线状态管理                                     │  │
│  │  └─ 离线消息处理                                     │  │
│  └───────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  数据访问层                                           │  │
│  │  ├─ 用户数据操作                                     │  │
│  │  ├─ 消息数据操作                                     │  │
│  │  ├─ 公钥数据操作                                     │  │
│  │  └─ 会话数据操作                                     │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
         ↕                    ↕                    ↕
    PostgreSQL            Kafka              Redis
    (持久化存储)        (消息队列)          (缓存/会话)
```

## 📂 项目结构详解

### 后端结构

```
internal/
├── config/
│   └── config.go              # 配置管理
│       ├── LoadConfig()       # 加载环境变量配置
│       └── getEnv()           # 获取环境变量
│
├── db/
│   ├── db.go                  # 数据库初始化
│   │   ├── InitDB()           # 初始化数据库连接
│   │   └── createTables()     # 创建数据库表
│   ├── user.go                # 用户数据操作
│   │   ├── CreateUser()       # 创建用户
│   │   ├── GetUserByUsername()# 按用户名查询
│   │   ├── GetUserByID()      # 按 ID 查询
│   │   └── VerifyPassword()   # 验证密码
│   ├── public_key.go          # 公钥数据操作
│   │   ├── SavePublicKey()    # 保存公钥
│   │   ├── GetPublicKey()     # 获取公钥
│   │   └── PublicKeyExists()  # 检查公钥是否存在
│   └── message.go             # 消息数据操作
│       ├── SaveMessage()      # 保存消息
│       ├── GetUnreadMessages()# 获取未读消息
│       └── MarkMessageAsRead()# 标记消息为已读
│
├── crypto/
│   └── rsa.go                 # RSA 加密/解密
│       ├── GenerateRSAKeyPair()      # 生成密钥对
│       ├── EncryptWithPublicKey()    # 公钥加密
│       └── DecryptWithPrivateKey()   # 私钥解密
│
├── kafka/
│   └── kafka.go               # Kafka 消息队列
│       ├── InitProducer()     # 初始化生产者
│       ├── PublishMessage()   # 发布消息
│       ├── NewConsumer()      # 创建消费者
│       └── ReadMessage()      # 读取消息
│
└── server/
    ├── server.go              # 服务器主体
    │   ├── NewServer()        # 创建服务器
    │   ├── Run()              # 启动服务器
    │   ├── RegisterClient()   # 注册客户端
    │   ├── UnregisterClient() # 注销客户端
    │   ├── GetClient()        # 获取客户端
    │   └── GetOnlineUsers()   # 获取在线用户
    ├── auth.go                # 认证处理
    │   ├── handleRegister()   # 处理注册
    │   ├── handleLogin()      # 处理登录
    │   ├── generateToken()    # 生成 JWT token
    │   └── authMiddleware()   # 认证中间件
    ├── websocket.go           # WebSocket 处理
    │   ├── handleWebSocket()  # 处理 WebSocket 连接
    │   ├── readPump()         # 读取消息循环
    │   ├── writePump()        # 写入消息循环
    │   └── handleMessageFromClient() # 处理客户端消息
    ├── websocket_auth.go      # WebSocket 认证
    │   └── wsAuthMiddleware() # WebSocket 认证中间件
    ├── keys.go                # 密钥管理
    │   ├── handleGenerateKeys()    # 生成密钥
    │   ├── handleUploadPublicKey() # 上传公钥
    │   └── handleGetPublicKey()    # 获取公钥
    └── messages.go            # 消息处理
        ├── handleSendMessage()     # 发送消息
        ├── handleGetUnreadMessages() # 获取未读消息
        └── handleMarkMessageAsRead() # 标记消息为已读
```

### 前端结构

```
web/src/
├── main.jsx                   # React 入口
├── App.jsx                    # 主应用组件
├── index.css                  # 全局样式
│
├── pages/
│   ├── AuthPage.jsx           # 认证页面
│   │   ├── 登录表单
│   │   ├── 注册表单
│   │   └── 切换登录/注册
│   ├── AuthPage.css
│   ├── ChatPage.jsx           # 聊天页面
│   │   ├── WebSocket 连接管理
│   │   ├── 消息接收处理
│   │   └── 在线用户管理
│   └── ChatPage.css
│
└── components/
    ├── UserList.jsx           # 用户列表组件
    │   ├── 当前用户信息
    │   ├── 在线用户列表
    │   ├── 用户搜索
    │   └── 退出登录
    ├── UserList.css
    ├── ChatWindow.jsx         # 聊天窗口组件
    │   ├── 聊天头部
    │   ├── 消息显示区域
    │   ├── 消息输入框
    │   └── 发送按钮
    └── ChatWindow.css
```

## 🔄 核心流程

### 1. 用户注册流程

```
用户输入用户名和密码
        ↓
前端发送 POST /api/auth/register
        ↓
后端验证用户名是否存在
        ↓
后端使用 bcrypt 加密密码
        ↓
后端保存用户到数据库
        ↓
后端生成 JWT token
        ↓
前端接收 token 和用户信息
        ↓
前端发送 POST /api/keys/generate 生成 RSA 密钥对
        ↓
后端生成 RSA 密钥对
        ↓
后端保存公钥到数据库
        ↓
前端接收公钥和私钥
        ↓
前端保存私钥到 localStorage
        ↓
用户登录成功，进入聊天页面
```

### 2. 消息发送流程

```
用户 A 输入消息
        ↓
前端获取用户 B 的公钥 (GET /api/keys/{userID})
        ↓
前端使用用户 B 的公钥加密消息
        ↓
前端通过 WebSocket 发送加密消息
        ↓
后端接收消息
        ↓
后端检查用户 B 是否在线
        ↓
    ├─ 如果在线：直接转发消息给用户 B
    │
    └─ 如果离线：
        ├─ 保存消息到数据库
        └─ 发布消息到 Kafka 队列
        ↓
用户 B 接收消息
        ↓
用户 B 使用自己的私钥解密消息
        ↓
消息显示在聊天窗口
```

### 3. 离线消息处理流程

```
用户 B 离线时，用户 A 发送消息
        ↓
后端将消息保存到数据库
        ↓
后端发布消息到 Kafka 队列
        ↓
用户 B 上线
        ↓
前端连接 WebSocket
        ↓
后端从 Kafka 读取离线消息
        ↓
后端发送离线消息给用户 B
        ↓
用户 B 接收并解密消息
```

## 🔐 安全机制

### 1. 密码安全
- 使用 bcrypt 算法加密密码
- 密码加盐存储
- 登录时验证密码哈希

### 2. 通信安全
- 使用 HTTPS/WSS 加密传输
- JWT token 认证
- Token 有效期限制

### 3. 消息加密
- RSA-2048 端到端加密
- 每个用户拥有唯一密钥对
- 服务器无法解密消息

### 4. 密钥管理
- 私钥存储在客户端
- 公钥存储在服务器
- 密钥定期轮换机制

## 📊 数据库设计

### users 表
```sql
id              SERIAL PRIMARY KEY
username        VARCHAR(255) UNIQUE NOT NULL
password_hash   VARCHAR(255) NOT NULL
created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
```

### public_keys 表
```sql
id              SERIAL PRIMARY KEY
user_id         INTEGER NOT NULL REFERENCES users(id)
public_key      TEXT NOT NULL
created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
UNIQUE(user_id)
```

### messages 表
```sql
id                  SERIAL PRIMARY KEY
sender_id           INTEGER NOT NULL REFERENCES users(id)
receiver_id         INTEGER NOT NULL REFERENCES users(id)
encrypted_content   TEXT NOT NULL
is_read             BOOLEAN DEFAULT FALSE
created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP
```

### sessions 表
```sql
id          SERIAL PRIMARY KEY
user_id     INTEGER NOT NULL REFERENCES users(id)
token       VARCHAR(255) UNIQUE NOT NULL
created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
expires_at  TIMESTAMP NOT NULL
```

## 🚀 性能优化

### 1. 数据库优化
- 为 username 添加唯一索引
- 为 user_id 添加外键索引
- 为 created_at 添加时间索引

### 2. 缓存优化
- 使用 Redis 缓存在线用户列表
- 缓存用户公钥
- 缓存会话信息

### 3. 消息队列优化
- 使用 Kafka 处理离线消息
- 异步消息处理
- 消息批量处理

### 4. WebSocket 优化
- 连接池管理
- 心跳检测
- 自动重连机制

## 🔧 扩展性

### 1. 水平扩展
- 使用负载均衡器分发请求
- 使用 Redis 共享会话
- 使用 Kafka 分布式消息队列

### 2. 功能扩展
- 群组聊天
- 文件传输
- 语音/视频通话
- 消息搜索
- 消息撤回

### 3. 集成扩展
- 第三方登录 (OAuth)
- 消息推送通知
- 分析统计
- 日志系统

## 📈 监控和日志

### 1. 日志记录
- 请求日志
- 错误日志
- 性能日志
- 安全日志

### 2. 监控指标
- 在线用户数
- 消息吞吐量
- 响应时间
- 错误率

### 3. 告警机制
- 服务异常告警
- 性能告警
- 安全告警

## 🧪 测试策略

### 1. 单元测试
- 加密/解密函数测试
- 数据库操作测试
- 业务逻辑测试

### 2. 集成测试
- API 端点测试
- WebSocket 连接测试
- 消息流程测试

### 3. 性能测试
- 并发连接测试
- 消息吞吐量测试
- 数据库性能测试

### 4. 安全测试
- SQL 注入测试
- XSS 测试
- CSRF 测试
- 加密强度测试

## 📚 技术栈

### 后端
- **语言**: Go 1.21
- **Web 框架**: Gin
- **数据库**: PostgreSQL
- **消息队列**: Kafka
- **缓存**: Redis
- **加密**: RSA, bcrypt
- **认证**: JWT

### 前端
- **框架**: React 18
- **构建工具**: Vite
- **HTTP 客户端**: Axios
- **样式**: CSS3
- **通信**: WebSocket

### 基础设施
- **容器化**: Docker
- **编排**: Docker Compose
- **版本控制**: Git

## 🎓 学习资源

### 后端相关
- [Go 官方文档](https://golang.org/doc/)
- [Gin 框架文档](https://gin-gonic.com/)
- [PostgreSQL 文档](https://www.postgresql.org/docs/)
- [Kafka 文档](https://kafka.apache.org/documentation/)

### 前端相关
- [React 官方文档](https://react.dev/)
- [Vite 文档](https://vitejs.dev/)
- [WebSocket API](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket)

### 安全相关
- [OWASP 安全指南](https://owasp.org/)
- [RSA 加密](https://en.wikipedia.org/wiki/RSA_(cryptosystem))
- [JWT 认证](https://jwt.io/)

---

**最后更新**: 2024 年 1 月
**版本**: 1.0.0
