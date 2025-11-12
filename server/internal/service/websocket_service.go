package service

import (
	"log"
	"sync"
	"time"

	"im-system/server/internal/model"

	"github.com/gorilla/websocket"
)

// WebSocketService WebSocket 服务接口
type WebSocketService interface {
	RegisterClient(userID int, username string, conn *websocket.Conn) *model.WSClient
	UnregisterClient(userID int)
	GetClient(userID int) *model.WSClient
	GetOnlineUsers() []int
	HandleMessage(client *model.WSClient, msg model.WSMessage)
	ReadPump(client *model.WSClient, conn *websocket.Conn)
	WritePump(client *model.WSClient, conn *websocket.Conn)
}

type websocketService struct {
	clients        map[int]*model.WSClient
	clientsMutex   sync.RWMutex
	messageService MessageService
	userService    UserService
}

// NewWebSocketService 创建 WebSocket 服务实例
func NewWebSocketService(messageService MessageService, userService UserService) WebSocketService {
	return &websocketService{
		clients:        make(map[int]*model.WSClient),
		messageService: messageService,
		userService:    userService,
	}
}

func (s *websocketService) RegisterClient(userID int, username string, conn *websocket.Conn) *model.WSClient {
	client := &model.WSClient{
		UserID:   userID,
		Username: username,
		Send:     make(chan interface{}, 256),
	}

	s.clientsMutex.Lock()
	s.clients[userID] = client
	s.clientsMutex.Unlock()

	log.Printf("User %s (ID: %d) connected", username, userID)
	return client
}

func (s *websocketService) UnregisterClient(userID int) {
	s.clientsMutex.Lock()
	if client, ok := s.clients[userID]; ok {
		close(client.Send)
		delete(s.clients, userID)
		log.Printf("User ID %d disconnected", userID)
	}
	s.clientsMutex.Unlock()
}

func (s *websocketService) GetClient(userID int) *model.WSClient {
	s.clientsMutex.RLock()
	defer s.clientsMutex.RUnlock()
	return s.clients[userID]
}

func (s *websocketService) GetOnlineUsers() []int {
	s.clientsMutex.RLock()
	defer s.clientsMutex.RUnlock()

	users := make([]int, 0, len(s.clients))
	for userID := range s.clients {
		users = append(users, userID)
	}
	return users
}

func (s *websocketService) HandleMessage(client *model.WSClient, msg model.WSMessage) {
	// 保存消息到数据库
	messageID, err := s.messageService.SendMessage(client.UserID, msg.ReceiverID, msg.Content)
	if err != nil {
		client.Send <- model.WSMessage{
			Type:    "error",
			Content: "Failed to save message",
		}
		return
	}

	// 检查接收者是否在线
	receiver := s.GetClient(msg.ReceiverID)
	if receiver != nil {
		// 在线，直接发送
		receiver.Send <- model.WSMessage{
			Type:      "message",
			SenderID:  client.UserID,
			Content:   msg.Content,
			MessageID: messageID,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		log.Printf("Message sent from %d to %d (online)", client.UserID, msg.ReceiverID)
	} else {
		// 离线，消息已保存到数据库
		log.Printf("Message saved for offline user %d", msg.ReceiverID)
	}

	// 发送确认
	client.Send <- model.WSMessage{
		Type:      "message_sent",
		Content:   "Message sent successfully",
		MessageID: messageID,
	}
}

func (s *websocketService) ReadPump(client *model.WSClient, conn *websocket.Conn) {
	defer func() {
		s.UnregisterClient(client.UserID)
		conn.Close()
	}()

	for {
		var msg model.WSMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			return
		}

		switch msg.Type {
		case "message":
			s.HandleMessage(client, msg)
		case "ping":
			client.Send <- model.WSMessage{Type: "pong"}
		}
	}
}

func (s *websocketService) WritePump(client *model.WSClient, conn *websocket.Conn) {
	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := conn.WriteJSON(message); err != nil {
				return
			}
		}
	}
}
