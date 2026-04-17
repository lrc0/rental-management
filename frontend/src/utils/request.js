// API请求封装 - 支持微信云托管
import cloudConfig from './config'

// 微信小程序云托管请求
const cloudContainerRequest = (options) => {
  return new Promise((resolve, reject) => {
    const token = uni.getStorageSync('token')

    // 构建请求头
    const header = {
      'Content-Type': 'application/json',
      ...options.header
    }

    // 添加Authorization
    if (token) {
      header['Authorization'] = 'Bearer ' + token
    }

    // #ifdef MP-WEIXIN
    // 微信小程序 - 使用 wx.cloud.callContainer
    wx.cloud.callContainer({
      config: {
        env: cloudConfig.env
      },
      path: cloudConfig.apiPrefix + options.url,
      method: options.method || 'GET',
      header: {
        'X-WX-SERVICE': cloudConfig.serviceName,
        ...header
      },
      data: options.data || {},
      success: function(res) {
        console.log('云托管请求成功:', res)
        if (res.statusCode === 200) {
          if (res.data.code === 0) {
            resolve(res.data.data)
          } else if (res.data.code === 401 || res.data.code === 10004 || res.data.code === 10005) {
            uni.removeStorageSync('token')
            uni.removeStorageSync('userInfo')
            uni.reLaunch({ url: '/pages/login/index' })
            reject(res.data)
          } else {
            uni.showToast({ title: res.data.message || '请求失败', icon: 'none' })
            reject(res.data)
          }
        } else {
          uni.showToast({ title: '网络错误: ' + res.statusCode, icon: 'none' })
          reject(res)
        }
      },
      fail: function(err) {
        console.error('云托管请求失败:', err)
        uni.showToast({ title: '请求失败', icon: 'none' })
        reject(err)
      }
    })
    // #endif

    // #ifdef H5
    // H5环境 - 使用普通HTTP请求
    uni.request({
      url: cloudConfig.apiPrefix + options.url,
      method: options.method || 'GET',
      data: options.data || {},
      header: header,
      success: function(res) {
        if (res.statusCode === 200) {
          if (res.data.code === 0) {
            resolve(res.data.data)
          } else if (res.data.code === 401 || res.data.code === 10004 || res.data.code === 10005) {
            uni.removeStorageSync('token')
            uni.removeStorageSync('userInfo')
            uni.reLaunch({ url: '/pages/login/index' })
            reject(res.data)
          } else {
            uni.showToast({ title: res.data.message || '请求失败', icon: 'none' })
            reject(res.data)
          }
        } else {
          uni.showToast({ title: '网络错误: ' + res.statusCode, icon: 'none' })
          reject(res)
        }
      },
      fail: function(err) {
        console.error('请求失败:', err)
        uni.showToast({ title: '网络请求失败', icon: 'none' })
        reject(err)
      }
    })
    // #endif

    // #ifndef MP-WEIXIN || H5
    // 其他环境（开发环境）
    uni.request({
      url: 'http://localhost:8080' + cloudConfig.apiPrefix + options.url,
      method: options.method || 'GET',
      data: options.data || {},
      header: header,
      success: function(res) {
        if (res.statusCode === 200) {
          if (res.data.code === 0) {
            resolve(res.data.data)
          } else if (res.data.code === 401 || res.data.code === 10004 || res.data.code === 10005) {
            uni.removeStorageSync('token')
            uni.removeStorageSync('userInfo')
            uni.reLaunch({ url: '/pages/login/index' })
            reject(res.data)
          } else {
            uni.showToast({ title: res.data.message || '请求失败', icon: 'none' })
            reject(res.data)
          }
        } else {
          uni.showToast({ title: '网络错误: ' + res.statusCode, icon: 'none' })
          reject(res)
        }
      },
      fail: function(err) {
        console.error('请求失败:', err)
        uni.showToast({ title: '网络请求失败', icon: 'none' })
        reject(err)
      }
    })
    // #endif
  })
}

export const api = {
  get: function(url, data) {
    return cloudContainerRequest({ url: url, method: 'GET', data: data })
  },
  post: function(url, data) {
    return cloudContainerRequest({ url: url, method: 'POST', data: data })
  },
  put: function(url, data) {
    return cloudContainerRequest({ url: url, method: 'PUT', data: data })
  },
  delete: function(url, data) {
    return cloudContainerRequest({ url: url, method: 'DELETE', data: data })
  }
}

export const authApi = {
  login: function(data) { return api.post('/auth/login', data) },
  register: function(data) { return api.post('/auth/register', data) },
  getProfile: function() { return api.get('/auth/profile') },
  updateProfile: function(data) { return api.put('/auth/profile', data) },
  changePassword: function(data) { return api.put('/auth/password', data) },
  getStatistics: function() { return api.get('/statistics') }
}

export const propertyApi = {
  getList: function(params) { return api.get('/properties', params) },
  getDetail: function(id) { return api.get('/properties/' + id) },
  create: function(data) { return api.post('/properties', data) },
  update: function(id, data) { return api.put('/properties/' + id, data) },
  delete: function(id) { return api.delete('/properties/' + id) }
}

export const roomApi = {
  getList: function(params) { return api.get('/rooms', params) },
  getDetail: function(id) { return api.get('/rooms/' + id) },
  create: function(data) { return api.post('/rooms', data) },
  update: function(id, data) { return api.put('/rooms/' + id, data) },
  updateStatus: function(id, status) { return api.put('/rooms/' + id + '/status', { status: status }) },
  delete: function(id) { return api.delete('/rooms/' + id) }
}

export const tenantApi = {
  getList: function(params) { return api.get('/tenants', params) },
  getDetail: function(id) { return api.get('/tenants/' + id) },
  create: function(data) { return api.post('/tenants', data) },
  update: function(id, data) { return api.put('/tenants/' + id, data) },
  delete: function(id) { return api.delete('/tenants/' + id) }
}

export const contractApi = {
  getList: function(params) { return api.get('/contracts', params) },
  getDetail: function(id) { return api.get('/contracts/' + id) },
  create: function(data) { return api.post('/contracts', data) },
  update: function(id, data) { return api.put('/contracts/' + id, data) },
  terminate: function(id, reason) { return api.put('/contracts/' + id + '/terminate', { reason: reason }) },
  delete: function(id) { return api.delete('/contracts/' + id) }
}

export const meterApi = {
  getList: function(params) { return api.get('/meter-readings', params) },
  create: function(data) { return api.post('/meter-readings', data) },
  delete: function(id) { return api.delete('/meter-readings/' + id) }
}

export const billApi = {
  getList: function(params) { return api.get('/bills', params) },
  getDetail: function(id) { return api.get('/bills/' + id) },
  create: function(data) { return api.post('/bills', data) },
  update: function(id, data) { return api.put('/bills/' + id, data) },
  delete: function(id) { return api.delete('/bills/' + id) },
  pay: function(id, data) { return api.put('/bills/' + id + '/pay', data) },
  preview: function(params) { return api.get('/bills/preview', params) },
  getStatistics: function(params) { return api.get('/bills/statistics', params) },
  getMonthlyStatistics: function(params) { return api.get('/bills/monthly-statistics', params) }
}

export const feeRateApi = {
  get: function() { return api.get('/fee-rates') },
  update: function(data) { return api.put('/fee-rates', data) }
}
