import React, { useState, useEffect, useRef } from 'react'
import './ChatWindow.css'
function ChatWindow({ selectedUser, messages, currentUser, onSendMessage, loading }) {
  const [inputValue, setInputValue] = useState('')
  const [isSending, setIsSending] = useState(false)
  const messagesEndRef = useRef(null)

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }

  useEffect(() => {
    scrollToBottom()
  }, [messages])

  const handleSendMessage = async (e) => {
    e.preventDefault()
    if (!inputValue.trim() || !selectedUser || isSending) return

    setIsSending(true)
    try {
      onSendMessage(selectedUser.id, inputValue)
      setInputValue('')
    } catch (err) {
      console.error('Failed to send message:', err)
    } finally {
      setIsSending(false)
    }
  }

  if (!selectedUser) {
    return (
      <div className="chat-window empty">
        <div className="empty-state">
          <div className="empty-icon">💬</div>
          <h2>选择一个用户开始聊天</h2>
          <p>从左侧用户列表中选择一个在线用户，开始发送消息</p>
        </div>
      </div>
    )
  }

  return (
    <div className="chat-window">
      <div className="chat-header">
        <div className="chat-header-info">
          <div className="chat-user-avatar">
            {selectedUser.username.charAt(0).toUpperCase()}
          </div>
          <div className="chat-user-details">
            <div className="chat-user-name">{selectedUser.username}</div>
            <div className="chat-user-status">在线</div>
          </div>
        </div>
        <div className="chat-header-actions">
          <button className="action-btn" title="视频通话">
            📹
          </button>
          <button className="action-btn" title="语音通话">
            📞
          </button>
          <button className="action-btn" title="更多选项">
            ⋮
          </button>
        </div>
      </div>

      <div className="chat-messages">
        {loading ? (
          <div className="loading-state">
            <div className="spinner"></div>
            <p>连接中...</p>
          </div>
        ) : messages.length === 0 ? (
          <div className="empty-messages">
            <div className="empty-icon">👋</div>
            <p>开始与 {selectedUser.username} 的对话</p>
          </div>
        ) : (
          <>
            {messages.map((msg, index) => (
              <div
                key={index}
                className={`message ${msg.is_own || msg.sender_id === currentUser.user_id ? 'own' : 'other'}`}
              >
                {!(msg.is_own || msg.sender_id === currentUser.user_id) && (
                  <div className="message-avatar">
                    {selectedUser.username.charAt(0).toUpperCase()}
                  </div>
                )}
                <div className="message-content">
                  <div className="message-text">{msg.content}</div>
                  <div className="message-time">
                    {new Date(msg.timestamp || msg.created_at).toLocaleTimeString('zh-CN', {
                      hour: '2-digit',
                      minute: '2-digit',
                    })}
                  </div>
                </div>
              </div>
            ))}
            <div ref={messagesEndRef} />
          </>
        )}
      </div>

      <form className="chat-input-area" onSubmit={handleSendMessage}>
        <div className="input-actions">
          <button type="button" className="action-btn" title="表情">
            😊
          </button>
          <button type="button" className="action-btn" title="文件">
            📎
          </button>
        </div>
        <input
          type="text"
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          placeholder="输入消息..."
          className="message-input"
          disabled={isSending || !selectedUser}
        />
        <button
          type="submit"
          className="send-btn"
          disabled={isSending || !inputValue.trim() || !selectedUser}
          title="发送"
        >
          {isSending ? '⏳' : '📤'}
        </button>
      </form>
    </div>
  )
}

export default ChatWindow
