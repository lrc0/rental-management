<template>
  <view class="page-container">
    <view class="stats-card">
      <view class="stat-item">
        <text class="stat-label">待收金额</text>
        <text class="stat-value warning">¥{{ stats.pending }}</text>
      </view>
      <view class="stat-divider"></view>
      <view class="stat-item">
        <text class="stat-label">已收金额</text>
        <text class="stat-value success">¥{{ stats.paid }}</text>
      </view>
    </view>

    <view class="bill-list">
      <view class="bill-card" v-for="bill in list" :key="bill.id">
        <view class="bill-header">
          <view class="bill-room">{{ bill.room?.room_number || '-' }}</view>
          <view class="bill-status" :class="getStatusClass(bill.status)">{{ getStatusText(bill.status) }}</view>
        </view>
        <view class="bill-body">
          <view class="bill-tenant" v-if="bill.tenant">租客：{{ bill.tenant.name }}</view>
          <view class="bill-month">{{ bill.bill_month }}</view>

          <!-- 费用明细 -->
          <view class="fee-detail">
            <view class="fee-row" v-if="bill.rent_fee > 0">
              <text class="fee-label">租金</text>
              <text class="fee-value">¥{{ bill.rent_fee }}</text>
            </view>
            <view class="fee-row" v-if="bill.water_fee > 0">
              <text class="fee-label">水费</text>
              <text class="fee-value">¥{{ bill.water_fee }}</text>
            </view>
            <view class="fee-row" v-if="bill.electricity_fee > 0">
              <text class="fee-label">电费</text>
              <text class="fee-value">¥{{ bill.electricity_fee }}</text>
            </view>
            <view class="fee-row" v-if="bill.gas_fee > 0">
              <text class="fee-label">气费</text>
              <text class="fee-value">¥{{ bill.gas_fee }}</text>
            </view>
          </view>

          <view class="bill-total">
            <text>合计</text>
            <text class="total-amount">¥{{ bill.amount }}</text>
          </view>
        </view>
        <view class="bill-actions">
          <button class="btn-pay" v-if="bill.status === 1" @click="payBill(bill)">收款</button>
          <button class="btn-delete" @click="deleteBill(bill)">删除</button>
        </view>
      </view>
      <view class="empty" v-if="list.length === 0 && !loading">
        <text class="empty-icon">💰</text>
        <text>暂无账单</text>
      </view>
    </view>

    <view class="add-btn" @click="addBill"><text>+</text></view>

    <!-- 创建账单弹窗 -->
    <view class="modal" v-if="showAdd">
      <view class="modal-content">
        <view class="popup-title">创建账单</view>
        <view class="form-item">
          <text class="form-label">选择房间</text>
          <picker :range="roomList" range-key="room_number" @change="onRoomChange">
            <view class="form-picker">{{ selectedRoom?.room_number || '请选择' }} ›</view>
          </picker>
        </view>
        <view class="form-item">
          <text class="form-label">账单月份</text>
          <picker mode="date" fields="month" :value="formData.bill_month" @change="(e) => formData.bill_month = e.detail.value">
            <view class="form-picker">{{ formData.bill_month || '请选择' }} ›</view>
          </picker>
        </view>

        <!-- 费用预览区域 -->
        <view class="fee-preview" v-if="previewData && selectedRoom">
          <view class="preview-title">费用明细（自动计算）</view>
          <view class="preview-item">
            <text class="preview-label">租金</text>
            <input class="preview-input" type="digit" v-model="formData.rent_fee" placeholder="0.00" @input="calculateTotal" />
          </view>
          <view class="preview-item" v-if="previewData.has_reading">
            <text class="preview-label">水费 <text class="usage" v-if="previewData.water_usage > 0">({{ previewData.water_usage }}吨 × {{ previewData.water_rate }}元)</text></text>
            <input class="preview-input" type="digit" v-model="formData.water_fee" placeholder="0.00" @input="calculateTotal" />
          </view>
          <view class="preview-item" v-else>
            <text class="preview-label">水费</text>
            <input class="preview-input" type="digit" v-model="formData.water_fee" placeholder="0.00" @input="calculateTotal" />
          </view>
          <view class="preview-item" v-if="previewData.has_reading">
            <text class="preview-label">电费 <text class="usage" v-if="previewData.electricity_usage > 0">({{ previewData.electricity_usage }}度 × {{ previewData.electricity_rate }}元)</text></text>
            <input class="preview-input" type="digit" v-model="formData.electricity_fee" placeholder="0.00" @input="calculateTotal" />
          </view>
          <view class="preview-item" v-else>
            <text class="preview-label">电费</text>
            <input class="preview-input" type="digit" v-model="formData.electricity_fee" placeholder="0.00" @input="calculateTotal" />
          </view>
          <view class="preview-item" v-if="previewData.has_reading">
            <text class="preview-label">气费 <text class="usage" v-if="previewData.gas_usage > 0">({{ previewData.gas_usage }}立方 × {{ previewData.gas_rate }}元)</text></text>
            <input class="preview-input" type="digit" v-model="formData.gas_fee" placeholder="0.00" @input="calculateTotal" />
          </view>
          <view class="preview-item" v-else>
            <text class="preview-label">气费</text>
            <input class="preview-input" type="digit" v-model="formData.gas_fee" placeholder="0.00" @input="calculateTotal" />
          </view>
          <view class="preview-total">
            <text class="preview-label">合计</text>
            <text class="preview-amount">¥{{ totalAmount }}</text>
          </view>
        </view>

        <!-- 无预览数据时的手动输入 -->
        <view class="fee-manual" v-if="!previewData && selectedRoom">
          <view class="form-item">
            <text class="form-label">租金</text>
            <input class="form-input" type="digit" v-model="formData.rent_fee" placeholder="0.00" />
          </view>
          <view class="form-item">
            <text class="form-label">水费</text>
            <input class="form-input" type="digit" v-model="formData.water_fee" placeholder="0.00" />
          </view>
          <view class="form-item">
            <text class="form-label">电费</text>
            <input class="form-input" type="digit" v-model="formData.electricity_fee" placeholder="0.00" />
          </view>
          <view class="form-item">
            <text class="form-label">气费</text>
            <input class="form-input" type="digit" v-model="formData.gas_fee" placeholder="0.00" />
          </view>
        </view>

        <button class="btn-primary" @click="submitForm" :disabled="submitting">{{ submitting ? '创建中...' : '创建' }}</button>
        <button class="btn-default" @click="showAdd = false">取消</button>
      </view>
    </view>

    <!-- 收款弹窗 -->
    <view class="modal" v-if="showPay">
      <view class="modal-content">
        <view class="popup-title">确认收款</view>
        <view class="pay-info">
          <view class="pay-room">{{ payingBill?.room?.room_number }}</view>
          <view class="pay-amount">¥{{ payingBill?.amount }}</view>
        </view>
        <view class="form-item">
          <text class="form-label">收款金额</text>
          <input class="form-input" type="digit" v-model="payForm.amount" placeholder="请输入" />
        </view>
        <button class="btn-primary" @click="confirmPay" :disabled="paying">{{ paying ? '处理中...' : '确认收款' }}</button>
        <button class="btn-default" @click="showPay = false">取消</button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { onPullDownRefresh } from '@dcloudio/uni-app'
import { billApi, roomApi } from '../../utils/request'

const loading = ref(false)
const submitting = ref(false)
const paying = ref(false)
const list = ref([])
const showAdd = ref(false)
const showPay = ref(false)
const roomList = ref([])
const selectedRoom = ref(null)
const payingBill = ref(null)
const previewData = ref(null)

const stats = reactive({ pending: '0.00', paid: '0.00' })
const formData = reactive({ bill_month: '', rent_fee: '', water_fee: '', electricity_fee: '', gas_fee: '' })
const payForm = reactive({ amount: '' })

const totalAmount = computed(() => {
  const rent = parseFloat(formData.rent_fee) || 0
  const water = parseFloat(formData.water_fee) || 0
  const electricity = parseFloat(formData.electricity_fee) || 0
  const gas = parseFloat(formData.gas_fee) || 0
  return (rent + water + electricity + gas).toFixed(2)
})

const getStatusText = (status) => ({ 1: '待支付', 2: '已支付', 3: '已逾期' }[status] || '未知')
const getStatusClass = (status) => ({ 1: 'warning', 2: 'success', 3: 'danger' }[status] || '')

onMounted(async () => {
  await loadList()
  await loadStats()
  await loadRooms()
})

// 监听房间和月份变化，自动计算费用
watch([selectedRoom, () => formData.bill_month], async ([room, month]) => {
  if (room && month) {
    await loadPreview()
  }
})

const loadList = async () => {
  loading.value = true
  try {
    const res = await billApi.getList({ page: 1, page_size: 100 })
    list.value = res.list || []
  } catch (error) { console.error(error) } finally { loading.value = false; uni.stopPullDownRefresh() }
}

const loadStats = async () => {
  try {
    const res = await billApi.getStatistics({})
    stats.pending = res.pending_amount?.toFixed(2) || '0.00'
    stats.paid = res.paid_amount?.toFixed(2) || '0.00'
  } catch (error) { console.error(error) }
}

const loadRooms = async () => {
  try {
    const res = await roomApi.getList({ page: 1, page_size: 100 })
    roomList.value = res.list || []
  } catch (error) { console.error(error) }
}

const loadPreview = async () => {
  if (!selectedRoom.value || !formData.bill_month) return
  try {
    const res = await billApi.preview({ room_id: selectedRoom.value.id, bill_month: formData.bill_month })
    previewData.value = res
    // 自动填充计算结果
    formData.rent_fee = res.rent_fee?.toFixed(2) || ''
    formData.water_fee = res.water_fee?.toFixed(2) || ''
    formData.electricity_fee = res.electricity_fee?.toFixed(2) || ''
    formData.gas_fee = res.gas_fee?.toFixed(2) || ''
  } catch (error) {
    console.error('加载费用预览失败', error)
    previewData.value = null
  }
}

const addBill = () => {
  selectedRoom.value = null
  previewData.value = null
  formData.bill_month = ''
  formData.rent_fee = ''
  formData.water_fee = ''
  formData.electricity_fee = ''
  formData.gas_fee = ''
  showAdd.value = true
}

const onRoomChange = (e) => {
  selectedRoom.value = roomList.value[e.detail.value]
}

const calculateTotal = () => {
  // 手动输入时重新计算总价（computed会自动处理）
}

const submitForm = async () => {
  if (!selectedRoom.value) { uni.showToast({ title: '请选择房间', icon: 'none' }); return }
  if (!formData.bill_month) { uni.showToast({ title: '请选择账单月份', icon: 'none' }); return }
  submitting.value = true
  try {
    await billApi.create({
      room_id: selectedRoom.value.id,
      bill_month: formData.bill_month,
      rent_fee: parseFloat(formData.rent_fee) || 0,
      water_fee: parseFloat(formData.water_fee) || 0,
      electricity_fee: parseFloat(formData.electricity_fee) || 0,
      gas_fee: parseFloat(formData.gas_fee) || 0,
      auto_calculate: false  // 已经在前端预览计算过了
    })
    uni.showToast({ title: '创建成功', icon: 'success' })
    showAdd.value = false
    loadList()
    loadStats()
  } catch (error) { console.error(error) } finally { submitting.value = false }
}

const deleteBill = (bill) => {
  uni.showModal({
    title: '确认删除',
    content: `确定要删除 ${bill.bill_month} 的账单吗？`,
    success: async (res) => {
      if (res.confirm) {
        try {
          await billApi.delete(bill.id)
          uni.showToast({ title: '删除成功', icon: 'success' })
          loadList()
          loadStats()
        } catch (error) { uni.showToast({ title: error.message || '删除失败', icon: 'none' }) }
      }
    }
  })
}

const payBill = (bill) => {
  payingBill.value = bill
  payForm.amount = String(bill.amount)
  showPay.value = true
}

const confirmPay = async () => {
  if (!payForm.amount) { uni.showToast({ title: '请输入收款金额', icon: 'none' }); return }
  paying.value = true
  try {
    await billApi.pay(payingBill.value.id, { amount: parseFloat(payForm.amount), payment_method: 1 })
    uni.showToast({ title: '收款成功', icon: 'success' })
    showPay.value = false
    loadList()
    loadStats()
  } catch (error) { console.error(error) } finally { paying.value = false }
}

onPullDownRefresh(() => { loadList(); loadStats() })
</script>

<style scoped>
.page-container {
  min-height: 100vh;
  background: #f8f8f8;
  padding-bottom: 120rpx;
}

.stats-card {
  display: flex;
  margin: 24rpx;
  background: #fff;
  border-radius: 16rpx;
  padding: 32rpx;
}

.stat-item { flex: 1; text-align: center; }
.stat-label { display: block; font-size: 26rpx; color: #999; margin-bottom: 8rpx; }
.stat-value { font-size: 40rpx; font-weight: 700; }
.stat-value.warning { color: #FF9800; }
.stat-value.success { color: #4CAF50; }
.stat-divider { width: 1rpx; background: #eee; }

.bill-list { padding: 0 24rpx; }

.bill-card {
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 24rpx;
}

.bill-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16rpx;
}

.bill-room { font-size: 32rpx; font-weight: 600; color: #333; }

.bill-status { padding: 8rpx 20rpx; border-radius: 8rpx; font-size: 24rpx; }
.bill-status.warning { background: #FFF3E0; color: #FF9800; }
.bill-status.success { background: #E8F5E9; color: #4CAF50; }
.bill-status.danger { background: #FFEBEE; color: #F44336; }

.bill-tenant { font-size: 26rpx; color: #666; margin-bottom: 8rpx; }
.bill-month { font-size: 24rpx; color: #999; margin-bottom: 16rpx; }

.fee-detail {
  background: #f8f8f8;
  border-radius: 12rpx;
  padding: 16rpx;
  margin-bottom: 12rpx;
}

.fee-row {
  display: flex;
  justify-content: space-between;
  padding: 8rpx 0;
  font-size: 26rpx;
}

.fee-label { color: #666; }
.fee-value { color: #333; }

.bill-total {
  display: flex;
  justify-content: space-between;
  padding-top: 16rpx;
  font-size: 28rpx;
}

.total-amount { color: #FF6B6B; font-weight: 700; font-size: 36rpx; }

.bill-actions { margin-top: 20rpx; display: flex; gap: 20rpx; }
.btn-pay { flex: 1; background: linear-gradient(135deg, #4CAF50 0%, #45a049 100%); color: #fff; border-radius: 8rpx; padding: 20rpx; font-size: 28rpx; }
.btn-delete { flex: 1; background: #f5f5f5; color: #FF3B30; border-radius: 8rpx; padding: 20rpx; font-size: 28rpx; }

.add-btn {
  position: fixed;
  right: 32rpx;
  bottom: 120rpx;
  width: 100rpx;
  height: 100rpx;
  background: linear-gradient(135deg, #007AFF 0%, #5856D6 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8rpx 24rpx rgba(0, 122, 255, 0.4);
}
.add-btn text { color: #fff; font-size: 56rpx; }

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

.modal-content { width: 100%; background: #fff; border-radius: 24rpx 24rpx 0 0; padding: 32rpx; max-height: 80vh; overflow-y: auto; }
.popup-title { font-size: 36rpx; font-weight: 600; text-align: center; margin-bottom: 32rpx; }
.modal-content .form-item { display: flex; align-items: center; padding: 24rpx 0; border-bottom: 1rpx solid #eee; }
.modal-content .form-label { width: 160rpx; color: #666; }
.modal-content .form-input { flex: 1; text-align: right; font-size: 28rpx; }
.form-picker { flex: 1; text-align: right; color: #333; }
.modal-content .btn-primary { margin-top: 32rpx; }
.modal-content .btn-default { margin-top: 16rpx; }

.fee-preview {
  background: #f8f8f8;
  border-radius: 12rpx;
  padding: 24rpx;
  margin: 24rpx 0;
}

.preview-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 20rpx;
}

.preview-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16rpx 0;
  border-bottom: 1rpx solid #eee;
}

.preview-label {
  font-size: 26rpx;
  color: #666;
}

.usage {
  font-size: 22rpx;
  color: #999;
}

.preview-input {
  width: 200rpx;
  text-align: right;
  font-size: 26rpx;
  background: #fff;
  border-radius: 8rpx;
  padding: 12rpx;
}

.preview-total {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 20rpx;
  margin-top: 16rpx;
}

.preview-amount {
  font-size: 40rpx;
  font-weight: 700;
  color: #FF6B6B;
}

.fee-manual {
  margin: 24rpx 0;
}

.pay-info { text-align: center; padding: 32rpx 0; border-bottom: 1rpx solid #eee; margin-bottom: 24rpx; }
.pay-room { font-size: 28rpx; color: #666; margin-bottom: 8rpx; }
.pay-amount { font-size: 56rpx; font-weight: 700; color: #FF6B6B; }

.empty { padding: 100rpx; text-align: center; color: #999; }
.empty-icon { display: block; font-size: 100rpx; margin-bottom: 16rpx; }
</style>
