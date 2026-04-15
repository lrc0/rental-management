<template>
  <view class="page-container">
    <view class="user-card">
      <view class="user-info">
        <view class="user-avatar">
          <text>{{ userStore.userInfo?.name?.charAt(0) || '房' }}</text>
        </view>
        <view class="user-detail">
          <text class="user-name">{{ userStore.userInfo?.name || '未登录' }}</text>
          <text class="user-phone">{{ userStore.userInfo?.phone || '' }}</text>
        </view>
      </view>
    </view>

    <view class="menu-list">
      <view class="menu-group">
        <view class="menu-item" @click="goToFeeRate">
          <text class="menu-icon">⚙️</text>
          <text class="menu-label">费率设置</text>
          <text class="menu-arrow">›</text>
        </view>
        <view class="menu-item" @click="changePassword">
          <text class="menu-icon">🔐</text>
          <text class="menu-label">修改密码</text>
          <text class="menu-arrow">›</text>
        </view>
      </view>

      <view class="menu-group">
        <view class="menu-item" @click="showAbout">
          <text class="menu-icon">ℹ️</text>
          <text class="menu-label">关于我们</text>
          <text class="menu-arrow">›</text>
        </view>
      </view>

      <view class="menu-group">
        <view class="menu-item logout" @click="handleLogout">
          <text class="menu-icon">🚪</text>
          <text class="menu-label">退出登录</text>
        </view>
      </view>
    </view>

    <view class="version">
      <text>租房管理系统 v1.0.0</text>
    </view>

    <!-- 费率设置弹窗 -->
    <view class="modal" v-if="showFeeRate">
      <view class="modal-content">
        <view class="popup-title">费率设置</view>
        <view class="form-item">
          <text class="form-label">水费(元/吨)</text>
          <input class="form-input" type="digit" v-model="feeRate.water_rate" placeholder="0.00" />
        </view>
        <view class="form-item">
          <text class="form-label">电费(元/度)</text>
          <input class="form-input" type="digit" v-model="feeRate.electricity_rate" placeholder="0.00" />
        </view>
        <view class="form-item">
          <text class="form-label">气费(元/立方)</text>
          <input class="form-input" type="digit" v-model="feeRate.gas_rate" placeholder="0.00" />
        </view>
        <button class="btn-primary" @click="saveFeeRate" :disabled="saving">{{ saving ? '保存中...' : '保存' }}</button>
        <button class="btn-default" @click="showFeeRate = false">取消</button>
      </view>
    </view>

    <!-- 修改密码弹窗 -->
    <view class="modal" v-if="showPassword">
      <view class="modal-content">
        <view class="popup-title">修改密码</view>
        <view class="form-item">
          <text class="form-label">原密码</text>
          <input class="form-input" password v-model="passwordForm.old_password" placeholder="请输入原密码" />
        </view>
        <view class="form-item">
          <text class="form-label">新密码</text>
          <input class="form-input" password v-model="passwordForm.new_password" placeholder="请输入新密码" />
        </view>
        <view class="form-item">
          <text class="form-label">确认密码</text>
          <input class="form-input" password v-model="passwordForm.confirm_password" placeholder="请再次输入" />
        </view>
        <button class="btn-primary" @click="submitPassword" :disabled="changingPwd">{{ changingPwd ? '修改中...' : '确定' }}</button>
        <button class="btn-default" @click="showPassword = false">取消</button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useUserStore } from '../../stores/user'
import { feeRateApi, authApi } from '../../utils/request'

const userStore = useUserStore()

const showFeeRate = ref(false)
const showPassword = ref(false)
const saving = ref(false)
const changingPwd = ref(false)

const feeRate = reactive({ water_rate: '', electricity_rate: '', gas_rate: '' })
const passwordForm = reactive({ old_password: '', new_password: '', confirm_password: '' })

onMounted(async () => { await loadFeeRate() })

const loadFeeRate = async () => {
  try {
    const res = await feeRateApi.get()
    feeRate.water_rate = String(res.water_rate || '')
    feeRate.electricity_rate = String(res.electricity_rate || '')
    feeRate.gas_rate = String(res.gas_rate || '')
  } catch (error) { console.error(error) }
}

const goToFeeRate = () => { showFeeRate.value = true }

const saveFeeRate = async () => {
  saving.value = true
  try {
    await feeRateApi.update({
      water_rate: parseFloat(feeRate.water_rate) || 0,
      electricity_rate: parseFloat(feeRate.electricity_rate) || 0,
      gas_rate: parseFloat(feeRate.gas_rate) || 0
    })
    uni.showToast({ title: '保存成功', icon: 'success' })
    showFeeRate.value = false
  } catch (error) { console.error(error) } finally { saving.value = false }
}

const changePassword = () => {
  passwordForm.old_password = ''
  passwordForm.new_password = ''
  passwordForm.confirm_password = ''
  showPassword.value = true
}

const submitPassword = async () => {
  if (!passwordForm.old_password) { uni.showToast({ title: '请输入原密码', icon: 'none' }); return }
  if (passwordForm.new_password.length < 6) { uni.showToast({ title: '新密码至少6位', icon: 'none' }); return }
  if (passwordForm.new_password !== passwordForm.confirm_password) { uni.showToast({ title: '两次密码不一致', icon: 'none' }); return }
  changingPwd.value = true
  try {
    await authApi.changePassword({ old_password: passwordForm.old_password, new_password: passwordForm.new_password })
    uni.showToast({ title: '修改成功', icon: 'success' })
    showPassword.value = false
  } catch (error) { console.error(error) } finally { changingPwd.value = false }
}

const showAbout = () => {
  uni.showModal({
    title: '关于我们',
    content: '租房管理系统\n版本：1.0.0\n\n专为个人房东打造的房源管理工具。',
    showCancel: false
  })
}

const handleLogout = () => {
  uni.showModal({
    title: '退出登录',
    content: '确定要退出登录吗？',
    success: (res) => { if (res.confirm) userStore.logout() }
  })
}
</script>

<style scoped>
.page-container {
  min-height: 100vh;
  background: #f8f8f8;
}

.user-card {
  background: linear-gradient(135deg, #007AFF 0%, #5856D6 100%);
  padding: 60rpx 32rpx;
}

.user-info { display: flex; align-items: center; }

.user-avatar {
  width: 120rpx;
  height: 120rpx;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 32rpx;
}

.user-avatar text { color: #fff; font-size: 48rpx; font-weight: 600; }

.user-detail { display: flex; flex-direction: column; }
.user-name { font-size: 40rpx; font-weight: 600; color: #fff; margin-bottom: 8rpx; }
.user-phone { font-size: 28rpx; color: rgba(255, 255, 255, 0.8); }

.menu-list { padding: 24rpx; }

.menu-group {
  background: #fff;
  border-radius: 16rpx;
  margin-bottom: 24rpx;
  overflow: hidden;
}

.menu-item {
  display: flex;
  align-items: center;
  padding: 32rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.menu-item:last-child { border-bottom: none; }
.menu-icon { font-size: 40rpx; margin-right: 24rpx; }
.menu-label { flex: 1; font-size: 32rpx; color: #333; }
.menu-arrow { font-size: 32rpx; color: #ccc; }

.menu-item.logout { justify-content: center; }
.menu-item.logout .menu-label { color: #F44336; text-align: center; }

.version { text-align: center; padding: 60rpx; color: #999; font-size: 24rpx; }

.modal {
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

.modal-content { width: 100%; background: #fff; border-radius: 24rpx 24rpx 0 0; padding: 32rpx; }
.popup-title { font-size: 36rpx; font-weight: 600; text-align: center; margin-bottom: 32rpx; }
.modal-content .form-item { display: flex; align-items: center; padding: 24rpx 0; border-bottom: 1rpx solid #eee; }
.modal-content .form-label { width: 200rpx; color: #666; }
.modal-content .form-input { flex: 1; font-size: 28rpx; text-align: right; }
.modal-content .btn-primary { margin-top: 32rpx; }
.modal-content .btn-default { margin-top: 16rpx; }
</style>
