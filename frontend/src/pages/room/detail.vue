<template>
  <view class="page-container">
    <!-- 房间基本信息 -->
    <view class="card">
      <view class="card-header">
        <view class="room-title">
          <text class="room-number">{{ room.room_number }}</text>
          <text class="room-property">{{ room.property?.name || '' }}</text>
        </view>
        <text class="status-tag" :class="getStatusClass(room.status)">{{ getStatusText(room.status) }}</text>
      </view>
      <view class="info-grid">
        <view class="info-item">
          <text class="info-label">{{ getRentTypeLabel(room.rent_type) }}</text>
          <text class="info-value price">¥{{ room.rent_amount || room.monthly_rent }}</text>
        </view>
        <view class="info-item">
          <text class="info-label">楼层</text>
          <text class="info-value">{{ room.floor || '-' }}层</text>
        </view>
      </view>
    </view>

    <!-- 租客信息 -->
    <view class="section" v-if="tenant.name">
      <view class="section-header">
        <text class="section-title">当前租客</text>
        <text class="section-action" @click="goTenant">查看详情 ›</text>
      </view>
      <view class="tenant-card">
        <view class="tenant-avatar">
          <text class="avatar-text">{{ tenant.name?.charAt(0) || '租' }}</text>
        </view>
        <view class="tenant-info">
          <text class="tenant-name">{{ tenant.name }}</text>
          <text class="tenant-phone">📱 {{ tenant.phone || '未填写' }}</text>
        </view>
      </view>
    </view>

    <!-- 最近抄表 -->
    <view class="section">
      <view class="section-header">
        <text class="section-title">抄表记录</text>
        <text class="section-action" @click="addMeter">+ 抄表</text>
      </view>
      <view class="meter-list" v-if="meters.length > 0">
        <view class="meter-item" v-for="item in meters" :key="item.id">
          <view class="meter-date">{{ formatDate(item.reading_date) }}</view>
          <view class="meter-values">
            <view class="meter-val" v-if="item.water_reading">
              <text class="val-icon">💧</text>
              <text class="val-num">{{ item.water_reading }}</text>
            </view>
            <view class="meter-val" v-if="item.electricity_reading">
              <text class="val-icon">⚡</text>
              <text class="val-num">{{ item.electricity_reading }}</text>
            </view>
            <view class="meter-val" v-if="item.gas_reading">
              <text class="val-icon">🔥</text>
              <text class="val-num">{{ item.gas_reading }}</text>
            </view>
          </view>
        </view>
      </view>
      <view class="empty-tip" v-else>暂无抄表记录</view>
    </view>

    <!-- 账单列表 -->
    <view class="section">
      <view class="section-header">
        <text class="section-title">账单记录</text>
        <text class="section-action" @click="addBill">+ 创建账单</text>
      </view>
      <view class="bill-list" v-if="bills.length > 0">
        <view class="bill-item" v-for="bill in bills" :key="bill.id" @click="payBill(bill)">
          <view class="bill-left">
            <text class="bill-month">{{ bill.bill_month }}</text>
            <text class="bill-status" :class="bill.status === 2 ? 'paid' : 'unpaid'">
              {{ bill.status === 2 ? '已支付' : '待支付' }}
            </text>
          </view>
          <view class="bill-right">
            <text class="bill-amount" :class="{ unpaid: bill.status === 1 }">¥{{ bill.amount }}</text>
            <text class="bill-arrow">›</text>
          </view>
        </view>
      </view>
      <view class="empty-tip" v-else>暂无账单记录</view>
    </view>

    <!-- 抄表弹窗 -->
    <view class="modal" v-if="showMeter">
      <view class="modal-content">
        <view class="popup-title">录入抄表</view>
        <view class="form-item">
          <text class="form-label">水表读数</text>
          <input class="form-input" type="digit" v-model="meterForm.water" placeholder="请输入" />
        </view>
        <view class="form-item">
          <text class="form-label">电表读数</text>
          <input class="form-input" type="digit" v-model="meterForm.electricity" placeholder="请输入" />
        </view>
        <view class="form-item">
          <text class="form-label">气表读数</text>
          <input class="form-input" type="digit" v-model="meterForm.gas" placeholder="请输入" />
        </view>
        <button class="btn-primary" @click="submitMeter" :disabled="submitting">{{ submitting ? '提交中...' : '提交' }}</button>
        <button class="btn-default" @click="showMeter = false">取消</button>
      </view>
    </view>

    <!-- 创建账单弹窗 -->
    <view class="modal" v-if="showBill">
      <view class="modal-content">
        <view class="popup-title">创建账单</view>
        <view class="form-item">
          <text class="form-label">账单月份</text>
          <picker mode="date" fields="month" :value="billForm.bill_month" @change="(e) => billForm.bill_month = e.detail.value">
            <view class="form-picker">{{ billForm.bill_month || '请选择' }} ›</view>
          </picker>
        </view>
        <view class="form-item">
          <text class="form-label">租金</text>
          <input class="form-input" type="digit" v-model="billForm.rent_fee" placeholder="0.00" />
        </view>
        <view class="form-item">
          <text class="form-label">水费</text>
          <input class="form-input" type="digit" v-model="billForm.water_fee" placeholder="0.00" />
        </view>
        <view class="form-item">
          <text class="form-label">电费</text>
          <input class="form-input" type="digit" v-model="billForm.electricity_fee" placeholder="0.00" />
        </view>
        <button class="btn-primary" @click="submitBill" :disabled="submitting">{{ submitting ? '创建中...' : '创建' }}</button>
        <button class="btn-default" @click="showBill = false">取消</button>
      </view>
    </view>

    <!-- 收款弹窗 -->
    <view class="modal" v-if="showPay">
      <view class="modal-content">
        <view class="popup-title">确认收款</view>
        <view class="pay-info">
          <text class="pay-month">{{ payingBill?.bill_month }}</text>
          <text class="pay-amount">¥{{ payingBill?.amount }}</text>
        </view>
        <view class="form-item">
          <text class="form-label">收款金额</text>
          <input class="form-input" type="digit" v-model="payAmount" placeholder="请输入" />
        </view>
        <button class="btn-primary" @click="confirmPay" :disabled="submitting">{{ submitting ? '处理中...' : '确认收款' }}</button>
        <button class="btn-default" @click="showPay = false">取消</button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { roomApi, meterApi, billApi } from '../../utils/request'

const roomId = ref(0)
const room = ref({})
const tenant = ref({})
const meters = ref([])
const bills = ref([])

const showMeter = ref(false)
const showBill = ref(false)
const showPay = ref(false)
const submitting = ref(false)
const payingBill = ref(null)
const payAmount = ref('')

const meterForm = reactive({ water: '', electricity: '', gas: '' })
const billForm = reactive({ bill_month: '', rent_fee: '', water_fee: '', electricity_fee: '' })

const getStatusText = (status) => ({ 1: '空置', 2: '已租', 3: '维修中' }[status] || '未知')
const getStatusClass = (status) => ({ 1: 'info', 2: 'success', 3: 'warning' }[status] || '')
const getRentTypeLabel = (type) => ({ 1: '月租金', 2: '季租金', 3: '年租金' }[type] || '月租金')
const formatDate = (date) => date ? date.substring(0, 10) : '-'

onMounted(async () => {
  const pages = getCurrentPages()
  roomId.value = pages[pages.length - 1].options?.id
  if (roomId.value) {
    await loadData()
  }
})

const loadData = async () => {
  await Promise.all([loadRoom(), loadMeters(), loadBills()])
}

const loadRoom = async () => {
  try {
    const res = await roomApi.getDetail(roomId.value)
    room.value = res || {}
    tenant.value = res.tenant || {}
    // 从房间列表接口获取租客信息
    const listRes = await roomApi.getList({ page: 1, page_size: 1 })
    const roomWithTenant = (listRes.list || []).find(r => r.id == roomId.value)
    if (roomWithTenant) {
      tenant.value = {
        id: roomWithTenant.tenant_id,
        name: roomWithTenant.tenant_name,
        phone: roomWithTenant.tenant_phone
      }
    }
  } catch (error) { console.error(error) }
}

const loadMeters = async () => {
  try {
    const res = await meterApi.getList({ room_id: parseInt(roomId.value) || 0, page: 1, page_size: 10 })
    meters.value = res.list || []
  } catch (error) { console.error(error) }
}

const loadBills = async () => {
  try {
    const res = await billApi.getList({ room_id: parseInt(roomId.value) || 0, page: 1, page_size: 10 })
    bills.value = res.list || []
  } catch (error) { console.error(error) }
}

const goTenant = () => {
  if (tenant.value.id) {
    uni.navigateTo({ url: `/pages/tenant/detail?id=${tenant.value.id}` })
  }
}

const addMeter = () => {
  meterForm.water = ''
  meterForm.electricity = ''
  meterForm.gas = ''
  showMeter.value = true
}

const submitMeter = async () => {
  submitting.value = true
  try {
    const today = new Date()
    const dateStr = `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}-${String(today.getDate()).padStart(2, '0')}`
    await meterApi.create({
      room_id: parseInt(roomId.value) || 0,
      reading_date: dateStr,
      water_reading: parseFloat(meterForm.water) || 0,
      electricity_reading: parseFloat(meterForm.electricity) || 0,
      gas_reading: parseFloat(meterForm.gas) || 0
    })
    uni.showToast({ title: '提交成功', icon: 'success' })
    showMeter.value = false
    loadMeters()
  } catch (error) { console.error(error) } finally { submitting.value = false }
}

const addBill = () => {
  const today = new Date()
  billForm.bill_month = `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}`
  billForm.rent_fee = String(room.value.monthly_rent || '')
  billForm.water_fee = ''
  billForm.electricity_fee = ''
  showBill.value = true
}

const submitBill = async () => {
  if (!billForm.bill_month) {
    uni.showToast({ title: '请选择账单月份', icon: 'none' })
    return
  }
  submitting.value = true
  try {
    await billApi.create({
      room_id: parseInt(roomId.value) || 0,
      bill_month: billForm.bill_month,
      rent_fee: parseFloat(billForm.rent_fee) || 0,
      water_fee: parseFloat(billForm.water_fee) || 0,
      electricity_fee: parseFloat(billForm.electricity_fee) || 0
    })
    uni.showToast({ title: '创建成功', icon: 'success' })
    showBill.value = false
    loadBills()
  } catch (error) { console.error(error) } finally { submitting.value = false }
}

const payBill = (bill) => {
  if (bill.status === 2) return // 已支付
  payingBill.value = bill
  payAmount.value = String(bill.amount)
  showPay.value = true
}

const confirmPay = async () => {
  if (!payAmount.value) {
    uni.showToast({ title: '请输入收款金额', icon: 'none' })
    return
  }
  submitting.value = true
  try {
    await billApi.pay(payingBill.value.id, { amount: parseFloat(payAmount.value), payment_method: 1 })
    uni.showToast({ title: '收款成功', icon: 'success' })
    showPay.value = false
    loadBills()
  } catch (error) { console.error(error) } finally { submitting.value = false }
}
</script>

<style scoped>
.page-container {
  min-height: 100vh;
  background: #f8f8f8;
  padding: 24rpx;
  padding-bottom: 40rpx;
}

.card {
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 24rpx;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24rpx;
}

.room-title {
  display: flex;
  flex-direction: column;
}

.room-number {
  font-size: 40rpx;
  font-weight: 700;
  color: #333;
}

.room-property {
  font-size: 26rpx;
  color: #999;
  margin-top: 8rpx;
}

.status-tag {
  padding: 8rpx 20rpx;
  border-radius: 8rpx;
  font-size: 24rpx;
}
.status-tag.success { background: #E8F5E9; color: #4CAF50; }
.status-tag.info { background: #E3F2FD; color: #2196F3; }
.status-tag.warning { background: #FFF3E0; color: #FF9800; }

.info-grid {
  display: flex;
  gap: 40rpx;
}

.info-item {
  display: flex;
  flex-direction: column;
}

.info-label {
  font-size: 24rpx;
  color: #999;
  margin-bottom: 8rpx;
}

.info-value {
  font-size: 32rpx;
  color: #333;
  font-weight: 500;
}

.info-value.price {
  color: #FF6B6B;
}

.section {
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 24rpx;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20rpx;
}

.section-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
}

.section-action {
  font-size: 26rpx;
  color: #007AFF;
}

.tenant-card {
  display: flex;
  align-items: center;
  background: #f8f8f8;
  border-radius: 12rpx;
  padding: 20rpx;
}

.tenant-avatar {
  width: 72rpx;
  height: 72rpx;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 20rpx;
}

.avatar-text {
  color: #fff;
  font-size: 32rpx;
  font-weight: 600;
}

.tenant-info {
  flex: 1;
}

.tenant-name {
  display: block;
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 4rpx;
}

.tenant-phone {
  font-size: 24rpx;
  color: #999;
}

.meter-list, .bill-list {
  background: #f8f8f8;
  border-radius: 12rpx;
}

.meter-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20rpx;
  border-bottom: 1rpx solid #eee;
}

.meter-item:last-child {
  border-bottom: none;
}

.meter-date {
  font-size: 26rpx;
  color: #666;
}

.meter-values {
  display: flex;
  gap: 24rpx;
}

.meter-val {
  display: flex;
  align-items: center;
  gap: 8rpx;
}

.val-icon {
  font-size: 24rpx;
}

.val-num {
  font-size: 26rpx;
  font-weight: 600;
  color: #333;
}

.bill-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20rpx;
  border-bottom: 1rpx solid #eee;
}

.bill-item:last-child {
  border-bottom: none;
}

.bill-left {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.bill-month {
  font-size: 28rpx;
  color: #333;
}

.bill-status {
  font-size: 22rpx;
  padding: 4rpx 12rpx;
  border-radius: 6rpx;
}

.bill-status.paid {
  background: #E8F5E9;
  color: #4CAF50;
}

.bill-status.unpaid {
  background: #FFF3E0;
  color: #FF9800;
}

.bill-right {
  display: flex;
  align-items: center;
}

.bill-amount {
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
}

.bill-amount.unpaid {
  color: #FF6B6B;
}

.bill-arrow {
  color: #ccc;
  font-size: 32rpx;
  margin-left: 8rpx;
}

.empty-tip {
  text-align: center;
  color: #999;
  font-size: 26rpx;
  padding: 40rpx;
}

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

.modal-content {
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

.modal-content .form-item {
  display: flex;
  align-items: center;
  padding: 24rpx 0;
  border-bottom: 1rpx solid #eee;
}

.modal-content .form-label {
  width: 180rpx;
  color: #666;
  font-size: 28rpx;
}

.modal-content .form-input {
  flex: 1;
  font-size: 28rpx;
  text-align: right;
}

.form-picker {
  flex: 1;
  text-align: right;
  color: #333;
}

.pay-info {
  text-align: center;
  padding: 24rpx 0;
  border-bottom: 1rpx solid #eee;
  margin-bottom: 24rpx;
}

.pay-month {
  display: block;
  font-size: 28rpx;
  color: #666;
  margin-bottom: 8rpx;
}

.pay-amount {
  font-size: 48rpx;
  font-weight: 700;
  color: #FF6B6B;
}

.modal-content .btn-primary {
  margin-top: 32rpx;
  background: #007AFF;
  color: #fff;
  border-radius: 12rpx;
  padding: 24rpx;
  font-size: 30rpx;
}

.modal-content .btn-default {
  margin-top: 16rpx;
  background: #f5f5f5;
  color: #666;
  border-radius: 12rpx;
  padding: 24rpx;
  font-size: 30rpx;
}
</style>
