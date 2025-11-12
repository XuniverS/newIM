# 🚀 IM 即时通讯系统

基于 **Go + React** 的现代化即时通讯系统，采用**客户端-服务端分离**架构，支持端到端 ECC 加密。

## ✨ 核心特性

- ✅ **客户端-服务端分离** - 模拟真实IM应用（如WhatsApp）
- ✅ **端到端加密** - ECC P-256 + ECDH + AES-256-GCM
- ✅ **私钥客户端生成** - 私钥永不通过网络传输
- ✅ **WebSocket 实时通信** - 低延迟双向通信
- ✅ **标准Go项目布局** - cmd、internal、pkg目录结构
- ✅ **分层架构** - Controller-Service-Repository模式
- ✅ **JWT 身份认证** - 安全的用户认证
- ✅ **离线消息支持** - 用户离线时消息保存到数据库

## 🏗️ 系统架构

```
用户浏览器 (React前端)
    ↓ 明文消息
客户端后端 (Go)
    ↓ 加密消息 (ECC+AES)
服务端 (Go)
    ↓ 转发加密消息（不解密）
客户端后端 (Go)
    ↓ 解密消息
用户浏览器 (React前端)
```

**安全特点**：
- 私钥在浏览器中生成（Web Crypto API）
- 私钥永远不会通过网络传输
- 服务端只存储公钥，无法解密消息

## 📂 项目结构

```
newIM/
├── server/              # 服务端（纯后端）
│   ├── cmd/server/      # 服务端入口
│   ├── internal/        # 私有代码
│   │   ├── controller/  # 控制器层
│   │   ├── service/     # 服务层
│   │   ├── repository/  # 数据访问层
│   │   ├── model/       # 数据模型
│   │   ├── middleware/  # 中间件
│   │   └── router/      # 路由配置
│   └── pkg/             # 可复用包
│
├── client/              # 客户端（前端+后端）
│   ├── cmd/client/      # 客户端后端入口
│   ├── internal/        # 私有代码
│   │   ├── controller/  # 控制器层
│   │   ├── service/     # 服务层（加密、通信）
│   │   └── model/       # 数据模型
│   ├── pkg/             # 可复用包
│   └── web/             # React前端
│       └── src/
│           ├── components/  # React组件
│           ├── pages/       # 页面组件
│           └── services/    # API服务
│
├── .env.example         # 环境变量示例
├── start-all.sh         # 启动脚本
└── README.md            # 本文件
```

## 🚀 快速启动

### 前置要求

- Go 1.21+
- Node.js 16+
- PostgreSQL 12+

### 3步启动

```bash
# 1. 创建数据库
psql postgres -c "CREATE DATABASE im_db;"

# 2. 配置环境变量
cp .env.example .env

# 3. 启动所有服务
./start-all.sh
```

访问：http://localhost:3000

### 启动选项

```bash
./start-all.sh          # 启动所有服务（默认）
./start-all.sh server   # 只启动服务端
./start-all.sh client   # 只启动客户端
./start-all.sh help     # 显示帮助
```

## 🔄 核心流程

### 用户注册

```
1. 用户输入用户名和密码
2. 前端发送到客户端后端 → 服务端
3. 服务端创建用户（密码bcrypt加密）
4. 服务端返回JWT token
5. 前端在浏览器中生成ECC密钥对（私钥不传输）
6. 前端保存私钥到localStorage
7. 前端上传公钥到服务端
```

### 消息发送

```
1. 用户A输入消息 "Hello"
2. 客户端后端A获取用户B的公钥
3. 客户端后端A加密消息
4. 发送加密消息到服务端
5. 服务端转发（不解密）
6. 客户端后端B解密消息
7. 用户B看到 "Hello"
```

## 🌐 环境配置

### 本地开发

```bash
# .env
SERVER_HOST=localhost
SERVER_PORT=8080
CLIENT_PORT=3001
```

### 生产部署

```bash
# .env.production
SERVER_HOST=api.yourdomain.com
SERVER_PORT=443
CLIENT_PORT=3001
VITE_API_URL=https://client.yourdomain.com
VITE_WS_URL=wss://client.yourdomain.com
```

## 📡 API 接口

### 服务端 (端口 8080)

- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `GET /api/users` - 获取所有用户
- `GET /api/users/online` - 获取在线用户
- `POST /api/keys/upload` - 上传公钥
- `GET /api/keys/:userID` - 获取用户公钥
- `POST /api/messages/send` - 发送消息
- `GET /api/messages/unread` - 获取未读消息
- `GET /api/ws` - WebSocket连接

### 客户端后端 (端口 3001)

提供类似API，自动处理加密解密

## 📊 数据库

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
  user_id INTEGER NOT NULL REFERENCES users(id),
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

## 🔐 安全特性

1. **端到端加密**
   - 消息在客户端加密，服务端无法解密
   - 使用 ECC P-256 + ECDH + AES-256-GCM

2. **密钥管理**
   - 私钥在浏览器中生成（Web Crypto API）
   - 私钥永远不会通过网络传输
   - 服务端只存储公钥

3. **认证授权**
   - JWT token认证
   - Token有效期24小时

4. **密码安全**
   - bcrypt加密存储

## 🎯 技术栈

### 后端
- **Go 1.21+** - 服务端和客户端后端
- **Gin** - Web框架
- **PostgreSQL** - 数据库
- **JWT** - 认证
- **WebSocket** - 实时通信

### 前端
- **React 18** - UI框架
- **Vite** - 构建工具
- **Axios** - HTTP客户端
- **Web Crypto API** - 密钥生成

## 🐛 常见问题

### 数据库连接失败
```bash
# 检查PostgreSQL是否运行
brew services list | grep postgresql

# 启动PostgreSQL
brew services start postgresql@15
```

### 端口被占用
```bash
# 检查端口
lsof -i :8080  # 服务端
lsof -i :3001  # 客户端后端
lsof -i :3000  # 前端

# 杀死进程
kill -9 <PID>
```

### WebSocket连接失败
- 确保服务端和客户端后端都已启动
- 检查token是否有效
- 查看浏览器控制台错误

## 📝 测试步骤

1. **注册用户Alice**
   - 访问 http://localhost:3000
   - 注册：alice / 123456

2. **注册用户Bob**
   - 新窗口（无痕模式）
   - 注册：bob / 123456

3. **发送消息**
   - Alice窗口：点击Bob，发送"Hello!"
   - Bob窗口：收到消息

4. **验证加密**
   ```bash
   psql im_db -c "SELECT encrypted_content FROM messages;"
   ```
   查看数据库中的消息是加密的！

## 🚀 生产部署

### Docker部署

```bash
# 构建镜像
docker-compose build

# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f
```

### 传统部署

```bash
# 编译服务端
cd server
go build -o server cmd/server/main.go

# 编译客户端后端
cd ../client
go build -o client cmd/client/main.go

# 构建前端
cd web
npm run build

# 使用systemd或supervisor管理服务
```

## 许可证

MIT License

## 学习资源

- [Go项目标准布局](https://github.com/golang-standards/project-layout)
- [Gin Web框架](https://gin-gonic.com/)
- [React官方文档](https://react.dev/)
- [Web Crypto API](https://developer.mozilla.org/en-US/docs/Web/API/Web_Crypto_API)

---

注意：这是一个学习项目，展示了真实IM应用的架构设计。生产环境使用需要进一步完善。
