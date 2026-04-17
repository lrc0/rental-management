<template>
  <view class="page-container">
    <view class="year-selector">
      <text class="year-btn" @click="prevYear">‹</text>
      <text class="year-text">{{ currentYear }}年</text>
      <text class="year-btn" @click="nextYear">›</text>
    </view>

    <view class="summary-card">
      <view class="summary-title">年度收入汇总</view>
      <view class="summary-amount">¥{{ yearlyTotal }}</view>
      <view class="summary-detail">
        <view class="detail-item">
          <text class="detail-label">已收款</text>
          <text class="detail-value success">¥{{ yearlyPaid }}</text>
        </view>
        <view class="detail-item">
          <text class="detail-label">待收款</text>
          <text class="detail-value warning">¥{{ yearlyPending }}</text>
        </view>
      </view>
    </view>

    <view class="detail-card">
      <view class="card-title">月度明细</view>
      <view class="month-list">
        <template v-for="item in monthlyData" :key="item.month">
          <view class="month-item" v-if="item && item.bill_count > 0">
            <view class="month-info">
              <text class="month-label">{{ item.month }}</text>
              <text class="month-count">{{ item.bill_count }}笔账单</text>
            </view>
            <view class="month-amounts">
              <text class="amount paid">¥{{ (item.paid_fee || 0).toFixed(2) }}</text>
              <text class="amount total">/ ¥{{ (item.total_fee || 0).toFixed(2) }}</text>
            </view>
          </view>
        </template>
        <view class="empty" v-if="!hasData">
          <text class="empty-icon">📊</text>
          <text>暂无统计数据</text>
        </view>
      </view>
    </view>

    <view class="quick-stats">
      <view class="stat-card" @click="goToBills">
        <view class="stat-icon" style="background: #FFF3E0;">💰</view>
        <view class="stat-info">
          <text class="stat-num">{{ pendingCount }}</text>
          <text class="stat-label">待收款账单</text>
        </view>
        <text class="stat-arrow">›</text>
      </view>
      <view class="stat-card" @click="goToBills">
        <view class="stat-icon" style="background: #E8F5E9;">✅</view>
        <view class="stat-info">
          <text class="stat-num">{{ paidCount }}</text>
          <text class="stat-label">已收款账单</text>
        </view>
        <text class="stat-arrow">›</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { billApi } from '../../utils/request'

const currentYear = ref(new Date().getFullYear())
const monthlyData = ref([])
const pendingCount = ref(0)
const paidCount = ref(0)

const yearlyTotal = computed(() => {
  const data = monthlyData.value || []
  return data.reduce((sum, item) => sum + (item?.total_fee || 0), 0).toFixed(2)
})

const yearlyPaid = computed(() => {
  const data = monthlyData.value || []
  return data.reduce((sum, item) => sum + (item?.paid_fee || 0), 0).toFixed(2)
})

const yearlyPending = computed(() => {
  return (parseFloat(yearlyTotal.value) - parseFloat(yearlyPaid.value)).toFixed(2)
})

const hasData = computed(() => {
  const data = monthlyData.value || []
  return data.some(item => item && item.bill_count > 0)
})

onMounted(() => loadData())
watch(currentYear, () => loadData())

const loadData = async () => {
  // 加载月度统计
  try {
    const monthRes = await billApi.getMonthlyStatistics({ year: currentYear.value })
    monthlyData.value = monthRes || []
  } catch (error) {
    console.error('加载月度统计失败', error)
    monthlyData.value = []
  }

  // 加载账单统计
  try {
    const statsRes = await billApi.getStatistics({})
    if (statsRes && typeof statsRes === 'object') {
      pendingCount.value = (statsRes.bill_count || 0) - (statsRes.paid_count || 0)
      paidCount.value = statsRes.paid_count || 0
    } else {
      pendingCount.value = 0
      paidCount.value = 0
    }
  } catch (error) {
    console.error('加载账单统计失败', error)
    pendingCount.value = 0
    paidCount.value = 0
  }
}

const prevYear = () => { currentYear.value-- }
const nextYear = () => { if (currentYear.value < new Date().getFullYear()) currentYear.value++ }

const goToBills = () => {
  uni.switchTab({ url: '/pages/bill/list' })
}
</script>

<style scoped>
.page-container {
  min-height: 100vh;
  background: #f8f8f8;
  padding: 24rpx;
  padding-bottom: 40rpx;
}

.year-selector {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 24rpx;
  background: #fff;
  border-radius: 16rpx;
  margin-bottom: 24rpx;
}

.year-btn {
  width: 64rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 40rpx;
  color: #007AFF;
  background: #f5f5f5;
  border-radius: 50%;
}
.year-text { font-size: 36rpx; font-weight: 600; color: #333; margin: 0 32rpx; }

.summary-card {
  background: linear-gradient(135deg, #007AFF 0%, #5856D6 100%);
  border-radius: 16rpx;
  padding: 32rpx;
  margin-bottom: 24rpx;
  color: #fff;
}

.summary-title { font-size: 28rpx; opacity: 0.8; margin-bottom: 16rpx; }
.summary-amount { font-size: 56rpx; font-weight: 700; margin-bottom: 24rpx; }

.summary-detail { display: flex; gap: 48rpx; }
.detail-item { display: flex; flex-direction: column; }
.detail-label { font-size: 24rpx; opacity: 0.8; }
.detail-value { font-size: 32rpx; font-weight: 600; margin-top: 8rpx; }
.detail-value.success { color: #81C784; }
.detail-value.warning { color: #FFD54F; }

.detail-card {
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 24rpx;
}
.card-title { font-size: 32rpx; font-weight: 600; color: #333; margin-bottom: 24rpx; }

.month-list { max-height: 500rpx; overflow-y: auto; }
.month-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20rpx 0;
  border-bottom: 1rpx solid #f0f0f0;
}
.month-item:last-child { border-bottom: none; }
.month-info { display: flex; flex-direction: column; }
.month-label { font-size: 28rpx; color: #333; font-weight: 500; }
.month-count { font-size: 24rpx; color: #999; margin-top: 4rpx; }
.month-amounts { text-align: right; }
.amount { font-size: 28rpx; }
.amount.paid { color: #4CAF50; font-weight: 600; }
.amount.total { color: #999; }

.quick-stats { display: flex; gap: 24rpx; }
.stat-card {
  flex: 1;
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  display: flex;
  align-items: center;
}
.stat-icon {
  width: 80rpx;
  height: 80rpx;
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 40rpx;
  margin-right: 20rpx;
}
.stat-info { flex: 1; display: flex; flex-direction: column; }
.stat-num { font-size: 40rpx; font-weight: 700; color: #333; }
.stat-label { font-size: 24rpx; color: #999; margin-top: 4rpx; }
.stat-arrow { font-size: 36rpx; color: #ccc; }

.empty { padding: 60rpx; text-align: center; color: #999; }
.empty-icon { display: block; font-size: 80rpx; margin-bottom: 16rpx; }
</style>
