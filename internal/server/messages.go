package server

import (
	"net/http"
	"strconv"

	"im-system/internal/db"

	"github.com/gin-gonic/gin"
)

type SendMessageRequest struct {
	ReceiverID int    `json:"receiver_id" binding:"required"`
	Content    string `json:"content" binding:"required"`
}

type MessageResponse struct {
	ID               int    `json:"id"`
	SenderID         int    `json:"sender_id"`
	ReceiverID       int    `json:"receiver_id"`
	EncryptedContent string `json:"encrypted_content"`
	IsRead           bool   `json:"is_read"`
	CreatedAt        string `json:"created_at"`
}

func (s *Server) handleSendMessage(c *gin.Context) {
	userID := s.getUserIDFromContext(c)

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 保存消息到数据库
	messageID, err := db.SaveMessage(s.db, userID, req.ReceiverID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
		return
	}

	// 检查接收者是否在线
	receiver := s.GetClient(req.ReceiverID)
	if receiver != nil {
		// 在线，直接发送
		receiver.Send <- WSMessage{
			Type:      "message",
			Content:   req.Content,
			MessageID: messageID,
		}
	} else {
		// 离线，放入 Kafka 队列
		// TODO: 实现 Kafka 消息队列
	}

	c.JSON(http.StatusOK, gin.H{
		"message_id": messageID,
		"status":     "sent",
	})
}

func (s *Server) handleGetUnreadMessages(c *gin.Context) {
	userID := s.getUserIDFromContext(c)

	messages, err := db.GetUnreadMessages(s.db, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	var response []MessageResponse
	for _, msg := range messages {
		response = append(response, MessageResponse{
			ID:               msg.ID,
			SenderID:         msg.SenderID,
			ReceiverID:       msg.ReceiverID,
			EncryptedContent: msg.EncryptedContent,
			IsRead:           msg.IsRead,
			CreatedAt:        msg.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"messages": response})
}

func (s *Server) handleMarkMessageAsRead(c *gin.Context) {
	messageIDStr := c.Param("messageID")
	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	if err := db.MarkMessageAsRead(s.db, messageID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark message as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message marked as read"})
}
