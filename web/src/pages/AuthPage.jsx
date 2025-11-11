import React, { useState } from 'react'
import axios from 'axios'
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
      const endpoint = isLogin ? '/api/auth/login' : '/api/auth/register'
      const response = await axios.post(endpoint, {
        username,
        password,
      })

      const { token, user_id, username: userName } = response.data

      // 如果是注册，需要生成密钥
      if (!isLogin) {
        try {
          const keysResponse = await axios.post(
            '/api/keys/generate',
            {},
            {
              headers: {
                Authorization: `Bearer ${token}`,
              },
            }
          )
          localStorage.setItem('privateKey', keysResponse.data.private_key)
        } catch (err) {
          console.error('Failed to generate keys:', err)
        }
      } else {
        // 登录时检查是否有私钥，如果没有则生成
        try {
          const keysResponse = await axios.post(
            '/api/keys/generate',
            {},
            {
              headers: {
                Authorization: `Bearer ${token}`,
              },
            }
          )
          localStorage.setItem('privateKey', keysResponse.data.private_key)
        } catch (err) {
          // 可能已经存在密钥，忽略错误
          console.log('Keys already exist or error:', err.message)
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
