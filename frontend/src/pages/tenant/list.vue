<template>
  <view class="page-container">
    <view class="tenant-list">
      <view class="tenant-card" v-for="tenant in list" :key="tenant.id" @click="showDetail(tenant)">
        <view class="tenant-avatar">
          <text class="avatar-text">{{ tenant.name?.charAt(0) || '租' }}</text>
        </view>
        <view class="tenant-info">
          <view class="tenant-name">{{ tenant.name }}</view>
          <view class="tenant-contact">📱 {{ tenant.phone || '未填写' }}</view>
          <view class="tenant-room" v-if="tenant.room_number">🏠 {{ tenant.property_name }} - {{ tenant.room_number }}</view>
        </view>
        <view class="tenant-status">
          <text class="status-tag" :class="tenant.status === 1 ? 'success' : 'default'">
            {{ tenant.status === 1 ? '正常' : '已退租' }}
          </text>
        </view>
      </view>
      <view class="empty" v-if="list.length === 0 && !loading">
        <text class="empty-icon">👥</text>
        <text>暂无租客</text>
      </view>
    </view>

    <view class="add-btn" @click="addTenant"><text>+</text></view>

    <!-- 添加弹窗 -->
    <view class="modal" v-if="showAdd">
      <view class="modal-content">
        <view class="popup-title">{{ editingTenant ? '编辑租客' : '添加租客' }}</view>
        <view class="form-item">
          <text class="form-label">姓名</text>
          <input class="form-input" v-model="formData.name" placeholder="请输入姓名" />
        </view>
        <view class="form-item">
          <text class="form-label">手机号</text>
          <input class="form-input" type="number" v-model="formData.phone" placeholder="请输入手机号" maxlength="11" />
        </view>
        <view class="form-item">
          <text class="form-label">身份证</text>
          <input class="form-input" v-model="formData.id_card" placeholder="请输入身份证号" maxlength="18" />
        </view>
        <button class="btn-primary" @click="submitForm" :disabled="submitting">{{ submitting ? '提交中...' : '确定' }}</button>
        <button class="btn-default" @click="showAdd = false">取消</button>
      </view>
    </view>

    <!-- 详情弹窗 -->
    <view class="modal" v-if="showDetailModal">
      <view class="modal-content">
        <view class="popup-title">租客详情</view>
        <view class="detail-info" v-if="currentTenant">
          <view class="detail-item"><text class="detail-label">姓名</text><text class="detail-value">{{ currentTenant.name }}</text></view>
          <view class="detail-item"><text class="detail-label">手机号</text><text class="detail-value">{{ currentTenant.phone || '-' }}</text></view>
          <view class="detail-item"><text class="detail-label">身份证</text><text class="detail-value">{{ currentTenant.id_card || '-' }}</text></view>
          <view class="detail-item" v-if="currentTenant.room_number">
            <text class="detail-label">入住房间</text>
            <text class="detail-value highlight">{{ currentTenant.property_name }} - {{ currentTenant.room_number }}</text>
          </view>
        </view>
        <button class="btn-primary" @click="editCurrentTenant">编辑</button>
        <button class="btn-danger" @click="deleteCurrentTenant">删除</button>
        <button class="btn-default" @click="showDetailModal = false">关闭</button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { onPullDownRefresh } from '@dcloudio/uni-app'
import { tenantApi } from '../../utils/request'

const loading = ref(false)
const submitting = ref(false)
const list = ref([])
const showAdd = ref(false)
const showDetailModal = ref(false)
const editingTenant = ref(null)
const currentTenant = ref(null)
const formData = reactive({ name: '', phone: '', id_card: '' })

onMounted(() => loadList())

const loadList = async () => {
  loading.value = true
  try {
    const res = await tenantApi.getList({ page: 1, page_size: 100 })
    list.value = res.list || []
  } catch (error) { console.error(error) } finally { loading.value = false; uni.stopPullDownRefresh() }
}

const addTenant = () => {
  editingTenant.value = null
  formData.name = ''
  formData.phone = ''
  formData.id_card = ''
  showAdd.value = true
}

const showDetail = (tenant) => {
  currentTenant.value = tenant
  showDetailModal.value = true
}

const editCurrentTenant = () => {
  showDetailModal.value = false
  editingTenant.value = currentTenant.value
  formData.name = currentTenant.value.name
  formData.phone = currentTenant.value.phone || ''
  formData.id_card = currentTenant.value.id_card || ''
  showAdd.value = true
}

const submitForm = async () => {
  if (!formData.name) { uni.showToast({ title: '请输入姓名', icon: 'none' }); return }
  submitting.value = true
  try {
    if (editingTenant.value) {
      await tenantApi.update(editingTenant.value.id, formData)
    } else {
      await tenantApi.create(formData)
    }
    uni.showToast({ title: '操作成功', icon: 'success' })
    showAdd.value = false
    loadList()
  } catch (error) { console.error(error) } finally { submitting.value = false }
}

const deleteCurrentTenant = async () => {
  uni.showModal({
    title: '确认删除',
    content: `确定删除租客"${currentTenant.value.name}"吗？`,
    success: async (res) => {
      if (res.confirm) {
        try {
          await tenantApi.delete(currentTenant.value.id)
          uni.showToast({ title: '删除成功', icon: 'success' })
          showDetailModal.value = false
          loadList()
        } catch (error) { uni.showToast({ title: error.message || '删除失败', icon: 'none' }) }
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

.tenant-list {
  padding: 24rpx;
}

.tenant-card {
  display: flex;
  align-items: center;
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 24rpx;
}

.tenant-avatar {
  width: 88rpx;
  height: 88rpx;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 24rpx;
}

.avatar-text {
  color: #fff;
  font-size: 36rpx;
  font-weight: 600;
}

.tenant-info {
  flex: 1;
}

.tenant-name {
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 8rpx;
}

.tenant-contact {
  font-size: 24rpx;
  color: #999;
}

.tenant-room {
  font-size: 24rpx;
  color: #007AFF;
  margin-top: 4rpx;
}

.status-tag.success { background: #E8F5E9; color: #4CAF50; }
.status-tag.default { background: #f5f5f5; color: #999; }

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

.modal-content .form-label { width: 140rpx; color: #666; }
.modal-content .form-input { flex: 1; font-size: 28rpx; }
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
.detail-value.highlight { color: #007AFF; font-weight: 500; }

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
