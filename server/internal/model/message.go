package model

import "time"

// Message 消息模型
type Message struct {
	ID               int       `json:"id"`
	SenderID         int       `json:"sender_id"`
	ReceiverID       int       `json:"receiver_id"`
	EncryptedContent string    `json:"encrypted_content"`
	IsRead           bool      `json:"is_read"`
	CreatedAt        time.Time `json:"created_at"`
}

// MessageDTO 消息传输对象
type MessageDTO struct {
	ID               int    `json:"id"`
	SenderID         int    `json:"sender_id"`
	ReceiverID       int    `json:"receiver_id"`
	EncryptedContent string `json:"encrypted_content"`
	IsRead           bool   `json:"is_read"`
	CreatedAt        string `json:"created_at"`
}
