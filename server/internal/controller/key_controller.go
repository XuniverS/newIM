package controller

import (
	"net/http"
	"strconv"

	"im-system/server/internal/service"

	"github.com/gin-gonic/gin"
)

// KeyController 密钥控制器
type KeyController struct {
	keyService service.KeyService
}

// NewKeyController 创建密钥控制器实例
func NewKeyController(keyService service.KeyService) *KeyController {
	return &KeyController{
		keyService: keyService,
	}
}

// GenerateKeysResponse 生成密钥响应
type GenerateKeysResponse struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

// UploadPublicKeyRequest 上传公钥请求
type UploadPublicKeyRequest struct {
	PublicKey string `json:"public_key" binding:"required"`
}

// GenerateKeys 生成密钥对
func (ctrl *KeyController) GenerateKeys(c *gin.Context) {
	userID := getUserIDFromContext(c)

	publicKey, privateKey, err := ctrl.keyService.GenerateKeys(userID)
	if err != nil {
		if err == service.ErrKeyAlreadyExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate keys"})
		}
		return
	}

	c.JSON(http.StatusOK, GenerateKeysResponse{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	})
}

// UploadPublicKey 上传公钥
func (ctrl *KeyController) UploadPublicKey(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req UploadPublicKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := ctrl.keyService.UploadPublicKey(userID, req.PublicKey); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload public key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Public key uploaded successfully"})
}

// GetPublicKey 获取用户公钥
func (ctrl *KeyController) GetPublicKey(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	publicKey, err := ctrl.keyService.GetPublicKey(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Public key not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"public_key": publicKey})
}
