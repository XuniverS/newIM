package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// 数据库配置
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Kafka 配置
	KafkaHost string
	KafkaPort string

	// JWT 配置
	JWTSecret string

	// Redis 配置
	RedisHost string
	RedisPort string

	// 服务器配置
	ServerPort string
}

// Load 加载配置
func Load() (*Config, error) {
	// 尝试加载 .env 文件
	_ = godotenv.Load("../../.env")
	_ = godotenv.Load(".env")

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "im_db"),
		KafkaHost:  getEnv("KAFKA_HOST", "localhost"),
		KafkaPort:  getEnv("KAFKA_PORT", "9092"),
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		RedisHost:  getEnv("REDIS_HOST", "localhost"),
		RedisPort:  getEnv("REDIS_PORT", "6379"),
		ServerPort: getEnv("PORT", "8080"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
