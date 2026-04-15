<template>
  <view class="page-container">
    <view class="property-list">
      <view class="property-item" v-for="item in list" :key="item.id" @click="goDetail(item.id)">
        <view class="property-info">
          <view class="property-name">{{ item.name }}</view>
          <view class="property-address">{{ item.address || '暂无地址' }}</view>
          <view class="property-meta">
            <text class="meta-item">{{ getPropertyType(item.property_type) }}</text>
            <text class="meta-item">{{ item.total_rooms }}间房</text>
            <text class="status-tag" :class="item.status === 1 ? 'success' : 'warning'">
              {{ item.status === 1 ? '正常' : '已下架' }}
            </text>
          </view>
        </view>
        <text class="property-arrow">›</text>
      </view>
      <view class="empty" v-if="list.length === 0 && !loading">
        <text class="empty-icon">🏠</text>
        <text>暂无房源</text>
      </view>
    </view>

    <!-- 添加按钮 -->
    <view class="add-btn" @click="addProperty"><text>+</text></view>

    <!-- 添加弹窗 -->
    <view class="modal" v-if="showAdd">
      <view class="modal-content">
        <view class="popup-title">添加房源</view>
        <view class="form-item">
          <text class="form-label">房源名称</text>
          <input class="form-input" v-model="formData.name" placeholder="如: 阳光小区3栋" />
        </view>
        <view class="form-item">
          <text class="form-label">地址</text>
          <input class="form-input" v-model="formData.address" placeholder="详细地址" />
        </view>
        <view class="form-item">
          <text class="form-label">类型</text>
          <picker :range="typeOptions" range-key="label" @change="onTypeChange">
            <view class="form-picker">{{ typeOptions[formData.property_type - 1]?.label }} ›</view>
          </picker>
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
import { propertyApi } from '../../utils/request'

const loading = ref(false)
const submitting = ref(false)
const list = ref([])
const showAdd = ref(false)

const formData = reactive({ name: '', address: '', property_type: 1 })
const typeOptions = [{ label: '整栋', value: 1 }, { label: '单套', value: 2 }, { label: '商铺', value: 3 }]

const getPropertyType = (type) => typeOptions.find(t => t.value === type)?.label || '未知'

onMounted(() => loadList())

const loadList = async (refresh = false) => {
  loading.value = true
  try {
    const res = await propertyApi.getList({ page: 1, page_size: 100 })
    list.value = res.list || []
  } catch (error) {
    console.error('加载失败', error)
  } finally {
    loading.value = false
    uni.stopPullDownRefresh()
  }
}

const addProperty = () => {
  formData.name = ''
  formData.address = ''
  formData.property_type = 1
  showAdd.value = true
}

const onTypeChange = (e) => { formData.property_type = typeOptions[e.detail.value].value }

const submitForm = async () => {
  if (!formData.name) { uni.showToast({ title: '请输入房源名称', icon: 'none' }); return }
  submitting.value = true
  try {
    await propertyApi.create(formData)
    uni.showToast({ title: '添加成功', icon: 'success' })
    showAdd.value = false
    loadList()
  } catch (error) {
    console.error('添加失败', error)
  } finally {
    submitting.value = false
  }
}

const goDetail = (id) => { uni.navigateTo({ url: `/pages/property/detail?id=${id}` }) }

onPullDownRefresh(() => loadList(true))
</script>

<style scoped>
.page-container {
  min-height: 100vh;
  background: #f8f8f8;
  padding-bottom: 120rpx;
}

.property-list {
  padding: 24rpx;
}

.property-item {
  display: flex;
  align-items: center;
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 24rpx;
}

.property-info {
  flex: 1;
}

.property-name {
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 8rpx;
}

.property-address {
  font-size: 26rpx;
  color: #999;
  margin-bottom: 12rpx;
}

.property-meta {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.meta-item {
  font-size: 24rpx;
  color: #666;
}

.property-arrow {
  font-size: 36rpx;
  color: #ccc;
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

.add-btn text {
  color: #fff;
  font-size: 56rpx;
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
  width: 160rpx;
  color: #666;
}

.modal-content .form-input {
  flex: 1;
  font-size: 28rpx;
}

.form-picker {
  flex: 1;
  text-align: right;
  color: #333;
}

.modal-content .btn-primary {
  margin-top: 32rpx;
}

.modal-content .btn-default {
  margin-top: 16rpx;
}

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
