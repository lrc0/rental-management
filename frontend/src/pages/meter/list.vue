<template>
  <view class="page-container">
    <view class="meter-list">
      <view class="meter-card" v-for="item in list" :key="item.id">
        <view class="meter-header">
          <view class="meter-room">{{ item.room?.room_number || '-' }}</view>
          <view class="meter-date">{{ formatDate(item.reading_date) }}</view>
        </view>
        <view class="meter-data">
          <view class="data-item">
            <view class="data-icon" style="background: #E3F2FD;">💧</view>
            <view class="data-info">
              <text class="data-label">水表</text>
              <text class="data-reading">{{ item.water_reading }}</text>
            </view>
          </view>
          <view class="data-item">
            <view class="data-icon" style="background: #FFF3E0;">⚡</view>
            <view class="data-info">
              <text class="data-label">电表</text>
              <text class="data-reading">{{ item.electricity_reading }}</text>
            </view>
          </view>
          <view class="data-item">
            <view class="data-icon" style="background: #FCE4EC;">🔥</view>
            <view class="data-info">
              <text class="data-label">气表</text>
              <text class="data-reading">{{ item.gas_reading }}</text>
            </view>
          </view>
        </view>
      </view>
      <view class="empty" v-if="list.length === 0 && !loading">
        <text class="empty-icon">⚡</text>
        <text>暂无抄表记录</text>
      </view>
    </view>

    <view class="add-btn" @click="addMeter"><text>+</text></view>

    <view class="modal" v-if="showAdd">
      <view class="modal-content">
        <view class="popup-title">录入抄表</view>
        <view class="form-item">
          <text class="form-label">选择房间</text>
          <picker :range="roomList" range-key="room_number" @change="onRoomChange">
            <view class="form-picker">{{ selectedRoom?.room_number || '请选择' }} ›</view>
          </picker>
        </view>
        <view class="form-item">
          <text class="form-label">水表读数</text>
          <input class="form-input" type="digit" v-model="formData.water_reading" placeholder="请输入" />
        </view>
        <view class="form-item">
          <text class="form-label">电表读数</text>
          <input class="form-input" type="digit" v-model="formData.electricity_reading" placeholder="请输入" />
        </view>
        <view class="form-item">
          <text class="form-label">气表读数</text>
          <input class="form-input" type="digit" v-model="formData.gas_reading" placeholder="请输入" />
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
import { meterApi, roomApi } from '../../utils/request'

const loading = ref(false)
const submitting = ref(false)
const list = ref([])
const showAdd = ref(false)
const roomList = ref([])
const selectedRoom = ref(null)
const formData = reactive({ water_reading: '', electricity_reading: '', gas_reading: '' })

const formatDate = (date) => date ? date.substring(0, 10) : '-'

onMounted(async () => {
  await loadList()
  await loadRooms()
})

const loadList = async () => {
  loading.value = true
  try {
    const res = await meterApi.getList({ page: 1, page_size: 100 })
    list.value = res.list || []
  } catch (error) { console.error(error) } finally { loading.value = false; uni.stopPullDownRefresh() }
}

const loadRooms = async () => {
  try {
    const res = await roomApi.getList({ page: 1, page_size: 100 })
    roomList.value = res.list || []
  } catch (error) { console.error(error) }
}

const addMeter = () => {
  selectedRoom.value = null
  formData.water_reading = ''
  formData.electricity_reading = ''
  formData.gas_reading = ''
  showAdd.value = true
}

const onRoomChange = (e) => { selectedRoom.value = roomList.value[e.detail.value] }

const submitForm = async () => {
  if (!selectedRoom.value) { uni.showToast({ title: '请选择房间', icon: 'none' }); return }
  submitting.value = true
  try {
    const today = new Date()
    const dateStr = `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}-${String(today.getDate()).padStart(2, '0')}`
    await meterApi.create({
      room_id: selectedRoom.value.id,
      reading_date: dateStr,
      water_reading: parseFloat(formData.water_reading) || 0,
      electricity_reading: parseFloat(formData.electricity_reading) || 0,
      gas_reading: parseFloat(formData.gas_reading) || 0
    })
    uni.showToast({ title: '提交成功', icon: 'success' })
    showAdd.value = false
    loadList()
  } catch (error) { console.error(error) } finally { submitting.value = false }
}

onPullDownRefresh(() => loadList())
</script>

<style scoped>
.page-container {
  min-height: 100vh;
  background: #f8f8f8;
  padding-bottom: 120rpx;
}

.meter-list { padding: 24rpx; }

.meter-card {
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 24rpx;
}

.meter-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20rpx;
}

.meter-room { font-size: 32rpx; font-weight: 600; color: #333; }
.meter-date { font-size: 24rpx; color: #999; }

.meter-data { display: flex; justify-content: space-around; padding: 20rpx 0; border-top: 1rpx solid #f0f0f0; }

.data-item { display: flex; flex-direction: column; align-items: center; }
.data-icon { width: 64rpx; height: 64rpx; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 32rpx; margin-bottom: 12rpx; }
.data-info { text-align: center; }
.data-label { display: block; font-size: 24rpx; color: #999; margin-bottom: 8rpx; }
.data-reading { font-size: 32rpx; font-weight: 600; color: #333; }

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

.modal-content { width: 100%; background: #fff; border-radius: 24rpx 24rpx 0 0; padding: 32rpx; }
.popup-title { font-size: 36rpx; font-weight: 600; text-align: center; margin-bottom: 32rpx; }
.modal-content .form-item { display: flex; align-items: center; padding: 24rpx 0; border-bottom: 1rpx solid #eee; }
.modal-content .form-label { width: 180rpx; color: #666; }
.modal-content .form-input { flex: 1; text-align: right; font-size: 28rpx; }
.form-picker { flex: 1; text-align: right; color: #333; }
.modal-content .btn-primary { margin-top: 32rpx; }
.modal-content .btn-default { margin-top: 16rpx; }

.empty { padding: 100rpx; text-align: center; color: #999; }
.empty-icon { display: block; font-size: 100rpx; margin-bottom: 16rpx; }
</style>
