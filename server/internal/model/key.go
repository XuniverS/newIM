package model

import "time"

// PublicKey 公钥模型
type PublicKey struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	PublicKey string    `json:"public_key"`
	CreatedAt time.Time `json:"created_at"`
}
