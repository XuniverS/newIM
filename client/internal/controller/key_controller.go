package controller

import (
	"net/http"
	"strconv"

	"im-system/client/internal/service"

	"github.com/gin-gonic/gin"
)

// KeyController 密钥控制器
type KeyController struct {
	serverService service.ServerService
	cryptoService service.CryptoService
}

// NewKeyController 创建密钥控制器实例
func NewKeyController(serverService service.ServerService, cryptoService service.CryptoService) *KeyController {
	return &KeyController{
		serverService: serverService,
		cryptoService: cryptoService,
	}
}

// GenerateKeys 生成密钥对
func (ctrl *KeyController) GenerateKeys(c *gin.Context) {
	token := getTokenFromHeader(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	keyPair, err := ctrl.serverService.GenerateKeys(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, keyPair)
}

// GetPublicKey 获取用户公钥
func (ctrl *KeyController) GetPublicKey(c *gin.Context) {
	token := getTokenFromHeader(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	userIDStr := c.Param("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	publicKey, err := ctrl.serverService.GetPublicKey(token, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"public_key": publicKey})
}
