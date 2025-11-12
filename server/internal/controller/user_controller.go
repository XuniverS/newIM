package controller

import (
	"net/http"

	"im-system/server/internal/service"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService service.UserService
	wsService   service.WebSocketService
}

// NewUserController 创建用户控制器实例
func NewUserController(userService service.UserService, wsService service.WebSocketService) *UserController {
	return &UserController{
		userService: userService,
		wsService:   wsService,
	}
}

// GetAllUsers 获取所有用户
func (ctrl *UserController) GetAllUsers(c *gin.Context) {
	currentUserID := getUserIDFromContext(c)

	users, err := ctrl.userService.GetAllUsers(currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// GetOnlineUsers 获取在线用户
func (ctrl *UserController) GetOnlineUsers(c *gin.Context) {
	users := ctrl.wsService.GetOnlineUsers()
	c.JSON(http.StatusOK, gin.H{"online_users": users})
}

// 辅助函数：从上下文获取用户ID
func getUserIDFromContext(c *gin.Context) int {
	userID, _ := c.Get("userID")
	return userID.(int)
}

// 辅助函数：从上下文获取用户名
func getUsernameFromContext(c *gin.Context) string {
	username, _ := c.Get("username")
	return username.(string)
}
