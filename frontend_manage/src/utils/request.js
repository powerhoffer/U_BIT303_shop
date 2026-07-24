import axios from 'axios'
import { Message } from 'element-ui'
import store from '@/store'
import router from '@/router'
import { getToken } from '@/utils/auth'

const service = axios.create({
  baseURL: process.env.VUE_APP_BASE_API || 'http://127.0.0.1:8000',
  timeout: 15000
})

service.interceptors.request.use(
  config => {
    const token = getToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => Promise.reject(error)
)

service.interceptors.response.use(
  response => {
    const res = response.data
    if (res.code !== 0) {
      if (res.code === 401) {
        store.dispatch('user/resetToken')
        if (router.currentRoute.path !== '/login') {
          router.replace(`/login?redirect=${router.currentRoute.fullPath}`)
        }
      }
      Message({
        message: res.message || 'Request failed',
        type: 'error',
        duration: 3000
      })
      return Promise.reject(res)
    }
    return res
  },
  error => {
    Message({
      message: error.message || 'Network error',
      type: 'error',
      duration: 3000
    })
    return Promise.reject(error)
  }
)

export default service
