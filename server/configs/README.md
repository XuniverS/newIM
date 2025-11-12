# Configs 配置文件目录
这个目录用于存放服务端的配置文件。

## 用途

- YAML/JSON 配置文件
- 不同环境的配置（开发、测试、生产）
- 数据库迁移脚本
- 初始化数据

## 示例文件

### config.yaml
```yaml
server:
  port: 8080
  host: localhost

database:
  host: localhost
  port: 5432
  name: im_db
  
logging:
  level: info
  format: json
```

### database/migrations/
- `001_create_users_table.sql`
- `002_create_messages_table.sql`
- `003_create_public_keys_table.sql`

## 注意

- 不要提交包含敏感信息的配置文件到Git
- 使用 `.example` 后缀提供配置模板
- 生产环境配置应该通过环境变量或密钥管理服务提供
