package repository

import (
	"database/sql"

	"im-system/server/internal/model"
)

// MessageRepository 消息数据访问接口
type MessageRepository interface {
	Save(senderID, receiverID int, encryptedContent string) (int, error)
	GetUnread(userID int) ([]model.Message, error)
	MarkAsRead(messageID int) error
	GetConversation(userID1, userID2 int, limit int) ([]model.Message, error)
}

type messageRepository struct {
	db *sql.DB
}

// NewMessageRepository 创建消息仓库实例
func NewMessageRepository(db *sql.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Save(senderID, receiverID int, encryptedContent string) (int, error) {
	var messageID int
	err := r.db.QueryRow(
		`INSERT INTO messages (sender_id, receiver_id, encrypted_content) 
		 VALUES ($1, $2, $3) RETURNING id`,
		senderID, receiverID, encryptedContent,
	).Scan(&messageID)

	return messageID, err
}

func (r *messageRepository) GetUnread(userID int) ([]model.Message, error) {
	rows, err := r.db.Query(
		`SELECT id, sender_id, receiver_id, encrypted_content, is_read, created_at 
		 FROM messages WHERE receiver_id = $1 AND is_read = FALSE 
		 ORDER BY created_at ASC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []model.Message
	for rows.Next() {
		var msg model.Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.EncryptedContent, &msg.IsRead, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}

func (r *messageRepository) MarkAsRead(messageID int) error {
	_, err := r.db.Exec(
		"UPDATE messages SET is_read = TRUE WHERE id = $1",
		messageID,
	)
	return err
}

func (r *messageRepository) GetConversation(userID1, userID2 int, limit int) ([]model.Message, error) {
	rows, err := r.db.Query(
		`SELECT id, sender_id, receiver_id, encrypted_content, is_read, created_at 
		 FROM messages 
		 WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1)
		 ORDER BY created_at DESC
		 LIMIT $3`,
		userID1, userID2, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []model.Message
	for rows.Next() {
		var msg model.Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.EncryptedContent, &msg.IsRead, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}
