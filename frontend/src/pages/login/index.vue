<template>
  <view class="login-container">
    <view class="login-header">
      <text class="title">🏠 租房管理系统</text>
      <text class="subtitle">轻松管理您的房源</text>
    </view>

    <view class="login-form">
      <view class="form-item">
        <text class="form-icon">👤</text>
        <input class="form-input" v-model="formData.username" placeholder="请输入账号" />
      </view>

      <view class="form-item">
        <text class="form-icon">🔒</text>
        <input class="form-input" :password="!showPassword" v-model="formData.password" placeholder="请输入密码" />
        <text class="form-icon" @click="showPassword = !showPassword">{{ showPassword ? '👁️' : '🙈' }}</text>
      </view>

      <button class="btn-primary login-btn" @click="handleLogin" :disabled="loading">
        {{ loading ? '登录中...' : '登录' }}
      </button>

      <view class="register-link" @click="goRegister">
        还没有账号？<text class="link">立即注册</text>
      </view>
    </view>

    <!-- 注册弹窗 -->
    <view class="register-modal" v-if="showRegister">
      <view class="register-content">
        <view class="popup-title">注册账号</view>
        <view class="form-item">
          <text class="form-label">账号</text>
          <input class="form-input" v-model="registerData.username" placeholder="请设置账号" />
        </view>
        <view class="form-item">
          <text class="form-label">姓名</text>
          <input class="form-input" v-model="registerData.name" placeholder="请输入姓名" />
        </view>
        <view class="form-item">
          <text class="form-label">密码</text>
          <input class="form-input" password v-model="registerData.password" placeholder="请设置密码(至少6位)" />
        </view>
        <view class="form-item">
          <text class="form-label">确认</text>
          <input class="form-input" password v-model="registerData.confirmPassword" placeholder="请再次输入密码" />
        </view>
        <button class="btn-primary" @click="handleRegister" :disabled="registerLoading">
          {{ registerLoading ? '注册中...' : '注册' }}
        </button>
        <button class="btn-default" @click="showRegister = false">取消</button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useUserStore } from '../../stores/user'

const userStore = useUserStore()
const loading = ref(false)
const registerLoading = ref(false)
const showPassword = ref(false)
const showRegister = ref(false)

const formData = reactive({ username: '', password: '' })
const registerData = reactive({ username: '', name: '', password: '', confirmPassword: '' })

const handleLogin = async () => {
  if (!formData.username) { uni.showToast({ title: '请输入账号', icon: 'none' }); return }
  if (!formData.password) { uni.showToast({ title: '请输入密码', icon: 'none' }); return }
  loading.value = true
  try {
    await userStore.login(formData)
    uni.showToast({ title: '登录成功', icon: 'success' })
    setTimeout(() => uni.switchTab({ url: '/pages/index/index' }), 1000)
  } catch (error) {
    console.error('登录失败', error)
  } finally {
    loading.value = false
  }
}

const goRegister = () => { showRegister.value = true }

const handleRegister = async () => {
  if (!registerData.username) { uni.showToast({ title: '请设置账号', icon: 'none' }); return }
  if (!registerData.name) { uni.showToast({ title: '请输入姓名', icon: 'none' }); return }
  if (!registerData.password || registerData.password.length < 6) { uni.showToast({ title: '密码至少6位', icon: 'none' }); return }
  if (registerData.password !== registerData.confirmPassword) { uni.showToast({ title: '两次密码不一致', icon: 'none' }); return }

  registerLoading.value = true
  try {
    await userStore.register({ username: registerData.username, name: registerData.name, password: registerData.password })
    uni.showToast({ title: '注册成功', icon: 'success' })
    showRegister.value = false
    formData.username = registerData.username
    formData.password = registerData.password
  } catch (error) {
    console.error('注册失败', error)
  } finally {
    registerLoading.value = false
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 100rpx 48rpx;
}

.login-header {
  text-align: center;
  margin-bottom: 80rpx;
}

.title {
  display: block;
  font-size: 48rpx;
  font-weight: 700;
  color: #fff;
  margin-bottom: 16rpx;
}

.subtitle {
  display: block;
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.8);
}

.login-form {
  background: #fff;
  border-radius: 24rpx;
  padding: 48rpx;
}

.form-item {
  display: flex;
  align-items: center;
  padding: 24rpx 0;
  border-bottom: 1rpx solid #eee;
}

.form-icon {
  font-size: 40rpx;
  margin-right: 16rpx;
}

.form-input {
  flex: 1;
  font-size: 30rpx;
}

.login-btn {
  margin-top: 48rpx;
}

.register-link {
  text-align: center;
  margin-top: 32rpx;
  font-size: 26rpx;
  color: #999;
}

.link {
  color: #007AFF;
}

.register-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: flex-end;
  z-index: 999;
}

.register-content {
  width: 100%;
  background: #fff;
  border-radius: 24rpx 24rpx 0 0;
  padding: 32rpx;
}

.popup-title {
  font-size: 36rpx;
  font-weight: 600;
  text-align: center;
  margin-bottom: 32rpx;
}

.register-content .form-label {
  width: 140rpx;
  color: #666;
}

.register-content .btn-primary {
  margin-top: 32rpx;
}

.register-content .btn-default {
  margin-top: 16rpx;
}
</style>
