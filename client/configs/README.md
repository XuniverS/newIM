# Configs 配置文件目录
这个目录用于存放客户端后端的配置文件。

## 用途

- 客户端配置文件
- 不同环境的配置
- 服务端连接配置

## 示例

### client.yaml
```yaml
server:
  host: localhost
  port: 8080
  
client:
  port: 3001
  
websocket:
  reconnect_interval: 5s
  max_retries: 10
```

## 注意

- 配置文件不应包含敏感信息
- 使用环境变量覆盖配置
