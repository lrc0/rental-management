// API基础配置
// H5部署时使用相对路径，通过nginx代理
// #ifdef H5
const BASE_URL = '/api/v1'
// #endif
// #ifndef H5
const BASE_URL = 'http://localhost:8080/api/v1'
// #endif

const request = (options) => {
  return new Promise((resolve, reject) => {
    const token = uni.getStorageSync('token')
    uni.request({
      url: BASE_URL + options.url,
      method: options.method || 'GET',
      data: options.data || {},
      header: {
        'Content-Type': 'application/json',
        'Authorization': token ? `Bearer ${token}` : '',
        ...options.header
      },
      success: (res) => {
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
          uni.showToast({ title: '网络错误', icon: 'none' })
          reject(res)
        }
      },
      fail: (err) => {
        uni.showToast({ title: '网络请求失败', icon: 'none' })
        reject(err)
      }
    })
  })
}

export const api = {
  get: (url, data) => request({ url, method: 'GET', data }),
  post: (url, data) => request({ url, method: 'POST', data }),
  put: (url, data) => request({ url, method: 'PUT', data }),
  delete: (url, data) => request({ url, method: 'DELETE', data })
}

export const authApi = {
  login: (data) => api.post('/auth/login', data),
  register: (data) => api.post('/auth/register', data),
  getProfile: () => api.get('/auth/profile'),
  updateProfile: (data) => api.put('/auth/profile', data),
  changePassword: (data) => api.put('/auth/password', data)
}

export const propertyApi = {
  getList: (params) => api.get('/properties', params),
  getDetail: (id) => api.get(`/properties/${id}`),
  create: (data) => api.post('/properties', data),
  update: (id, data) => api.put(`/properties/${id}`, data),
  delete: (id) => api.delete(`/properties/${id}`)
}

export const roomApi = {
  getList: (params) => api.get('/rooms', params),
  getDetail: (id) => api.get(`/rooms/${id}`),
  create: (data) => api.post('/rooms', data),
  update: (id, data) => api.put(`/rooms/${id}`, data),
  updateStatus: (id, status) => api.put(`/rooms/${id}/status`, { status }),
  delete: (id) => api.delete(`/rooms/${id}`)
}

export const tenantApi = {
  getList: (params) => api.get('/tenants', params),
  getDetail: (id) => api.get(`/tenants/${id}`),
  create: (data) => api.post('/tenants', data),
  update: (id, data) => api.put(`/tenants/${id}`, data),
  delete: (id) => api.delete(`/tenants/${id}`)
}

export const contractApi = {
  getList: (params) => api.get('/contracts', params),
  getDetail: (id) => api.get(`/contracts/${id}`),
  create: (data) => api.post('/contracts', data),
  terminate: (id, reason) => api.put(`/contracts/${id}/terminate`, { reason })
}

export const meterApi = {
  getList: (params) => api.get('/meter-readings', params),
  create: (data) => api.post('/meter-readings', data)
}

export const billApi = {
  getList: (params) => api.get('/bills', params),
  getDetail: (id) => api.get(`/bills/${id}`),
  create: (data) => api.post('/bills', data),
  pay: (id, data) => api.put(`/bills/${id}/pay`, data),
  getStatistics: (params) => api.get('/bills/statistics', params),
  getMonthlyStatistics: (params) => api.get('/bills/monthly-statistics', params)
}

export const feeRateApi = {
  get: () => api.get('/fee-rates'),
  update: (data) => api.put('/fee-rates', data)
}
