package server

import (
	"database/sql"
	"net/http"
	"sync"

	"im-system/internal/config"
	"im-system/internal/kafka"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Server struct {
	db              *sql.DB
	kafkaProducer   *kafka.Producer
	config          *config.Config
	clients         map[int]*Client
	clientsMutex    sync.RWMutex
	upgrader        websocket.Upgrader
	broadcastChan   chan interface{}
}

type Client struct {
	UserID   int
	Username string
	Conn     *websocket.Conn
	Send     chan interface{}
}

func NewServer(db *sql.DB, kafkaProducer *kafka.Producer, cfg *config.Config) *Server {
	return &Server{
		db:            db,
		kafkaProducer: kafkaProducer,
		config:        cfg,
		clients:       make(map[int]*Client),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // 允许跨域
			},
		},
		broadcastChan: make(chan interface{}, 256),
	}
}

func (s *Server) Run(addr string) error {
	router := gin.Default()

	// 添加 CORS 中间件
	router.Use(corsMiddleware())

	// 静态文件服务
	router.Static("/static", "./web/dist")
	router.StaticFile("/", "./web/dist/index.html")

	// API 路由
	api := router.Group("/api")
	{
		// 认证相关
		api.POST("/auth/register", s.handleRegister)
		api.POST("/auth/login", s.handleLogin)

		// WebSocket 连接
		api.GET("/ws", s.wsAuthMiddleware, s.handleWebSocket)

		// 公钥相关
		api.POST("/keys/generate", s.authMiddleware(), s.handleGenerateKeys)
		api.GET("/keys/:userID", s.authMiddleware(), s.handleGetPublicKey)
		api.POST("/keys/upload", s.authMiddleware(), s.handleUploadPublicKey)

		// 消息相关
		api.POST("/messages/send", s.authMiddleware(), s.handleSendMessage)
		api.GET("/messages/unread", s.authMiddleware(), s.handleGetUnreadMessages)
		api.POST("/messages/:messageID/read", s.authMiddleware(), s.handleMarkMessageAsRead)

		// 用户相关
		api.GET("/users/online", s.authMiddleware(), s.handleGetOnlineUsers)
	}

	return router.Run(addr)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (s *Server) RegisterClient(client *Client) {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()
	s.clients[client.UserID] = client
}

func (s *Server) UnregisterClient(userID int) {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()
	delete(s.clients, userID)
}

func (s *Server) GetClient(userID int) *Client {
	s.clientsMutex.RLock()
	defer s.clientsMutex.RUnlock()
	return s.clients[userID]
}

func (s *Server) GetOnlineUsers() []int {
	s.clientsMutex.RLock()
	defer s.clientsMutex.RUnlock()
	users := make([]int, 0, len(s.clients))
	for userID := range s.clients {
		users = append(users, userID)
	}
	return users
}
