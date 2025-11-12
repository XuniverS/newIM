import React, { useState } from 'react'
import { authAPI, keyAPI } from '../services/api'
import './AuthPage.css'
function AuthPage({ onLogin }) {
  const [isLogin, setIsLogin] = useState(true)
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      // 调用客户端后端API
      const apiCall = isLogin ? authAPI.login : authAPI.register
      const response = await apiCall(username, password)

      const { token, user_id, username: userName } = response.data

      // 在客户端生成密钥对（更安全，私钥不经过网络传输）
      if (!isLogin) {
        // 注册时生成新密钥对
        try {
          const { generateKeyPair } = await import('../utils/crypto')
          const { publicKey, privateKey } = await generateKeyPair()
          
          // 保存私钥到本地
          localStorage.setItem('privateKey', privateKey)
          
          // 上传公钥到服务端
          await keyAPI.uploadPublicKey(publicKey)
        } catch (err) {
          console.error('Failed to generate keys:', err)
          setError('密钥生成失败，请重试')
          return
        }
      } else {
        // 登录时检查是否有私钥
        const existingPrivateKey = localStorage.getItem('privateKey')
        if (!existingPrivateKey) {
          // 如果没有私钥，生成新的
          try {
            const { generateKeyPair } = await import('../utils/crypto')
            const { publicKey, privateKey } = await generateKeyPair()
            localStorage.setItem('privateKey', privateKey)
            await keyAPI.uploadPublicKey(publicKey)
          } catch (err) {
            console.error('Failed to generate keys:', err)
          }
        }
      }

      onLogin(token, { user_id, username: userName })
      setUsername('')
      setPassword('')
    } catch (err) {
      setError(err.response?.data?.error || 'An error occurred')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="auth-container">
      <div className="auth-card">
        <div className="auth-header">
          <h1>💬 IM 即时通讯</h1>
          <p>安全、快速、可靠的消息传递</p>
        </div>

        <form onSubmit={handleSubmit} className="auth-form">
          <div className="form-group">
            <label htmlFor="username">用户名</label>
            <input
              id="username"
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="输入用户名"
              required
              disabled={loading}
            />
          </div>

          <div className="form-group">
            <label htmlFor="password">密码</label>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="输入密码"
              required
              disabled={loading}
            />
          </div>

          {error && <div className="error-message">{error}</div>}

          <button type="submit" className="submit-btn" disabled={loading}>
            {loading ? '处理中...' : isLogin ? '登录' : '注册'}
          </button>
        </form>

        <div className="auth-footer">
          <p>
            {isLogin ? '还没有账户？' : '已有账户？'}
            <button
              type="button"
              className="toggle-btn"
              onClick={() => {
                setIsLogin(!isLogin)
                setError('')
              }}
              disabled={loading}
            >
              {isLogin ? '注册' : '登录'}
            </button>
          </p>
        </div>

        <div className="auth-info">
          <h3>🔐 安全特性</h3>
          <ul>
            <li>✓ RSA 端到端加密</li>
            <li>✓ 离线消息队列</li>
            <li>✓ 实时 WebSocket 通信</li>
            <li>✓ 密钥自动生成与管理</li>
          </ul>
        </div>
      </div>
    </div>
  )
}

export default AuthPage
