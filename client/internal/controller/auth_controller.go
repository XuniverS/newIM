package controller

import (
	"net/http"

	"im-system/client/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthController 认证控制器
type AuthController struct {
	serverService service.ServerService
	cryptoService service.CryptoService
}

// NewAuthController 创建认证控制器实例
func NewAuthController(serverService service.ServerService, cryptoService service.CryptoService) *AuthController {
	return &AuthController{
		serverService: serverService,
		cryptoService: cryptoService,
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register 注册
func (ctrl *AuthController) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 调用服务端注册
	authResp, err := ctrl.serverService.Register(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 自动生成密钥对
	keyPair, err := ctrl.serverService.GenerateKeys(authResp.Token)
	if err != nil {
		// 密钥生成失败不影响注册，只记录错误
		c.JSON(http.StatusOK, gin.H{
			"token":    authResp.Token,
			"user_id":  authResp.UserID,
			"username": authResp.Username,
			"error":    "Failed to generate keys: " + err.Error(),
		})
		return
	}

	// 返回token和私钥（私钥需要客户端保存）
	c.JSON(http.StatusOK, gin.H{
		"token":       authResp.Token,
		"user_id":     authResp.UserID,
		"username":    authResp.Username,
		"private_key": keyPair.PrivateKey,
	})
}

// Login 登录
func (ctrl *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 调用服务端登录
	authResp, err := ctrl.serverService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 尝试生成密钥对（如果用户还没有）
	keyPair, err := ctrl.serverService.GenerateKeys(authResp.Token)
	if err != nil {
		// 可能已经存在密钥，这是正常的
		c.JSON(http.StatusOK, gin.H{
			"token":    authResp.Token,
			"user_id":  authResp.UserID,
			"username": authResp.Username,
		})
		return
	}

	// 返回token和私钥
	c.JSON(http.StatusOK, gin.H{
		"token":       authResp.Token,
		"user_id":     authResp.UserID,
		"username":    authResp.Username,
		"private_key": keyPair.PrivateKey,
	})
}
