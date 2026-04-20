import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../api'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref(null)

  const setToken = (newToken) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  const setUserInfo = (info) => {
    userInfo.value = info
  }

  const logout = () => {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
  }

  const login = async (username, password) => {
    const res = await api.post('/api/user/login', { username, password })
    if (res.data.code === 200) {
      setToken(res.data.data.token)
      await fetchUserInfo()
      return true
    }
    throw new Error(res.data.message)
  }

  const register = async (username, password, email) => {
    const res = await api.post('/api/user/register', { username, password, email })
    return res.data.code === 200
  }

  const fetchUserInfo = async () => {
    try {
      const res = await api.get('/api/user/info')
      if (res.data.code === 200) {
        setUserInfo(res.data.data)
      }
    } catch (e) {
      console.error('获取用户信息失败', e)
    }
  }

  return {
    token,
    userInfo,
    login,
    register,
    logout,
    fetchUserInfo,
    setToken
  }
})
