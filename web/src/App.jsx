import React, { useState, useEffect } from 'react'
import './App.css'
import AuthPage from './pages/AuthPage'
import ChatPage from './pages/ChatPage'
function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(false)
  const [user, setUser] = useState(null)

  useEffect(() => {
    // 检查本地存储中是否有 token
    const token = localStorage.getItem('token')
    const userData = localStorage.getItem('user')
    if (token && userData) {
      setIsLoggedIn(true)
      setUser(JSON.parse(userData))
    }
  }, [])

  const handleLogin = (token, userData) => {
    localStorage.setItem('token', token)
    localStorage.setItem('user', JSON.stringify(userData))
    setIsLoggedIn(true)
    setUser(userData)
  }

  const handleLogout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    localStorage.removeItem('privateKey')
    setIsLoggedIn(false)
    setUser(null)
  }

  return (
    <div className="app">
      {isLoggedIn ? (
        <ChatPage user={user} onLogout={handleLogout} />
      ) : (
        <AuthPage onLogin={handleLogin} />
      )}
    </div>
  )
}

export default App
