package websocket

import (
	"sync"
)

// Client WebSocket 客户端
type Client struct {
	ID       string
	UserID   string
	Conn     interface{} // *websocket.Conn
	Send     chan interface{}
	IsOnline bool
}

// Hub WebSocket 连接管理中心
type Hub struct {
	clients    map[string]*Client
	broadcast  chan interface{}
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		broadcast:  make(chan interface{}, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// RegisterClient 注册客户端
func (h *Hub) RegisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client.UserID] = client
	client.IsOnline = true
}

// UnregisterClient 注销客户端
func (h *Hub) UnregisterClient(userID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if client, ok := h.clients[userID]; ok {
		client.IsOnline = false
		delete(h.clients, userID)
	}
}

// IsUserOnline 检查用户是否在线
func (h *Hub) IsUserOnline(userID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	client, ok := h.clients[userID]
	return ok && client.IsOnline
}

// GetClient 获取客户端
func (h *Hub) GetClient(userID string) *Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.clients[userID]
}

// BroadcastMessage 广播消息
func (h *Hub) BroadcastMessage(msg interface{}) {
	h.broadcast <- msg
}
