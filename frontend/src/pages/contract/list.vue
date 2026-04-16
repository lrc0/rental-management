<template>
  <view class="page-container">
    <view class="contract-list">
      <view class="contract-card" v-for="contract in list" :key="contract.id" @click="showDetail(contract)">
        <view class="contract-header">
          <view class="contract-room">{{ contract.room?.room_number || '-' }}</view>
          <view class="contract-status" :class="getStatusClass(contract.status)">{{ getStatusText(contract.status) }}</view>
        </view>
        <view class="contract-body">
          <view class="contract-tenant">租客：{{ contract.tenant?.name || '-' }}</view>
          <view class="contract-info">
            <view class="info-item">
              <text class="info-label">租期</text>
              <text class="info-value">{{ formatDate(contract.start_date) }} ~ {{ formatDate(contract.end_date) }}</text>
            </view>
            <view class="info-item">
              <text class="info-label">月租金</text>
              <text class="info-value price">¥{{ contract.monthly_rent }}</text>
            </view>
          </view>
        </view>
        <view class="contract-actions" v-if="contract.status === 1">
          <button class="btn-terminate" @click.stop="terminateContract(contract)">解约</button>
        </view>
      </view>
      <view class="empty" v-if="list.length === 0 && !loading">
        <text class="empty-icon">📝</text>
        <text>暂无合同</text>
      </view>
    </view>

    <view class="add-btn" @click="addContract"><text>+</text></view>

    <!-- 签订合同弹窗 -->
    <view class="modal" v-if="showAdd">
      <view class="modal-content">
        <view class="popup-title">{{ editingContract ? '编辑合同' : '签订合同' }}</view>
        <view class="form-item" v-if="!editingContract">
          <text class="form-label">选择房间</text>
          <picker :range="availableRooms" range-key="room_number" @change="onRoomChange">
            <view class="form-picker">{{ selectedRoom?.room_number || '请选择' }} ›</view>
          </picker>
        </view>
        <view class="form-item" v-if="!editingContract">
          <text class="form-label">选择租客</text>
          <picker :range="tenantList" range-key="name" @change="onTenantChange">
            <view class="form-picker">{{ selectedTenant?.name || '请选择' }} ›</view>
          </picker>
        </view>
        <view class="form-item">
          <text class="form-label">开始日期</text>
          <picker mode="date" :value="formData.start_date" @change="(e) => formData.start_date = e.detail.value">
            <view class="form-picker">{{ formData.start_date || '请选择' }} ›</view>
          </picker>
        </view>
        <view class="form-item">
          <text class="form-label">结束日期</text>
          <picker mode="date" :value="formData.end_date" @change="(e) => formData.end_date = e.detail.value">
            <view class="form-picker">{{ formData.end_date || '请选择' }} ›</view>
          </picker>
        </view>
        <view class="form-item">
          <text class="form-label">月租金</text>
          <input class="form-input" type="digit" v-model="formData.monthly_rent" placeholder="请输入" />
        </view>
        <button class="btn-primary" @click="submitForm" :disabled="submitting">{{ submitting ? '提交中...' : '确定' }}</button>
        <button class="btn-default" @click="showAdd = false">取消</button>
      </view>
    </view>

    <!-- 详情弹窗 -->
    <view class="modal" v-if="showDetailModal">
      <view class="modal-content">
        <view class="popup-title">合同详情</view>
        <view class="detail-info" v-if="currentContract">
          <view class="detail-item"><text class="detail-label">房间</text><text class="detail-value">{{ currentContract.room?.room_number || '-' }}</text></view>
          <view class="detail-item"><text class="detail-label">租客</text><text class="detail-value">{{ currentContract.tenant?.name || '-' }}</text></view>
          <view class="detail-item"><text class="detail-label">租期</text><text class="detail-value">{{ formatDate(currentContract.start_date) }} ~ {{ formatDate(currentContract.end_date) }}</text></view>
          <view class="detail-item"><text class="detail-label">月租金</text><text class="detail-value price">¥{{ currentContract.monthly_rent }}</text></view>
          <view class="detail-item"><text class="detail-label">押金</text><text class="detail-value">¥{{ currentContract.deposit || 0 }}</text></view>
          <view class="detail-item"><text class="detail-label">状态</text><text class="detail-value" :class="getStatusClass(currentContract.status)">{{ getStatusText(currentContract.status) }}</text></view>
        </view>
        <button class="btn-primary" @click="editCurrentContract" v-if="currentContract?.status === 1">编辑</button>
        <button class="btn-danger" @click="deleteCurrentContract" v-if="currentContract?.status !== 1">删除</button>
        <button class="btn-default" @click="showDetailModal = false">关闭</button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { onPullDownRefresh } from '@dcloudio/uni-app'
import { contractApi, roomApi, tenantApi } from '../../utils/request'

const loading = ref(false)
const submitting = ref(false)
const list = ref([])
const showAdd = ref(false)
const showDetailModal = ref(false)
const editingContract = ref(null)
const currentContract = ref(null)
const availableRooms = ref([])
const tenantList = ref([])
const selectedRoom = ref(null)
const selectedTenant = ref(null)
const formData = reactive({ start_date: '', end_date: '', monthly_rent: '' })

const getStatusText = (status) => ({ 1: '生效中', 2: '已到期', 3: '已解约' }[status] || '未知')
const getStatusClass = (status) => ({ 1: 'success', 2: 'default', 3: 'warning' }[status] || '')
const formatDate = (date) => date ? date.substring(0, 10) : '-'

onMounted(async () => {
  await loadList()
  await loadAvailableRooms()
  await loadTenants()
})

const loadList = async () => {
  loading.value = true
  try {
    const res = await contractApi.getList({ page: 1, page_size: 100 })
    list.value = res.list || []
  } catch (error) { console.error(error) } finally { loading.value = false; uni.stopPullDownRefresh() }
}

const loadAvailableRooms = async () => {
  try {
    const res = await roomApi.getList({ status: 1, page: 1, page_size: 100 })
    availableRooms.value = res.list || []
  } catch (error) { console.error(error) }
}

const loadTenants = async () => {
  try {
    const res = await tenantApi.getList({ status: 1, page: 1, page_size: 100 })
    tenantList.value = res.list || []
  } catch (error) { console.error(error) }
}

const addContract = () => {
  editingContract.value = null
  selectedRoom.value = null
  selectedTenant.value = null
  formData.start_date = ''
  formData.end_date = ''
  formData.monthly_rent = ''
  showAdd.value = true
}

const onRoomChange = (e) => {
  selectedRoom.value = availableRooms.value[e.detail.value]
  formData.monthly_rent = String(selectedRoom.value?.monthly_rent || '')
}

const onTenantChange = (e) => { selectedTenant.value = tenantList.value[e.detail.value] }

const submitForm = async () => {
  if (!editingContract.value) {
    if (!selectedRoom.value) { uni.showToast({ title: '请选择房间', icon: 'none' }); return }
    if (!selectedTenant.value) { uni.showToast({ title: '请选择租客', icon: 'none' }); return }
  }
  if (!formData.start_date || !formData.end_date) { uni.showToast({ title: '请选择租期', icon: 'none' }); return }
  submitting.value = true
  try {
    if (editingContract.value) {
      await contractApi.update(editingContract.value.id, {
        start_date: formData.start_date,
        end_date: formData.end_date,
        monthly_rent: parseFloat(formData.monthly_rent) || 0
      })
    } else {
      await contractApi.create({
        room_id: selectedRoom.value.id,
        tenant_id: selectedTenant.value.id,
        start_date: formData.start_date,
        end_date: formData.end_date,
        monthly_rent: parseFloat(formData.monthly_rent) || 0,
        deposit: 0,
        payment_day: 1
      })
    }
    uni.showToast({ title: editingContract.value ? '修改成功' : '签订成功', icon: 'success' })
    showAdd.value = false
    loadList()
    loadAvailableRooms()
  } catch (error) { console.error(error) } finally { submitting.value = false }
}

const showDetail = (contract) => {
  currentContract.value = contract
  showDetailModal.value = true
}

const editCurrentContract = () => {
  showDetailModal.value = false
  editingContract.value = currentContract.value
  formData.start_date = currentContract.value.start_date?.substring(0, 10) || ''
  formData.end_date = currentContract.value.end_date?.substring(0, 10) || ''
  formData.monthly_rent = String(currentContract.value.monthly_rent || '')
  showAdd.value = true
}

const deleteCurrentContract = async () => {
  uni.showModal({
    title: '确认删除',
    content: `确定删除该合同吗？`,
    success: async (res) => {
      if (res.confirm) {
        try {
          await contractApi.delete(currentContract.value.id)
          uni.showToast({ title: '删除成功', icon: 'success' })
          showDetailModal.value = false
          loadList()
        } catch (error) { uni.showToast({ title: error.message || '删除失败', icon: 'none' }) }
      }
    }
  })
}

const terminateContract = (contract) => {
  uni.showModal({
    title: '确认解约',
    content: `确定解除与${contract.tenant?.name}的合同吗？`,
    success: async (res) => {
      if (res.confirm) {
        try {
          await contractApi.terminate(contract.id, '房东主动解约')
          uni.showToast({ title: '解约成功', icon: 'success' })
          loadList()
          loadAvailableRooms()
        } catch (error) { uni.showToast({ title: error.message || '解约失败', icon: 'none' }) }
      }
    }
  })
}

onPullDownRefresh(() => loadList())
</script>

<style scoped>
.page-container {
  min-height: 100vh;
  background: #f8f8f8;
  padding-bottom: 120rpx;
}

.contract-list { padding: 24rpx; }

.contract-card {
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 24rpx;
}

.contract-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20rpx;
}

.contract-room { font-size: 36rpx; font-weight: 600; color: #333; }

.contract-status { padding: 8rpx 20rpx; border-radius: 8rpx; font-size: 24rpx; }
.contract-status.success { background: #E8F5E9; color: #4CAF50; }
.contract-status.default { background: #f5f5f5; color: #999; }
.contract-status.warning { background: #FFF3E0; color: #FF9800; }

.contract-tenant { font-size: 26rpx; color: #666; margin-bottom: 20rpx; }

.contract-info .info-item { display: flex; justify-content: space-between; padding: 8rpx 0; }
.info-label { color: #999; font-size: 26rpx; }
.info-value { color: #333; font-size: 26rpx; }
.info-value.price { color: #FF6B6B; font-weight: 600; }

.contract-actions { margin-top: 20rpx; }
.btn-terminate { width: 100%; background: #FFF3E0; color: #FF9800; border-radius: 8rpx; padding: 20rpx; font-size: 28rpx; }

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
  max-height: 80vh;
  overflow-y: auto;
}

.popup-title { font-size: 36rpx; font-weight: 600; text-align: center; margin-bottom: 32rpx; }
.modal-content .form-item { display: flex; align-items: center; padding: 24rpx 0; border-bottom: 1rpx solid #eee; }
.modal-content .form-label { width: 160rpx; color: #666; }
.modal-content .form-input { flex: 1; text-align: right; font-size: 28rpx; }
.form-picker { flex: 1; text-align: right; color: #333; }
.modal-content .btn-primary { margin-top: 32rpx; }
.modal-content .btn-default { margin-top: 16rpx; }

.detail-item {
  display: flex;
  justify-content: space-between;
  padding: 16rpx 0;
  border-bottom: 1rpx solid #f0f0f0;
}
.detail-label { color: #999; }
.detail-value { color: #333; }
.detail-value.price { color: #FF6B6B; font-weight: 600; }
.detail-value.success { color: #4CAF50; }
.detail-value.default { color: #999; }
.detail-value.warning { color: #FF9800; }

.btn-danger {
  background: #FFEBEE;
  color: #F44336;
  border-radius: 12rpx;
  padding: 24rpx;
  text-align: center;
  font-size: 30rpx;
  margin-top: 16rpx;
}

.empty { padding: 100rpx; text-align: center; color: #999; }
.empty-icon { display: block; font-size: 100rpx; margin-bottom: 16rpx; }
</style>
