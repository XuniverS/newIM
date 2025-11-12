package model

import "time"

// User 用户模型
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // 不在JSON中显示
	CreatedAt time.Time `json:"created_at"`
}

// UserPublicInfo 用户公开信息
type UserPublicInfo struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}
