<template>
  <view class="page-container">
    <view class="card">
      <view class="card-header">
        <text class="card-title">{{ property.name }}</text>
        <text class="status-tag" :class="property.status === 1 ? 'success' : 'warning'">
          {{ property.status === 1 ? '正常' : '已下架' }}
        </text>
      </view>
      <view class="info-item">
        <text class="info-label">地址</text>
        <text class="info-value">{{ property.address || '暂无' }}</text>
      </view>
      <view class="info-item">
        <text class="info-label">类型</text>
        <text class="info-value">{{ getPropertyType(property.property_type) }}</text>
      </view>
      <view class="info-item">
        <text class="info-label">房间数</text>
        <text class="info-value">{{ property.total_rooms }}间</text>
      </view>
    </view>

    <view class="section">
      <view class="section-header">
        <text class="section-title">房间列表</text>
        <text class="section-action" @click="addRoom">+ 添加房间</text>
      </view>
      <view class="room-list">
        <view class="room-item" v-for="room in rooms" :key="room.id" @click="goRoomDetail(room.id)">
          <view class="room-left">
            <view class="room-number">{{ room.room_number }}</view>
            <view class="room-meta">
              <text class="room-rent">¥{{ room.rent_amount || room.monthly_rent }}/{{ getRentTypeShort(room.rent_type) }}</text>
              <text class="room-status" :class="getRoomStatusClass(room.status)">
                {{ getRoomStatus(room.status) }}
              </text>
            </view>
            <view class="room-tenant" v-if="room.tenant_name" @click.stop="goTenant(room.tenant_id)">
              <text class="tenant-label">租客：</text>
              <text class="tenant-name">{{ room.tenant_name }}</text>
              <text class="tenant-phone" v-if="room.tenant_phone">（{{ room.tenant_phone }}）</text>
              <text class="tenant-arrow">›</text>
            </view>
          </view>
          <text class="room-arrow">›</text>
        </view>
        <view class="empty" v-if="rooms.length === 0">
          <text class="empty-icon">🚪</text>
          <text>暂无房间</text>
        </view>
      </view>
    </view>

    <!-- 添加/编辑房间弹窗 -->
    <view class="modal" v-if="showRoom">
      <view class="modal-content">
        <view class="popup-title">{{ editingRoom ? '编辑房间' : '添加房间' }}</view>
        <view class="form-item">
          <text class="form-label">房间号</text>
          <input class="form-input" v-model="roomForm.room_number" placeholder="如: 101" />
        </view>
        <view class="form-item">
          <text class="form-label">楼层</text>
          <input class="form-input" type="number" v-model="roomForm.floor" placeholder="楼层" />
        </view>
        <view class="form-item">
          <text class="form-label">租金类型</text>
          <picker :range="rentTypes" range-key="label" @change="onRentTypeChange">
            <view class="form-picker">{{ rentTypes[roomForm.rent_type - 1]?.label }} ›</view>
          </picker>
        </view>
        <view class="form-item">
          <text class="form-label">{{ rentTypes[roomForm.rent_type - 1]?.label || '租金' }}</text>
          <input class="form-input" type="digit" v-model="roomForm.rent_amount" placeholder="租金金额" />
        </view>
        <button class="btn-primary" @click="submitRoom" :disabled="submitting">{{ submitting ? '提交中...' : '确定' }}</button>
        <button class="btn-default" @click="showRoom = false">取消</button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { propertyApi, roomApi } from '../../utils/request'

const propertyId = ref(0)
const property = ref({})
const rooms = ref([])
const submitting = ref(false)
const showRoom = ref(false)
const editingRoom = ref(null)

const roomForm = reactive({ room_number: '', floor: '', rent_type: 1, rent_amount: '' })
const typeOptions = [{ label: '整栋', value: 1 }, { label: '单套', value: 2 }, { label: '商铺', value: 3 }]
const rentTypes = [{ label: '月租金', value: 1 }, { label: '季租金', value: 2 }, { label: '年租金', value: 3 }]

const getPropertyType = (type) => typeOptions.find(t => t.value === type)?.label || '未知'
const getRoomStatus = (status) => ({ 1: '空置', 2: '已租', 3: '维修中' }[status] || '未知')
const getRoomStatusClass = (status) => ({ 1: 'info', 2: 'success', 3: 'warning' }[status] || '')
const getRentTypeShort = (type) => ({ 1: '月', 2: '季', 3: '年' }[type] || '月')

onMounted(async () => {
  const pages = getCurrentPages()
  propertyId.value = pages[pages.length - 1].options?.id
  if (propertyId.value) {
    await loadProperty()
    await loadRooms()
  }
})

const loadProperty = async () => {
  try { property.value = await propertyApi.getDetail(propertyId.value) } catch (error) { console.error(error) }
}

const loadRooms = async () => {
  try {
    const res = await roomApi.getList({ property_id: propertyId.value, page: 1, page_size: 100 })
    rooms.value = res.list || []
  } catch (error) { console.error(error) }
}

const addRoom = () => {
  editingRoom.value = null
  roomForm.room_number = ''
  roomForm.floor = ''
  roomForm.rent_type = 1
  roomForm.rent_amount = ''
  showRoom.value = true
}

const editRoom = (room) => {
  editingRoom.value = room
  roomForm.room_number = room.room_number
  roomForm.floor = String(room.floor || '')
  roomForm.rent_type = room.rent_type || 1
  roomForm.rent_amount = String(room.rent_amount || room.monthly_rent || '')
  showRoom.value = true
}

const onRentTypeChange = (e) => {
  roomForm.rent_type = rentTypes[e.detail.value].value
}

const goTenant = (tenantId) => {
  if (tenantId) {
    uni.navigateTo({ url: `/pages/tenant/detail?id=${tenantId}` })
  }
}

const goRoomDetail = (roomId) => {
  uni.navigateTo({ url: `/pages/room/detail?id=${roomId}` })
}

const submitRoom = async () => {
  if (!roomForm.room_number) { uni.showToast({ title: '请输入房间号', icon: 'none' }); return }
  submitting.value = true
  try {
    const data = {
      property_id: propertyId.value,
      room_number: roomForm.room_number,
      floor: parseInt(roomForm.floor) || 0,
      rent_type: roomForm.rent_type,
      rent_amount: parseFloat(roomForm.rent_amount) || 0,
      monthly_rent: roomForm.rent_type === 1 ? parseFloat(roomForm.rent_amount) || 0 : 0
    }
    if (editingRoom.value) {
      await roomApi.update(editingRoom.value.id, data)
    } else {
      await roomApi.create(data)
    }
    uni.showToast({ title: '操作成功', icon: 'success' })
    showRoom.value = false
    await loadRooms()
    await loadProperty()
  } catch (error) { console.error(error) } finally { submitting.value = false }
}
</script>

<style scoped>
.page-container {
  min-height: 100vh;
  background: #f8f8f8;
  padding: 24rpx;
  padding-bottom: 160rpx;
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
  align-items: center;
  margin-bottom: 24rpx;
}

.card-title {
  font-size: 36rpx;
  font-weight: 600;
  color: #333;
}

.status-tag {
  padding: 6rpx 16rpx;
  border-radius: 8rpx;
  font-size: 24rpx;
}
.status-tag.success { background: #E8F5E9; color: #4CAF50; }
.status-tag.warning { background: #FFF3E0; color: #FF9800; }

.info-item {
  display: flex;
  padding: 16rpx 0;
  border-bottom: 1rpx solid #f0f0f0;
}

.info-item:last-child { border-bottom: none; }
.info-label { width: 140rpx; color: #999; }
.info-value { flex: 1; color: #333; }

.section {
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24rpx;
}

.section-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
}

.section-action {
  color: #007AFF;
  font-size: 28rpx;
}

.room-item {
  display: flex;
  align-items: center;
  padding: 20rpx 0;
  border-bottom: 1rpx solid #f0f0f0;
}

.room-item:last-child { border-bottom: none; }
.room-left { flex: 1; }
.room-number { font-size: 32rpx; font-weight: 600; color: #333; margin-bottom: 8rpx; }
.room-meta { display: flex; align-items: center; gap: 16rpx; margin-bottom: 8rpx; }
.room-rent { font-size: 26rpx; color: #666; }
.room-status { padding: 4rpx 12rpx; border-radius: 6rpx; font-size: 22rpx; }
.room-status.success { background: #E8F5E9; color: #4CAF50; }
.room-status.info { background: #E3F2FD; color: #2196F3; }
.room-status.warning { background: #FFF3E0; color: #FF9800; }

.room-tenant {
  font-size: 24rpx;
  color: #007AFF;
  background: #E3F2FD;
  padding: 8rpx 16rpx;
  border-radius: 8rpx;
  display: inline-flex;
  align-items: center;
}
.tenant-label { color: #666; }
.tenant-name { color: #007AFF; font-weight: 500; }
.tenant-phone { color: #999; margin-left: 4rpx; }
.tenant-arrow { color: #007AFF; font-size: 28rpx; margin-left: 8rpx; }

.room-arrow { font-size: 32rpx; color: #ccc; }

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
.modal-content .btn-primary { margin-top: 32rpx; }
.modal-content .btn-default { margin-top: 16rpx; }

.empty { padding: 60rpx; text-align: center; color: #999; }
.empty-icon { display: block; font-size: 80rpx; margin-bottom: 16rpx; }
</style>
