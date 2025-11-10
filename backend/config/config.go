package config

// Config 应用配置
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Queue    QueueConfig
	Auth     AuthConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string
	Port int
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// QueueConfig 消息队列配置
type QueueConfig struct {
	Type string // redis, rabbitmq, etc.
	Host string
	Port int
}

// AuthConfig 认证配置
type AuthConfig struct {
	JWTSecret string
	TokenTTL  int // seconds
}
