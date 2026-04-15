<template>
  <view class="page-container">
    <view class="tenant-header">
      <view class="tenant-avatar">
        <text class="avatar-text">{{ tenant.name?.charAt(0) || '租' }}</text>
      </view>
      <view class="tenant-basic">
        <text class="tenant-name">{{ tenant.name }}</text>
        <text class="tenant-phone">📱 {{ tenant.phone || '未填写' }}</text>
      </view>
      <text class="status-tag" :class="tenant.status === 1 ? 'success' : 'default'">
        {{ tenant.status === 1 ? '正常' : '已退租' }}
      </text>
    </view>

    <view class="section" v-if="tenant.room_number">
      <view class="section-title">入住信息</view>
      <view class="room-card" @click="goRoom">
        <view class="room-info">
          <text class="room-number">{{ tenant.room_number }}</text>
          <text class="room-property">{{ tenant.property_name }}</text>
        </view>
        <text class="room-arrow">›</text>
      </view>
    </view>

    <view class="section">
      <view class="section-title">基本信息</view>
      <view class="info-list">
        <view class="info-item">
          <text class="info-label">姓名</text>
          <text class="info-value">{{ tenant.name || '-' }}</text>
        </view>
        <view class="info-item">
          <text class="info-label">手机号</text>
          <text class="info-value">{{ tenant.phone || '-' }}</text>
        </view>
        <view class="info-item">
          <text class="info-label">身份证</text>
          <text class="info-value">{{ tenant.id_card || '-' }}</text>
        </view>
        <view class="info-item">
          <text class="info-label">紧急联系人</text>
          <text class="info-value">{{ tenant.emergency_contact || '-' }}</text>
        </view>
        <view class="info-item">
          <text class="info-label">紧急联系电话</text>
          <text class="info-value">{{ tenant.emergency_phone || '-' }}</text>
        </view>
        <view class="info-item">
          <text class="info-label">备注</text>
          <text class="info-value">{{ tenant.remark || '-' }}</text>
        </view>
      </view>
    </view>

    <view class="actions">
      <button class="btn-edit" @click="editTenant">编辑租客</button>
      <button class="btn-delete" @click="deleteTenant">删除租客</button>
    </view>

    <!-- 编辑弹窗 -->
    <view class="modal" v-if="showEdit">
      <view class="modal-content">
        <view class="popup-title">编辑租客</view>
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
        <view class="form-item">
          <text class="form-label">紧急联系人</text>
          <input class="form-input" v-model="formData.emergency_contact" placeholder="紧急联系人姓名" />
        </view>
        <view class="form-item">
          <text class="form-label">紧急联系电话</text>
          <input class="form-input" type="number" v-model="formData.emergency_phone" placeholder="紧急联系电话" />
        </view>
        <view class="form-item">
          <text class="form-label">备注</text>
          <input class="form-input" v-model="formData.remark" placeholder="备注信息" />
        </view>
        <button class="btn-primary" @click="submitEdit" :disabled="submitting">{{ submitting ? '保存中...' : '保存' }}</button>
        <button class="btn-default" @click="showEdit = false">取消</button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { tenantApi } from '../../utils/request'

const tenantId = ref(0)
const tenant = ref({})
const showEdit = ref(false)
const submitting = ref(false)
const formData = reactive({
  name: '',
  phone: '',
  id_card: '',
  emergency_contact: '',
  emergency_phone: '',
  remark: ''
})

onMounted(async () => {
  const pages = getCurrentPages()
  tenantId.value = pages[pages.length - 1].options?.id
  if (tenantId.value) {
    await loadTenant()
  }
})

const loadTenant = async () => {
  try {
    const res = await tenantApi.getDetail(tenantId.value)
    tenant.value = res || {}
  } catch (error) {
    console.error(error)
    uni.showToast({ title: '加载失败', icon: 'none' })
  }
}

const goRoom = () => {
  if (tenant.value.room_id) {
    uni.navigateTo({ url: `/pages/room/list` })
  }
}

const editTenant = () => {
  formData.name = tenant.value.name || ''
  formData.phone = tenant.value.phone || ''
  formData.id_card = tenant.value.id_card || ''
  formData.emergency_contact = tenant.value.emergency_contact || ''
  formData.emergency_phone = tenant.value.emergency_phone || ''
  formData.remark = tenant.value.remark || ''
  showEdit.value = true
}

const submitEdit = async () => {
  if (!formData.name) {
    uni.showToast({ title: '请输入姓名', icon: 'none' })
    return
  }
  submitting.value = true
  try {
    await tenantApi.update(tenantId.value, formData)
    uni.showToast({ title: '保存成功', icon: 'success' })
    showEdit.value = false
    await loadTenant()
  } catch (error) {
    console.error(error)
  } finally {
    submitting.value = false
  }
}

const deleteTenant = () => {
  uni.showModal({
    title: '确认删除',
    content: `确定删除租客"${tenant.value.name}"吗？`,
    success: async (res) => {
      if (res.confirm) {
        try {
          await tenantApi.delete(tenantId.value)
          uni.showToast({ title: '删除成功', icon: 'success' })
          setTimeout(() => {
            uni.navigateBack()
          }, 1000)
        } catch (error) {
          uni.showToast({ title: error.message || '删除失败', icon: 'none' })
        }
      }
    }
  })
}
</script>

<style scoped>
.page-container {
  min-height: 100vh;
  background: #f8f8f8;
  padding-bottom: 200rpx;
}

.tenant-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 60rpx 32rpx;
  display: flex;
  align-items: center;
}

.tenant-avatar {
  width: 120rpx;
  height: 120rpx;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 24rpx;
}

.avatar-text {
  color: #fff;
  font-size: 48rpx;
  font-weight: 600;
}

.tenant-basic {
  flex: 1;
}

.tenant-name {
  display: block;
  font-size: 40rpx;
  font-weight: 600;
  color: #fff;
  margin-bottom: 8rpx;
}

.tenant-phone {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.8);
}

.status-tag {
  padding: 8rpx 20rpx;
  border-radius: 8rpx;
  font-size: 24rpx;
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
}
.status-tag.success { background: rgba(76, 175, 80, 0.8); }
.status-tag.default { background: rgba(158, 158, 158, 0.8); }

.section {
  background: #fff;
  margin: 24rpx;
  border-radius: 16rpx;
  padding: 24rpx;
}

.section-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 24rpx;
}

.room-card {
  display: flex;
  align-items: center;
  background: #f8f8f8;
  border-radius: 12rpx;
  padding: 24rpx;
}

.room-info {
  flex: 1;
}

.room-number {
  display: block;
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 8rpx;
}

.room-property {
  font-size: 26rpx;
  color: #666;
}

.room-arrow {
  font-size: 36rpx;
  color: #ccc;
}

.info-list {
  background: #f8f8f8;
  border-radius: 12rpx;
}

.info-item {
  display: flex;
  justify-content: space-between;
  padding: 20rpx 24rpx;
  border-bottom: 1rpx solid #eee;
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  color: #999;
  font-size: 28rpx;
}

.info-value {
  color: #333;
  font-size: 28rpx;
}

.actions {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: #fff;
  padding: 24rpx;
  display: flex;
  gap: 24rpx;
  box-shadow: 0 -2rpx 10rpx rgba(0, 0, 0, 0.05);
}

.btn-edit {
  flex: 1;
  background: #007AFF;
  color: #fff;
  border-radius: 12rpx;
  padding: 24rpx;
  font-size: 30rpx;
}

.btn-delete {
  flex: 1;
  background: #FFEBEE;
  color: #F44336;
  border-radius: 12rpx;
  padding: 24rpx;
  font-size: 30rpx;
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
