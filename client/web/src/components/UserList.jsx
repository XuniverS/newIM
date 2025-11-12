import React, { useState, useEffect } from 'react'
import { userAPI } from '../services/api'
import './UserList.css'
function UserList({ currentUser, onlineUsers, onSelectUser, selectedUser, onLogout }) {
  const [allUsers, setAllUsers] = useState([])
  const [searchTerm, setSearchTerm] = useState('')
  const [loading, setLoading] = useState(true)
  const token = localStorage.getItem('token')

  useEffect(() => {
    fetchAllUsers()
  }, [])

  const fetchAllUsers = async () => {
    try {
      const response = await userAPI.getAllUsers()
      setAllUsers(response.data.users || [])
    } catch (err) {
      console.error('Failed to fetch users:', err)
    } finally {
      setLoading(false)
    }
  }

  const filteredUsers = allUsers.filter((user) => {
    const username = user.username.toLowerCase()
    const search = searchTerm.toLowerCase()
    return username.includes(search)
  })

  const isUserOnline = (userId) => {
    return onlineUsers.includes(userId)
  }

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
            所有用户 ({filteredUsers.length})
          </div>
          {loading ? (
            <div className="loading">加载中...</div>
          ) : filteredUsers.length === 0 ? (
            <div className="empty-state">
              <p>暂无用户</p>
              <small>邀请朋友加入聊天</small>
            </div>
          ) : (
            <div className="users-list">
              {filteredUsers.map((user) => {
                const online = isUserOnline(user.id)
                return (
                  <div
                    key={user.id}
                    className={`user-item ${selectedUser?.id === user.id ? 'active' : ''}`}
                    onClick={() => onSelectUser(user)}
                  >
                    <div className="user-item-avatar">
                      {user.username.charAt(0).toUpperCase()}
                    </div>
                    <div className="user-item-info">
                      <div className="user-item-name">{user.username}</div>
                      <div className={`user-item-status ${online ? 'online' : 'offline'}`}>
                        {online ? '在线' : '离线'}
                      </div>
                    </div>
                    {online && <div className="user-item-indicator"></div>}
                  </div>
                )
              })}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default UserList
