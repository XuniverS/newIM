package service

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"im-system/client/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocketService WebSocket服务
type WebSocketService struct {
	serverService ServerService
	cryptoService CryptoService
	clients       map[*websocket.Conn]*ClientInfo
	clientsMutex  sync.RWMutex
	upgrader      websocket.Upgrader
}

// ClientInfo 客户端信息
type ClientInfo struct {
	Token      string
	PrivateKey string
	ServerConn *websocket.Conn
}

// NewWebSocketService 创建WebSocket服务实例
func NewWebSocketService(serverService ServerService, cryptoService CryptoService) *WebSocketService {
	return &WebSocketService{
		serverService: serverService,
		cryptoService: cryptoService,
		clients:       make(map[*websocket.Conn]*ClientInfo),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// HandleWebSocket 处理WebSocket连接
func (s *WebSocketService) HandleWebSocket(c *gin.Context) {
	// 从查询参数或header获取token和privateKey
	token := c.Query("token")
	if token == "" {
		token = c.GetHeader("Authorization")
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
	}

	privateKey := c.Query("privateKey")
	if privateKey == "" {
		privateKey = c.GetHeader("X-Private-Key")
	}

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	// 升级到WebSocket
	clientConn, err := s.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// 连接到服务端WebSocket
	serverWSURL := s.serverService.GetServerWSURL() + "?token=" + token
	serverConn, _, err := websocket.DefaultDialer.Dial(serverWSURL, nil)
	if err != nil {
		log.Printf("Failed to connect to server WebSocket: %v", err)
		clientConn.Close()
		return
	}

	// 保存客户端信息
	clientInfo := &ClientInfo{
		Token:      token,
		PrivateKey: privateKey,
		ServerConn: serverConn,
	}

	s.clientsMutex.Lock()
	s.clients[clientConn] = clientInfo
	s.clientsMutex.Unlock()

	log.Printf("WebSocket connection established")

	// 启动消息转发
	go s.forwardFromClient(clientConn, serverConn, clientInfo)
	go s.forwardFromServer(serverConn, clientConn, clientInfo)
}

// forwardFromClient 从客户端转发消息到服务端
func (s *WebSocketService) forwardFromClient(clientConn, serverConn *websocket.Conn, info *ClientInfo) {
	defer func() {
		s.cleanup(clientConn, serverConn)
	}()

	for {
		var msg model.WSMessage
		err := clientConn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Client WebSocket error: %v", err)
			}
			return
		}

		// 如果是消息类型，需要加密
		if msg.Type == "message" && msg.Content != "" {
			// 获取接收者公钥
			publicKey, err := s.serverService.GetPublicKey(info.Token, msg.ReceiverID)
			if err != nil {
				log.Printf("Failed to get public key: %v", err)
				clientConn.WriteJSON(model.WSMessage{
					Type:    "error",
					Content: "Failed to get receiver's public key",
				})
				continue
			}

			// 加密消息
			encrypted, err := s.cryptoService.Encrypt(publicKey, msg.Content)
			if err != nil {
				log.Printf("Failed to encrypt message: %v", err)
				clientConn.WriteJSON(model.WSMessage{
					Type:    "error",
					Content: "Failed to encrypt message",
				})
				continue
			}

			msg.Content = encrypted
		}

		// 转发到服务端
		if err := serverConn.WriteJSON(msg); err != nil {
			log.Printf("Failed to forward to server: %v", err)
			return
		}
	}
}

// forwardFromServer 从服务端转发消息到客户端
func (s *WebSocketService) forwardFromServer(serverConn, clientConn *websocket.Conn, info *ClientInfo) {
	defer func() {
		s.cleanup(clientConn, serverConn)
	}()

	for {
		_, message, err := serverConn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Server WebSocket error: %v", err)
			}
			return
		}

		var msg model.WSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		// 如果是消息类型且有私钥，需要解密
		if msg.Type == "message" && msg.Content != "" && info.PrivateKey != "" {
			decrypted, err := s.cryptoService.Decrypt(info.PrivateKey, msg.Content)
			if err != nil {
				log.Printf("Failed to decrypt message: %v", err)
				// 即使解密失败，也转发原始消息
			} else {
				msg.Content = decrypted
			}
		}

		// 转发到客户端
		if err := clientConn.WriteJSON(msg); err != nil {
			log.Printf("Failed to forward to client: %v", err)
			return
		}
	}
}

// cleanup 清理连接
func (s *WebSocketService) cleanup(clientConn, serverConn *websocket.Conn) {
	s.clientsMutex.Lock()
	delete(s.clients, clientConn)
	s.clientsMutex.Unlock()

	clientConn.Close()
	serverConn.Close()
	log.Printf("WebSocket connection closed")
}
