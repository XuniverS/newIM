# WebSocket 客户端包
这个目录用于存放WebSocket客户端相关的代码。

## 用途

- WebSocket 连接管理
- 自动重连逻辑
- 心跳检测
- 消息队列

## 示例

```go
package websocket

import (
    "github.com/gorilla/websocket"
    "time"
)

// Client WebSocket客户端
type Client struct {
    conn *websocket.Conn
    url  string
    reconnectInterval time.Duration
}

// NewClient 创建WebSocket客户端
func NewClient(url string) *Client {
    return &Client{
        url: url,
        reconnectInterval: 5 * time.Second,
    }
}

// Connect 连接到服务器
func (c *Client) Connect() error {
    conn, _, err := websocket.DefaultDialer.Dial(c.url, nil)
    if err != nil {
        return err
    }
    c.conn = conn
    return nil
}

// Send 发送消息
func (c *Client) Send(data []byte) error {
    return c.conn.WriteMessage(websocket.TextMessage, data)
}

// Receive 接收消息
func (c *Client) Receive() ([]byte, error) {
    _, message, err := c.conn.ReadMessage()
    return message, err
}

// Close 关闭连接
func (c *Client) Close() error {
    return c.conn.Close()
}
```

## 功能

- ✅ 自动重连
- ✅ 心跳检测
- ✅ 消息队列
- ✅ 错误处理

## 注意

- 处理网络异常
- 实现优雅关闭
- 避免内存泄漏
