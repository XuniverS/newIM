package db

import (
	"database/sql"
)

type Message struct {
	ID               int
	SenderID         int
	ReceiverID       int
	EncryptedContent string
	IsRead           bool
	CreatedAt        string
}

func SaveMessage(db *sql.DB, senderID, receiverID int, encryptedContent string) (int, error) {
	var messageID int
	err := db.QueryRow(
		`INSERT INTO messages (sender_id, receiver_id, encrypted_content) 
		 VALUES ($1, $2, $3) RETURNING id`,
		senderID, receiverID, encryptedContent,
	).Scan(&messageID)

	return messageID, err
}

func GetUnreadMessages(db *sql.DB, userID int) ([]Message, error) {
	rows, err := db.Query(
		`SELECT id, sender_id, receiver_id, encrypted_content, is_read, created_at 
		 FROM messages WHERE receiver_id = $1 AND is_read = FALSE 
		 ORDER BY created_at ASC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.EncryptedContent, &msg.IsRead, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}

func MarkMessageAsRead(db *sql.DB, messageID int) error {
	_, err := db.Exec(
		"UPDATE messages SET is_read = TRUE WHERE id = $1",
		messageID,
	)
	return err
}
