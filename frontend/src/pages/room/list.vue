<template>
  <view class="page-container">
    <view class="room-list">
      <view class="room-card" v-for="room in list" :key="room.id" @click="goDetail(room.id)">
        <view class="room-header">
          <view class="room-number">{{ room.room_number }}</view>
          <view class="room-status" :class="getStatusClass(room.status)">{{ getStatusText(room.status) }}</view>
        </view>
        <view class="room-info">
          <view class="info-row">
            <text class="info-label">所属房源</text>
            <text class="info-value">{{ room.property_name || '-' }}</text>
          </view>
          <view class="info-row">
            <text class="info-label">{{ getRentTypeLabel(room.rent_type) }}</text>
            <text class="info-value price">¥{{ room.rent_amount || room.monthly_rent }}</text>
          </view>
          <view class="info-row tenant-row" v-if="room.tenant_name" @click.stop="goTenant(room.tenant_id)">
            <text class="info-label">当前租客</text>
            <view class="tenant-info">
              <text class="info-value tenant">{{ room.tenant_name }} {{ room.tenant_phone }}</text>
              <text class="tenant-arrow">›</text>
            </view>
          </view>
        </view>
        <view class="room-actions">
          <button class="action-btn" @click.stop="changeStatus(room)">
            {{ room.status === 1 ? '标记已租' : room.status === 2 ? '标记空置' : '标记维修完成' }}
          </button>
        </view>
      </view>
      <view class="empty" v-if="list.length === 0 && !loading">
        <text class="empty-icon">🚪</text>
        <text>暂无房间</text>
      </view>
    </view>

    <view class="add-btn" @click="addRoom"><text>+</text></view>

    <!-- 添加房间弹窗 -->
    <view class="modal" v-if="showAdd">
      <view class="modal-content">
        <view class="popup-title">添加房间</view>
        <view class="form-item">
          <text class="form-label">所属房源</text>
          <picker :range="propertyList" range-key="name" @change="onPropertyChange">
            <view class="form-picker">{{ selectedProperty?.name || '请选择' }} ›</view>
          </picker>
        </view>
        <view class="form-item">
          <text class="form-label">房间号</text>
          <input class="form-input" v-model="formData.room_number" placeholder="如: 101" />
        </view>
        <view class="form-item">
          <text class="form-label">楼层</text>
          <input class="form-input" type="number" v-model="formData.floor" placeholder="楼层" />
        </view>
        <view class="form-item">
          <text class="form-label">租金类型</text>
          <picker :range="rentTypes" range-key="label" @change="onRentTypeChange">
            <view class="form-picker">{{ rentTypes[formData.rent_type - 1]?.label }} ›</view>
          </picker>
        </view>
        <view class="form-item">
          <text class="form-label">{{ rentTypes[formData.rent_type - 1]?.label || '租金' }}</text>
          <input class="form-input" type="digit" v-model="formData.rent_amount" placeholder="租金金额" />
        </view>
        <button class="btn-primary" @click="submitForm" :disabled="submitting">{{ submitting ? '提交中...' : '提交' }}</button>
        <button class="btn-default" @click="showAdd = false">取消</button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { onPullDownRefresh } from '@dcloudio/uni-app'
import { roomApi, propertyApi } from '../../utils/request'

const loading = ref(false)
const submitting = ref(false)
const list = ref([])
const showAdd = ref(false)
const propertyList = ref([])
const selectedProperty = ref(null)
const formData = reactive({ room_number: '', floor: '', rent_type: 1, rent_amount: '' })

const rentTypes = [
  { label: '月租金', value: 1 },
  { label: '季租金', value: 2 },
  { label: '年租金', value: 3 }
]

const getStatusText = (status) => ({ 1: '空置', 2: '已租', 3: '维修中' }[status] || '未知')
const getStatusClass = (status) => ({ 1: 'info', 2: 'success', 3: 'warning' }[status] || '')
const getRentTypeLabel = (type) => ({ 1: '月租', 2: '季租', 3: '年租' }[type] || '月租')

onMounted(async () => {
  await loadList()
  await loadProperties()
})

const loadList = async () => {
  loading.value = true
  try {
    const res = await roomApi.getList({ page: 1, page_size: 100 })
    list.value = res.list || []
  } catch (error) { console.error(error) } finally { loading.value = false; uni.stopPullDownRefresh() }
}

const loadProperties = async () => {
  try {
    const res = await propertyApi.getList({ page: 1, page_size: 100 })
    propertyList.value = res.list || []
  } catch (error) { console.error(error) }
}

const addRoom = () => {
  selectedProperty.value = null
  formData.room_number = ''
  formData.floor = ''
  formData.rent_type = 1
  formData.rent_amount = ''
  showAdd.value = true
}

const onPropertyChange = (e) => {
  selectedProperty.value = propertyList.value[e.detail.value]
}

const onRentTypeChange = (e) => {
  formData.rent_type = rentTypes[e.detail.value].value
}

const submitForm = async () => {
  if (!selectedProperty.value) { uni.showToast({ title: '请选择房源', icon: 'none' }); return }
  if (!formData.room_number) { uni.showToast({ title: '请输入房间号', icon: 'none' }); return }
  submitting.value = true
  try {
    await roomApi.create({
      property_id: selectedProperty.value.id,
      room_number: formData.room_number,
      floor: parseInt(formData.floor) || 0,
      rent_type: formData.rent_type,
      rent_amount: parseFloat(formData.rent_amount) || 0,
      monthly_rent: formData.rent_type === 1 ? parseFloat(formData.rent_amount) || 0 : 0
    })
    uni.showToast({ title: '添加成功', icon: 'success' })
    showAdd.value = false
    loadList()
  } catch (error) { console.error(error) } finally { submitting.value = false }
}

const changeStatus = async (room) => {
  const newStatus = room.status === 1 ? 2 : room.status === 2 ? 1 : 1
  uni.showModal({
    title: '确认操作',
    content: `确定将房间状态改为"${getStatusText(newStatus)}"吗？`,
    success: async (res) => {
      if (res.confirm) {
        try {
          await roomApi.updateStatus(room.id, newStatus)
          uni.showToast({ title: '操作成功', icon: 'success' })
          room.status = newStatus
        } catch (error) { console.error(error) }
      }
    }
  })
}

const goTenant = (tenantId) => {
  if (tenantId) {
    uni.navigateTo({ url: `/pages/tenant/detail?id=${tenantId}` })
  }
}

const goDetail = (roomId) => {
  uni.navigateTo({ url: `/pages/room/detail?id=${roomId}` })
}

onPullDownRefresh(() => loadList())
</script>

<style scoped>
.page-container {
  min-height: 100vh;
  background: #f8f8f8;
  padding-bottom: 120rpx;
}

.room-list {
  padding: 24rpx;
}

.room-card {
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 24rpx;
}

.room-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20rpx;
}

.room-number {
  font-size: 36rpx;
  font-weight: 600;
  color: #333;
}

.room-status {
  padding: 8rpx 20rpx;
  border-radius: 8rpx;
  font-size: 24rpx;
}

.room-status.success { background: #E8F5E9; color: #4CAF50; }
.room-status.info { background: #E3F2FD; color: #2196F3; }
.room-status.warning { background: #FFF3E0; color: #FF9800; }

.room-info {
  border-top: 1rpx solid #f0f0f0;
  padding-top: 20rpx;
}

.info-row {
  display: flex;
  justify-content: space-between;
  padding: 8rpx 0;
}

.info-label { color: #999; font-size: 26rpx; }
.info-value { color: #333; font-size: 26rpx; }
.info-value.price { color: #FF6B6B; font-weight: 600; }
.info-value.tenant { color: #007AFF; }

.tenant-row {
  cursor: pointer;
}
.tenant-info {
  display: flex;
  align-items: center;
}
.tenant-arrow {
  color: #007AFF;
  font-size: 32rpx;
  margin-left: 8rpx;
}

.room-actions {
  margin-top: 20rpx;
  padding-top: 20rpx;
  border-top: 1rpx solid #f0f0f0;
}

.action-btn {
  width: 100%;
  background: #007AFF;
  color: #fff;
  border-radius: 8rpx;
  padding: 20rpx;
  font-size: 28rpx;
}

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

.modal-content .form-label { width: 160rpx; color: #666; }
.modal-content .form-input { flex: 1; font-size: 28rpx; text-align: right; }
.form-picker { flex: 1; text-align: right; color: #333; }
.modal-content .btn-primary { margin-top: 32rpx; }
.modal-content .btn-default { margin-top: 16rpx; }

.empty {
  padding: 100rpx;
  text-align: center;
  color: #999;
}

.empty-icon {
  display: block;
  font-size: 100rpx;
  margin-bottom: 16rpx;
}
</style>
