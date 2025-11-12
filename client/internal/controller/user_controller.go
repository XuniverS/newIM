package controller

import (
	"net/http"

	"im-system/client/internal/service"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	serverService service.ServerService
}

// NewUserController 创建用户控制器实例
func NewUserController(serverService service.ServerService) *UserController {
	return &UserController{
		serverService: serverService,
	}
}

// GetAllUsers 获取所有用户
func (ctrl *UserController) GetAllUsers(c *gin.Context) {
	token := getTokenFromHeader(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	users, err := ctrl.serverService.GetAllUsers(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// GetOnlineUsers 获取在线用户
func (ctrl *UserController) GetOnlineUsers(c *gin.Context) {
	token := getTokenFromHeader(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	users, err := ctrl.serverService.GetOnlineUsers(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"online_users": users})
}

// 辅助函数：从header获取token
func getTokenFromHeader(c *gin.Context) string {
	token := c.GetHeader("Authorization")
	if len(token) > 7 && token[:7] == "Bearer " {
		return token[7:]
	}
	return token
}
