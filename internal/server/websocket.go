package server

import (
	"log"

	"im-system/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WSMessage struct {
	Type      string `json:"type"`
	ReceiverID int    `json:"receiver_id,omitempty"`
	Content   string `json:"content,omitempty"`
	MessageID int    `json:"message_id,omitempty"`
}

func (s *Server) handleWebSocket(c *gin.Context) {
	userID := s.getUserIDFromContext(c)
	username := s.getUsernameFromContext(c)

	conn, err := s.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		UserID:   userID,
		Username: username,
		Conn:     conn,
		Send:     make(chan interface{}, 256),
	}

	s.RegisterClient(client)
	defer func() {
		s.UnregisterClient(userID)
		conn.Close()
	}()

	// 启动读取 goroutine
	go s.readPump(client)
	// 启动写入 goroutine
	go s.writePump(client)
}

func (s *Server) readPump(client *Client) {
	defer func() {
		client.Conn.Close()
	}()

	for {
		var msg WSMessage
		err := client.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			return
		}

		switch msg.Type {
		case "message":
			s.handleMessageFromClient(client, msg)
		case "ping":
			client.Send <- WSMessage{Type: "pong"}
		}
	}
}

func (s *Server) writePump(client *Client) {
	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := client.Conn.WriteJSON(message); err != nil {
				return
			}
		}
	}
}

func (s *Server) handleMessageFromClient(client *Client, msg WSMessage) {
	// 获取接收者的公钥
	_, err := db.GetPublicKey(s.db, msg.ReceiverID)
	if err != nil {
		client.Send <- WSMessage{
			Type: "error",
			Content: "Receiver public key not found",
		}
		return
	}

	// 保存消息到数据库
	_, err = db.SaveMessage(s.db, client.UserID, msg.ReceiverID, msg.Content)
	if err != nil {
		client.Send <- WSMessage{
			Type: "error",
			Content: "Failed to save message",
		}
		return
	}

	// 检查接收者是否在线
	receiver := s.GetClient(msg.ReceiverID)
	if receiver != nil {
		// 在线，直接发送
		receiver.Send <- WSMessage{
			Type:    "message",
			Content: msg.Content,
		}
	} else {
		// 离线，放入 Kafka 队列
		// TODO: 实现 Kafka 消息队列
	}

	// 发送确认
	client.Send <- WSMessage{
		Type: "message_sent",
		Content: "Message sent successfully",
	}
}
