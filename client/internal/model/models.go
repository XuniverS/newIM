package model

// AuthRequest 认证请求
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	Token    string `json:"token"`
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
}

// User 用户信息
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// Message 消息
type Message struct {
	ID               int    `json:"id"`
	SenderID         int    `json:"sender_id"`
	ReceiverID       int    `json:"receiver_id"`
	EncryptedContent string `json:"encrypted_content"`
	Content          string `json:"content"` // 解密后的内容
	IsRead           bool   `json:"is_read"`
	CreatedAt        string `json:"created_at"`
}

// KeyPair 密钥对
type KeyPair struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

// WSMessage WebSocket 消息
type WSMessage struct {
	Type       string `json:"type"`
	ReceiverID int    `json:"receiver_id,omitempty"`
	SenderID   int    `json:"sender_id,omitempty"`
	Content    string `json:"content,omitempty"`
	MessageID  int    `json:"message_id,omitempty"`
	Timestamp  string `json:"timestamp,omitempty"`
}
