<script>
import cloudConfig from './utils/config'

export default {
  onLaunch() {
    console.log('App Launch')

    // #ifdef MP-WEIXIN
    // 微信小程序 - 初始化云开发环境
    if (wx.cloud) {
      wx.cloud.init({
        env: cloudConfig.env,
        traceUser: true
      })
      console.log('微信云开发初始化成功, env:', cloudConfig.env)
    } else {
      console.warn('请使用 2.2.3 或以上的基础库以使用云能力')
    }
    // #endif

    // 检查登录状态
    const token = uni.getStorageSync('token')
    if (!token) {
      uni.reLaunch({ url: '/pages/login/index' })
    }
  },
  onShow() {
    console.log('App Show')
  },
  onHide() {
    console.log('App Hide')
  }
}
</script>

<style>
page {
  background-color: #f8f8f8;
  font-family: -apple-system, BlinkMacSystemFont, 'Helvetica Neue', Helvetica, sans-serif;
  font-size: 28rpx;
  color: #333;
}

.container { padding: 24rpx; }

.card {
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 24rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.05);
}

.btn-primary {
  background: linear-gradient(135deg, #007AFF 0%, #0055FF 100%);
  color: #fff;
  border-radius: 12rpx;
  padding: 24rpx;
  text-align: center;
  font-size: 32rpx;
  font-weight: 500;
}

.btn-default {
  background: #fff;
  color: #007AFF;
  border: 2rpx solid #007AFF;
  border-radius: 12rpx;
  padding: 24rpx;
  text-align: center;
  font-size: 32rpx;
}

.form-item {
  display: flex;
  align-items: center;
  padding: 24rpx 0;
  border-bottom: 1rpx solid #eee;
}

.form-label {
  width: 180rpx;
  color: #666;
  font-size: 28rpx;
}

.form-input {
  flex: 1;
  font-size: 28rpx;
  color: #333;
}

.status-tag {
  display: inline-block;
  padding: 4rpx 16rpx;
  border-radius: 8rpx;
  font-size: 24rpx;
}

.status-tag.success { background: #E8F5E9; color: #4CAF50; }
.status-tag.warning { background: #FFF3E0; color: #FF9800; }
.status-tag.danger { background: #FFEBEE; color: #F44336; }
.status-tag.info { background: #E3F2FD; color: #2196F3; }

.empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 100rpx 0;
  color: #999;
}

.empty-icon {
  font-size: 100rpx;
  margin-bottom: 20rpx;
}
</style>
