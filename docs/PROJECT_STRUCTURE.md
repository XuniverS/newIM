# IM 项目结构说明
## 项目概述
这是一个即时通讯（IM）系统，支持实时消息传输、身份认证、消息加密和离线消息队列。

## 目录结构

```
newIM/
├── cmd/                          # 可执行程序入口
│   ├── server/                   # 服务器程序
│   │   └── main.go              # 服务器主程序
│   └── client/                   # 客户端程序（CLI）
│       └── main.go              # 客户端主程序
│
├── backend/                       # 后端核心逻辑
│   ├── api/                      # API 处理层
│   │   └── handler.go           # HTTP/WebSocket 处理器
│   │
│   ├── auth/                     # 身份认证模块
│   │   └── auth.go              # 用户认证、令牌管理
│   │
│   ├── message/                  # 消息处理模块
│   │   └── message.go           # 消息业务逻辑
│   │
│   ├── websocket/                # WebSocket 连接管理
│   │   └── websocket.go         # 客户端连接、在线状态管理
│   │
│   ├── queue/                    # 消息队列模块
│   │   └── queue.go             # 离线消息队列、消费者
│   │
│   ├── database/                 # 数据库模块
│   │   └── database.go          # 数据库连接、仓库接口
│   │
│   ├── encryption/               # 加密模块
│   │   └── encryption.go        # 消息加密、解密
│   │
│   └── config/                   # 配置模块
│       └── config.go            # 应用配置管理
│
├── frontend/                      # 前端代码
│   ├── src/                      # 源代码
│   └── public/                   # 静态资源
│
├── docs/                          # 文档
│   └── PROJECT_STRUCTURE.md      # 项目结构说明
│
├── go.mod                         # Go 模块文件
├── main.go                        # 项目根入口（可选）
└── .gitignore                     # Git 忽略文件
```

## 核心模块说明

### 1. 身份认证 (auth)
- 用户注册和登录
- JWT 令牌生成和验证
- 令牌刷新机制

### 2. WebSocket 连接管理 (websocket)
- 客户端连接管理
- 在线状态跟踪
- 消息广播

### 3. 消息处理 (message)
- 消息存储和检索
- 消息状态管理（pending, sent, delivered, read）
- 消息业务逻辑

### 4. 消息队列 (queue)
- 离线消息暂存
- 消息消费
- 队列管理

### 5. 数据库 (database)
- 用户数据存储
- 消息历史存储
- 仓库模式实现

### 6. 加密 (encryption)
- 消息加密/解密
- 支持多种加密算法
- 初始化向量（IV）管理

### 7. API 处理 (api)
- HTTP 路由处理
- WebSocket 升级处理
- 请求验证和响应格式化

## 工作流程

1. **用户认证**
   - 用户通过前端注册/登录
   - 服务器验证身份并返回 JWT 令牌
   - 客户端使用令牌建立 WebSocket 连接

2. **消息发送**
   - 客户端加密消息
   - 发送到服务器
   - 服务器验证并转发

3. **在线消息转发**
   - 如果接收者在线，直接通过 WebSocket 转发
   - 消息加密传输

4. **离线消息处理**
   - 如果接收者离线，消息存入队列
   - 消息同时存储到数据库
   - 接收者上线后，服务器从队列消费消息并发送

## 技术栈

- **后端**: Go
- **前端**: Web（HTML/CSS/JavaScript）
- **数据库**: 待定（MySQL/PostgreSQL）
- **消息队列**: 待定（Redis/RabbitMQ）
- **实时通信**: WebSocket
- **加密**: AES/RSA（待定）
- **认证**: JWT

## 下一步

1. 实现具体的数据库连接（MySQL/PostgreSQL）
2. 选择并集成消息队列（Redis/RabbitMQ）
3. 实现加密算法（AES-256-GCM）
4. 开发前端 Web 应用
5. 编写单元测试和集成测试
6. 部署和性能优化
