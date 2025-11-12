package main

import (
	"fmt"
	"log"
	"os"

	"im-system/server/internal/config"
	"im-system/server/internal/repository"
	"im-system/server/internal/router"
	"im-system/server/internal/service"
	"im-system/server/pkg/logger"
)

func main() {
	// 初始化日志
	logger.Init()
	logger.Info("Starting IM Server...")

	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	db, err := repository.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// 初始化 Repository 层
	userRepo := repository.NewUserRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	keyRepo := repository.NewKeyRepository(db)

	// 初始化 Service 层
	userService := service.NewUserService(userRepo, cfg)
	messageService := service.NewMessageService(messageRepo, userRepo)
	keyService := service.NewKeyService(keyRepo, userRepo)
	wsService := service.NewWebSocketService(messageService, userService)

	// 初始化路由
	r := router.SetupRouter(cfg, userService, messageService, keyService, wsService)

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info(fmt.Sprintf("🚀 IM Server starting on port %s", port))
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
