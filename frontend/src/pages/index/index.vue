<template>
  <view class="index-container">
    <!-- 头部统计 -->
    <view class="header">
      <view class="header-top">
        <view class="greeting">
          <text class="greeting-text">您好，{{ userStore.userInfo?.name || '房东' }}</text>
          <text class="greeting-sub">{{ todayStr }}</text>
        </view>
      </view>
      <view class="stats-card">
        <view class="stat-item" @click="goPage('/pages/property/list')">
          <text class="stat-num">{{ stats.propertyCount }}</text>
          <text class="stat-label">房源数</text>
        </view>
        <view class="stat-item" @click="goPage('/pages/room/list')">
          <text class="stat-num">{{ stats.roomCount }}</text>
          <text class="stat-label">房间数</text>
        </view>
        <view class="stat-item" @click="goPage('/pages/tenant/list')">
          <text class="stat-num">{{ stats.tenantCount }}</text>
          <text class="stat-label">租客数</text>
        </view>
        <view class="stat-item" @click="goPage('/pages/bill/statistics')">
          <text class="stat-num highlight">{{ stats.monthlyIncome }}</text>
          <text class="stat-label">本月收入</text>
        </view>
      </view>
    </view>

    <!-- 快捷功能 -->
    <view class="section">
      <view class="section-title">快捷功能</view>
      <view class="quick-actions">
        <view class="action-item" @click="goPage('/pages/property/list')">
          <view class="action-icon icon-home">🏠</view>
          <text class="action-text">房源管理</text>
        </view>
        <view class="action-item" @click="goPage('/pages/room/list')">
          <view class="action-icon icon-room">🚪</view>
          <text class="action-text">房间管理</text>
        </view>
        <view class="action-item" @click="goPage('/pages/tenant/list')">
          <view class="action-icon icon-tenant">👥</view>
          <text class="action-text">租客管理</text>
        </view>
        <view class="action-item" @click="goPage('/pages/contract/list')">
          <view class="action-icon icon-contract">📝</view>
          <text class="action-text">合同管理</text>
        </view>
        <view class="action-item" @click="goPage('/pages/meter/list')">
          <view class="action-icon icon-meter">⚡</view>
          <text class="action-text">抄表记录</text>
        </view>
        <view class="action-item" @click="goPage('/pages/bill/list')">
          <view class="action-icon icon-bill">💰</view>
          <text class="action-text">账单管理</text>
        </view>
        <view class="action-item" @click="goPage('/pages/bill/statistics')">
          <view class="action-icon icon-stats">📊</view>
          <text class="action-text">收入统计</text>
        </view>
        <view class="action-item" @click="addMeter">
          <view class="action-icon icon-add">➕</view>
          <text class="action-text">快速抄表</text>
        </view>
      </view>
    </view>

    <!-- 待办事项 -->
    <view class="section">
      <view class="section-title">待处理</view>
      <view class="todo-list">
        <view class="todo-item" v-for="item in todos" :key="item.id">
          <view class="todo-icon">💰</view>
          <view class="todo-content">
            <text class="todo-title">{{ item.title }}</text>
            <text class="todo-desc">{{ item.desc }}</text>
          </view>
        </view>
        <view class="empty" v-if="todos.length === 0">
          <text class="empty-icon">✅</text>
          <text>暂无待处理事项</text>
        </view>
      </view>
    </view>

    <!-- 快速抄表弹窗 -->
    <view class="meter-modal" v-if="showMeter">
      <view class="meter-content">
        <view class="popup-title">快速抄表</view>
        <view class="form-item">
          <text class="form-label">选择房间</text>
          <picker :range="roomList" range-key="room_number" @change="onRoomChange">
            <view class="form-picker">{{ roomList[roomIndex]?.room_number || '请选择' }} ›</view>
          </picker>
        </view>
        <view class="form-item">
          <text class="form-label">水表读数</text>
          <input class="form-input" type="digit" v-model="meterData.water" placeholder="请输入" />
        </view>
        <view class="form-item">
          <text class="form-label">电表读数</text>
          <input class="form-input" type="digit" v-model="meterData.electricity" placeholder="请输入" />
        </view>
        <view class="form-item">
          <text class="form-label">气表读数</text>
          <input class="form-input" type="digit" v-model="meterData.gas" placeholder="请输入" />
        </view>
        <button class="btn-primary" @click="submitMeter" :disabled="submitting">{{ submitting ? '提交中...' : '提交' }}</button>
        <button class="btn-default" @click="showMeter = false">取消</button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useUserStore } from '../../stores/user'
import { authApi, billApi, meterApi, roomApi } from '../../utils/request'

const userStore = useUserStore()

const todayStr = computed(() => {
  const now = new Date()
  const weekDays = ['日', '一', '二', '三', '四', '五', '六']
  return `${now.getMonth() + 1}月${now.getDate()}日 周${weekDays[now.getDay()]}`
})

const stats = reactive({ propertyCount: 0, roomCount: 0, tenantCount: 0, monthlyIncome: '0.00' })
const todos = ref([])
const showMeter = ref(false)
const roomList = ref([])
const roomIndex = ref(0)
const submitting = ref(false)
const meterData = reactive({ water: '', electricity: '', gas: '' })

onMounted(async () => {
  await Promise.all([loadStats(), loadTodos(), loadRooms()])
})

const loadStats = async () => {
  try {
    const res = await authApi.getStatistics()
    if (res) {
      stats.propertyCount = res.property_count || 0
      stats.roomCount = res.room_count || 0
      stats.tenantCount = res.tenant_count || 0
      stats.monthlyIncome = (res.monthly_income || 0).toFixed(2)
    }
  } catch (error) {
    console.error('加载统计数据失败', error)
  }
}

const loadTodos = async () => {
  try {
    const bills = await billApi.getList({ status: 1, page: 1, page_size: 5 })
    todos.value = (bills.list || []).map(bill => ({
      id: bill.id,
      title: `${bill.room?.room_number || ''} 待收款`,
      desc: `金额: ¥${bill.amount} | ${bill.bill_month}`
    }))
  } catch (error) {
    console.error('加载待办失败', error)
  }
}

const loadRooms = async () => {
  try {
    const res = await roomApi.getList({ page: 1, page_size: 100 })
    roomList.value = res.list || []
  } catch (error) {
    console.error('加载房间失败', error)
  }
}

const goPage = (url) => {
  // tabBar 页面需要用 switchTab
  const tabBarPages = ['/pages/index/index', '/pages/property/list', '/pages/bill/list', '/pages/profile/index']
  if (tabBarPages.includes(url)) {
    uni.switchTab({ url })
  } else {
    uni.navigateTo({ url })
  }
}

const addMeter = () => { showMeter.value = true }

const onRoomChange = (e) => { roomIndex.value = e.detail.value }

const submitMeter = async () => {
  if (!roomList.value[roomIndex.value]) {
    uni.showToast({ title: '请选择房间', icon: 'none' })
    return
  }
  submitting.value = true
  try {
    const today = new Date()
    const dateStr = `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}-${String(today.getDate()).padStart(2, '0')}`
    await meterApi.create({
      room_id: roomList.value[roomIndex.value].id,
      reading_date: dateStr,
      water_reading: parseFloat(meterData.water) || 0,
      electricity_reading: parseFloat(meterData.electricity) || 0,
      gas_reading: parseFloat(meterData.gas) || 0
    })
    uni.showToast({ title: '提交成功', icon: 'success' })
    showMeter.value = false
    meterData.water = ''
    meterData.electricity = ''
    meterData.gas = ''
  } catch (error) {
    console.error('提交失败', error)
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.index-container {
  min-height: 100vh;
  background: #f8f8f8;
}

.header {
  background: linear-gradient(135deg, #007AFF 0%, #5856D6 100%);
  padding: 32rpx 24rpx 48rpx;
}

.header-top {
  margin-bottom: 32rpx;
}

.greeting-text {
  display: block;
  font-size: 36rpx;
  font-weight: 600;
  color: #fff;
}

.greeting-sub {
  display: block;
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.8);
  margin-top: 8rpx;
}

.stats-card {
  background: #fff;
  border-radius: 16rpx;
  padding: 32rpx;
  display: flex;
  justify-content: space-around;
}

.stat-item {
  text-align: center;
}

.stat-num {
  display: block;
  font-size: 44rpx;
  font-weight: 700;
  color: #333;
}

.stat-num.highlight {
  color: #007AFF;
}

.stat-label {
  display: block;
  font-size: 24rpx;
  color: #999;
  margin-top: 8rpx;
}

.section {
  padding: 24rpx;
}

.section-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 24rpx;
}

.quick-actions {
  display: flex;
  flex-wrap: wrap;
  background: #fff;
  border-radius: 16rpx;
  padding: 16rpx;
}

.action-item {
  width: 25%;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 24rpx 0;
}

.action-icon {
  width: 88rpx;
  height: 88rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 40rpx;
  margin-bottom: 12rpx;
  box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.1);
}

.icon-home { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
.icon-room { background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%); }
.icon-tenant { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }
.icon-contract { background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); }
.icon-meter { background: linear-gradient(135deg, #fa709a 0%, #fee140 100%); }
.icon-bill { background: linear-gradient(135deg, #a8edea 0%, #fed6e3 100%); }
.icon-stats { background: linear-gradient(135deg, #d299c2 0%, #fef9d7 100%); }
.icon-add { background: linear-gradient(135deg, #89f7fe 0%, #66a6ff 100%); }

.action-text {
  font-size: 24rpx;
  color: #333;
}

.todo-list {
  background: #fff;
  border-radius: 16rpx;
}

.todo-item {
  display: flex;
  align-items: center;
  padding: 24rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.todo-icon {
  font-size: 36rpx;
  margin-right: 24rpx;
}

.todo-content {
  flex: 1;
}

.todo-title {
  display: block;
  font-size: 30rpx;
  color: #333;
}

.todo-desc {
  display: block;
  font-size: 24rpx;
  color: #999;
  margin-top: 8rpx;
}

.meter-modal {
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

.meter-content {
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

.meter-content .form-item {
  display: flex;
  align-items: center;
  padding: 24rpx 0;
  border-bottom: 1rpx solid #eee;
}

.meter-content .form-label {
  width: 180rpx;
  color: #666;
}

.meter-content .form-input {
  flex: 1;
  text-align: right;
}

.form-picker {
  flex: 1;
  text-align: right;
  color: #333;
}

.meter-content .btn-primary {
  margin-top: 32rpx;
}

.meter-content .btn-default {
  margin-top: 16rpx;
}

.empty {
  padding: 60rpx;
  text-align: center;
  color: #999;
}
</style>
