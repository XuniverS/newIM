package controller

import (
	"net/http"

	"im-system/client/internal/service"

	"github.com/gin-gonic/gin"
)

// MessageController 消息控制器
type MessageController struct {
	serverService service.ServerService
	wsService     *service.WebSocketService
	cryptoService service.CryptoService
}

// NewMessageController 创建消息控制器实例
func NewMessageController(
	serverService service.ServerService,
	wsService *service.WebSocketService,
	cryptoService service.CryptoService,
) *MessageController {
	return &MessageController{
		serverService: serverService,
		wsService:     wsService,
		cryptoService: cryptoService,
	}
}

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	ReceiverID int    `json:"receiver_id" binding:"required"`
	Content    string `json:"content" binding:"required"`
}

// SendMessage 发送消息
func (ctrl *MessageController) SendMessage(c *gin.Context) {
	token := getTokenFromHeader(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 获取接收者公钥
	publicKey, err := ctrl.serverService.GetPublicKey(token, req.ReceiverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get receiver's public key"})
		return
	}

	// 加密消息
	encryptedContent, err := ctrl.cryptoService.Encrypt(publicKey, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt message"})
		return
	}

	// 发送到服务端
	messageID, err := ctrl.serverService.SendMessage(token, req.ReceiverID, encryptedContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message_id": messageID,
		"status":     "sent",
	})
}

// GetUnreadMessages 获取未读消息
func (ctrl *MessageController) GetUnreadMessages(c *gin.Context) {
	token := getTokenFromHeader(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	privateKey := c.GetHeader("X-Private-Key")

	// 获取未读消息
	messages, err := ctrl.serverService.GetUnreadMessages(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 如果有私钥，解密消息
	if privateKey != "" {
		for i := range messages {
			if messages[i].EncryptedContent != "" {
				decrypted, err := ctrl.cryptoService.Decrypt(privateKey, messages[i].EncryptedContent)
				if err == nil {
					messages[i].Content = decrypted
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}
