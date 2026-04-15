import { defineStore } from 'pinia'
import { ref } from 'vue'
import { authApi } from '../utils/request'

export const useUserStore = defineStore('user', () => {
  const token = ref(uni.getStorageSync('token') || '')
  const userInfo = ref(uni.getStorageSync('userInfo') || null)
  const isLoggedIn = ref(!!token.value)

  const login = async (loginData) => {
    const res = await authApi.login(loginData)
    token.value = res.token
    userInfo.value = res.user
    isLoggedIn.value = true
    uni.setStorageSync('token', res.token)
    uni.setStorageSync('userInfo', res.user)
    return res
  }

  const register = async (registerData) => {
    return await authApi.register(registerData)
  }

  const logout = () => {
    token.value = ''
    userInfo.value = null
    isLoggedIn.value = false
    uni.removeStorageSync('token')
    uni.removeStorageSync('userInfo')
    uni.reLaunch({ url: '/pages/login/index' })
  }

  const getProfile = async () => {
    const res = await authApi.getProfile()
    userInfo.value = res
    uni.setStorageSync('userInfo', res)
    return res
  }

  return { token, userInfo, isLoggedIn, login, register, logout, getProfile }
})
