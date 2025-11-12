package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// 服务端配置
	ServerHost string
	ServerPort string

	// 客户端配置
	ClientPort string
}

// Load 加载配置
func Load() (*Config, error) {
	// 尝试加载 .env 文件
	_ = godotenv.Load("../../.env")
	_ = godotenv.Load(".env")

	return &Config{
		ServerHost: getEnv("SERVER_HOST", "localhost"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		ClientPort: getEnv("CLIENT_PORT", "3001"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetServerURL 获取服务端URL
func (c *Config) GetServerURL() string {
	return "http://" + c.ServerHost + ":" + c.ServerPort
}

// GetServerWSURL 获取服务端WebSocket URL
func (c *Config) GetServerWSURL() string {
	return "ws://" + c.ServerHost + ":" + c.ServerPort
}
