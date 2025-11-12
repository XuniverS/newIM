import axios from 'axios'

// 客户端后端API地址
// 开发环境：http://localhost:3001
// 生产环境：https://client.yourdomain.com（通过环境变量配置）
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:3001'

// 创建axios实例
const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
})

// 请求拦截器 - 添加token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }

    // 如果需要私钥（用于解密）
    const privateKey = localStorage.getItem('privateKey')
    if (privateKey && config.headers['X-Need-Private-Key']) {
      config.headers['X-Private-Key'] = privateKey
      delete config.headers['X-Need-Private-Key']
    }

    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器 - 处理错误
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      // Token过期，清除本地存储并跳转到登录页
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      localStorage.removeItem('privateKey')
      window.location.href = '/'
    }
    return Promise.reject(error)
  }
)

// 认证API
export const authAPI = {
  register: (username, password) =>
    api.post('/api/auth/register', { username, password }),
  login: (username, password) =>
    api.post('/api/auth/login', { username, password }),
}

// 用户API
export const userAPI = {
  getAllUsers: () => api.get('/api/users'),
  getOnlineUsers: () => api.get('/api/users/online'),
}

// 密钥API
export const keyAPI = {
  uploadPublicKey: (publicKey) => api.post('/api/keys/upload', { public_key: publicKey }),
  getPublicKey: (userID) => api.get(`/api/keys/${userID}`),
}

// 消息API
export const messageAPI = {
  sendMessage: (receiverID, content) =>
    api.post('/api/messages/send', { receiver_id: receiverID, content }),
  getUnreadMessages: () =>
    api.get('/api/messages/unread', {
      headers: { 'X-Need-Private-Key': 'true' },
    }),
}

export default api
