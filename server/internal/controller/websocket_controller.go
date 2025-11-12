package controller

import (
	"log"
	"net/http"

	"im-system/server/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocketController WebSocket 控制器
type WebSocketController struct {
	wsService service.WebSocketService
	upgrader  websocket.Upgrader
}

// NewWebSocketController 创建 WebSocket 控制器实例
func NewWebSocketController(wsService service.WebSocketService) *WebSocketController {
	return &WebSocketController{
		wsService: wsService,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // 允许跨域
			},
		},
	}
}

// HandleWebSocket 处理 WebSocket 连接
func (ctrl *WebSocketController) HandleWebSocket(c *gin.Context) {
	userID := getUserIDFromContext(c)
	username := getUsernameFromContext(c)

	conn, err := ctrl.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := ctrl.wsService.RegisterClient(userID, username, conn)
	defer func() {
		ctrl.wsService.UnregisterClient(userID)
		conn.Close()
	}()

	// 启动读取和写入 goroutine
	go ctrl.wsService.WritePump(client, conn)
	ctrl.wsService.ReadPump(client, conn)
}
