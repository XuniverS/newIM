import React, { useState, useEffect } from 'react'
import axios from 'axios'
import './UserList.css'
function UserList({ currentUser, onlineUsers, onSelectUser, selectedUser, onLogout }) {
  const [allUsers, setAllUsers] = useState([])
  const [searchTerm, setSearchTerm] = useState('')
  const [loading, setLoading] = useState(true)
  const token = localStorage.getItem('token')

  useEffect(() => {
    // 这里可以从服务器获取所有用户列表
    // 目前我们使用在线用户列表
    setAllUsers(onlineUsers)
    setLoading(false)
  }, [onlineUsers])

  const filteredUsers = allUsers.filter((userId) => {
    const userStr = userId.toString()
    return userStr.includes(searchTerm) && userId !== currentUser.user_id
  })

  return (
    <div className="user-list">
      <div className="user-list-header">
        <div className="current-user-info">
          <div className="user-avatar">
            {currentUser.username.charAt(0).toUpperCase()}
          </div>
          <div className="user-details">
            <div className="user-name">{currentUser.username}</div>
            <div className="user-status">在线</div>
          </div>
        </div>
        <button className="logout-btn" onClick={onLogout} title="退出登录">
          🚪
        </button>
      </div>

      <div className="user-list-search">
        <input
          type="text"
          placeholder="搜索用户..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="search-input"
        />
      </div>

      <div className="user-list-content">
        <div className="online-users-section">
          <div className="section-title">
            在线用户 ({filteredUsers.length})
          </div>
          {loading ? (
            <div className="loading">加载中...</div>
          ) : filteredUsers.length === 0 ? (
            <div className="empty-state">
              <p>暂无在线用户</p>
              <small>邀请朋友加入聊天</small>
            </div>
          ) : (
            <div className="users-list">
              {filteredUsers.map((userId) => (
                <div
                  key={userId}
                  className={`user-item ${selectedUser?.id === userId ? 'active' : ''}`}
                  onClick={() => onSelectUser({ id: userId, username: `User ${userId}` })}
                >
                  <div className="user-item-avatar">
                    {`User ${userId}`.charAt(0).toUpperCase()}
                  </div>
                  <div className="user-item-info">
                    <div className="user-item-name">User {userId}</div>
                    <div className="user-item-status">在线</div>
                  </div>
                  <div className="user-item-indicator"></div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default UserList
