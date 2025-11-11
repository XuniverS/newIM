package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	KafkaHost  string
	KafkaPort  string
	JWTSecret  string
	RedisHost  string
	RedisPort  string
}

func LoadConfig() *Config {
	// 尝试加载 .env 文件
	_ = godotenv.Load()

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
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
