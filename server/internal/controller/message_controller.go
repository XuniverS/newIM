package controller

import (
	"net/http"
	"strconv"

	"im-system/server/internal/service"

	"github.com/gin-gonic/gin"
)

// MessageController 消息控制器
type MessageController struct {
	messageService service.MessageService
}

// NewMessageController 创建消息控制器实例
func NewMessageController(messageService service.MessageService) *MessageController {
	return &MessageController{
		messageService: messageService,
	}
}

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	ReceiverID int    `json:"receiver_id" binding:"required"`
	Content    string `json:"content" binding:"required"`
}

// SendMessage 发送消息
func (ctrl *MessageController) SendMessage(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	messageID, err := ctrl.messageService.SendMessage(userID, req.ReceiverID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message_id": messageID,
		"status":     "sent",
	})
}

// GetUnreadMessages 获取未读消息
func (ctrl *MessageController) GetUnreadMessages(c *gin.Context) {
	userID := getUserIDFromContext(c)

	messages, err := ctrl.messageService.GetUnreadMessages(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

// MarkMessageAsRead 标记消息为已读
func (ctrl *MessageController) MarkMessageAsRead(c *gin.Context) {
	messageIDStr := c.Param("messageID")
	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	if err := ctrl.messageService.MarkAsRead(messageID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark message as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message marked as read"})
}
