package main

import (
	"fmt"
	"log"
	"os"

	"im-system/client/internal/config"
	"im-system/client/internal/controller"
	"im-system/client/internal/service"
	"im-system/client/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化日志
	logger.Init()
	logger.Info("Starting IM Client...")

	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化服务层
	serverService := service.NewServerService(cfg)
	cryptoService := service.NewCryptoService()
	wsService := service.NewWebSocketService(serverService, cryptoService)

	// 初始化控制器
	authCtrl := controller.NewAuthController(serverService, cryptoService)
	messageCtrl := controller.NewMessageController(serverService, wsService, cryptoService)
	userCtrl := controller.NewUserController(serverService)
	keyCtrl := controller.NewKeyController(serverService, cryptoService)

	// 设置路由
	router := setupRouter(authCtrl, messageCtrl, userCtrl, keyCtrl, wsService)

	// 启动服务器
	port := os.Getenv("CLIENT_PORT")
	if port == "" {
		port = "3001"
	}

	logger.Info(fmt.Sprintf("🚀 IM Client Backend starting on port %s", port))
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func setupRouter(
	authCtrl *controller.AuthController,
	messageCtrl *controller.MessageController,
	userCtrl *controller.UserController,
	keyCtrl *controller.KeyController,
	wsService *service.WebSocketService,
) *gin.Engine {
	router := gin.Default()

	// CORS 中间件
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 静态文件服务（前端）
	router.Static("/assets", "./web/dist/assets")
	router.StaticFile("/", "./web/dist/index.html")
	router.StaticFile("/index.html", "./web/dist/index.html")

	// API 路由
	api := router.Group("/api")
	{
		// 认证
		api.POST("/auth/register", authCtrl.Register)
		api.POST("/auth/login", authCtrl.Login)

		// WebSocket
		api.GET("/ws", wsService.HandleWebSocket)

		// 用户
		api.GET("/users", userCtrl.GetAllUsers)
		api.GET("/users/online", userCtrl.GetOnlineUsers)

		// 密钥
		api.POST("/keys/generate", keyCtrl.GenerateKeys)
		api.GET("/keys/:userID", keyCtrl.GetPublicKey)

		// 消息
		api.POST("/messages/send", messageCtrl.SendMessage)
		api.GET("/messages/unread", messageCtrl.GetUnreadMessages)
	}

	return router
}
