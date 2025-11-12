package router

import (
	"im-system/server/internal/config"
	"im-system/server/internal/controller"
	"im-system/server/internal/middleware"
	"im-system/server/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter(
	cfg *config.Config,
	userService service.UserService,
	messageService service.MessageService,
	keyService service.KeyService,
	wsService service.WebSocketService,
) *gin.Engine {
	router := gin.Default()

	// 添加 CORS 中间件
	router.Use(middleware.CORSMiddleware())

	// 初始化控制器
	authCtrl := controller.NewAuthController(userService)
	userCtrl := controller.NewUserController(userService, wsService)
	messageCtrl := controller.NewMessageController(messageService)
	keyCtrl := controller.NewKeyController(keyService)
	wsCtrl := controller.NewWebSocketController(wsService)

	// API 路由组
	api := router.Group("/api")
	{
		// 认证路由（无需认证）
		auth := api.Group("/auth")
		{
			auth.POST("/register", authCtrl.Register)
			auth.POST("/login", authCtrl.Login)
		}

		// WebSocket 路由（需要认证）
		api.GET("/ws", middleware.WSAuthMiddleware(userService), wsCtrl.HandleWebSocket)

		// 需要认证的路由
		authenticated := api.Group("")
		authenticated.Use(middleware.AuthMiddleware(userService))
		{
			// 用户路由
			users := authenticated.Group("/users")
			{
				users.GET("", userCtrl.GetAllUsers)
				users.GET("/online", userCtrl.GetOnlineUsers)
			}

			// 密钥路由
			keys := authenticated.Group("/keys")
			{
				keys.POST("/generate", keyCtrl.GenerateKeys)
				keys.POST("/upload", keyCtrl.UploadPublicKey)
				keys.GET("/:userID", keyCtrl.GetPublicKey)
			}

			// 消息路由
			messages := authenticated.Group("/messages")
			{
				messages.POST("/send", messageCtrl.SendMessage)
				messages.GET("/unread", messageCtrl.GetUnreadMessages)
				messages.POST("/:messageID/read", messageCtrl.MarkMessageAsRead)
			}
		}
	}

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}
