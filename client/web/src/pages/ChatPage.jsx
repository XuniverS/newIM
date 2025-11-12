import React, { useState, useEffect, useRef } from 'react'
import { userAPI } from '../services/api'
import './ChatPage.css'
import UserList from '../components/UserList'
import ChatWindow from '../components/ChatWindow'
function ChatPage({ user, onLogout }) {
  const [users, setUsers] = useState([])
  const [selectedUser, setSelectedUser] = useState(null)
  const [messages, setMessages] = useState({})
  const [onlineUsers, setOnlineUsers] = useState([])
  const [loading, setLoading] = useState(true)
  const wsRef = useRef(null)
  const token = localStorage.getItem('token')

  useEffect(() => {
    // 初始化 WebSocket 连接
    initWebSocket()
    // 获取在线用户列表
    fetchOnlineUsers()

    return () => {
      if (wsRef.current) {
        wsRef.current.close()
      }
    }
  }, [])

  const initWebSocket = () => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const privateKey = localStorage.getItem('privateKey') || ''
    const wsUrl = `${protocol}//${window.location.host}/api/ws?token=${token}&privateKey=${encodeURIComponent(privateKey)}`

    wsRef.current = new WebSocket(wsUrl)

    wsRef.current.onopen = () => {
      console.log('WebSocket connected')
      setLoading(false)
    }

    wsRef.current.onmessage = (event) => {
      const message = JSON.parse(event.data)
      handleWebSocketMessage(message)
    }

    wsRef.current.onerror = (error) => {
      console.error('WebSocket error:', error)
    }

    wsRef.current.onclose = () => {
      console.log('WebSocket disconnected')
      // 尝试重新连接
      setTimeout(initWebSocket, 3000)
    }
  }

  const handleWebSocketMessage = (message) => {
    if (message.type === 'message') {
      // 接收消息（客户端后端已经解密）
      const senderID = message.sender_id || message.from_user_id
      
      setMessages((prev) => ({
        ...prev,
        [senderID]: [...(prev[senderID] || []), message],
      }))
    } else if (message.type === 'pong') {
      // 心跳响应
      console.log('Pong received')
    }
  }

  const fetchOnlineUsers = async () => {
    try {
      const response = await userAPI.getOnlineUsers()
      setOnlineUsers(response.data.online_users || [])
    } catch (err) {
      console.error('Failed to fetch online users:', err)
    }
  }

  const sendMessage = (receiverID, content) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      // 发送明文消息（客户端后端会加密）
      wsRef.current.send(
        JSON.stringify({
          type: 'message',
          receiver_id: receiverID,
          content: content,
        })
      )

      // 本地添加消息
      setMessages((prev) => ({
        ...prev,
        [receiverID]: [
          ...(prev[receiverID] || []),
          {
            type: 'message',
            content: content,
            sender_id: user.user_id,
            timestamp: new Date().toISOString(),
            is_own: true,
          },
        ],
      }))
    }
  }

  const handleLogout = () => {
    if (wsRef.current) {
      wsRef.current.close()
    }
    onLogout()
  }

  return (
    <div className="chat-page">
      <div className="chat-container">
        <UserList
          currentUser={user}
          onlineUsers={onlineUsers}
          onSelectUser={setSelectedUser}
          selectedUser={selectedUser}
          onLogout={handleLogout}
        />
        <ChatWindow
          selectedUser={selectedUser}
          messages={messages[selectedUser?.id] || []}
          currentUser={user}
          onSendMessage={sendMessage}
          loading={loading}
        />
      </div>
    </div>
  )
}

export default ChatPage
