import request from '@/utils/request'
import { query } from './query'

export function roleList(params) {
  return request({ url: `/backend/role/list?${query(params)}`, method: 'get' })
}

export function roleDetail(params) {
  return request({ url: `/backend/role/detail?${query(params)}`, method: 'get' })
}

export function roleCreate(data) {
  return request({ url: '/backend/role/create', method: 'post', data })
}

export function roleUpdate(data) {
  return request({ url: '/backend/role/update', method: 'post', data })
}

export function roleStatus(data) {
  return request({ url: '/backend/role/status', method: 'post', data })
}

export function rolePermissions(data) {
  return request({ url: '/backend/role/permissions', method: 'post', data })
}
