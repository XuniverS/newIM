package model

// WSMessage WebSocket 消息
type WSMessage struct {
	Type       string `json:"type"`
	ReceiverID int    `json:"receiver_id,omitempty"`
	SenderID   int    `json:"sender_id,omitempty"`
	Content    string `json:"content,omitempty"`
	MessageID  int    `json:"message_id,omitempty"`
	Timestamp  string `json:"timestamp,omitempty"`
}

// WSClient WebSocket 客户端
type WSClient struct {
	UserID   int
	Username string
	Send     chan interface{}
}
