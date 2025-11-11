package main

import (
	"fmt"
	"log"
	"os"

	"im-system/internal/config"
	"im-system/internal/db"
	"im-system/internal/kafka"
	"im-system/internal/server"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库
	database, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// 初始化 Kafka
	kafkaProducer, err := kafka.InitProducer(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Kafka: %v", err)
	}
	defer kafkaProducer.Close()

	// 创建服务器
	srv := server.NewServer(database, kafkaProducer, cfg)

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 IM Server starting on port %s\n", port)
	if err := srv.Run(":" + port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
