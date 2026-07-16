import request from '@/utils/request'
import { query } from './query'

export function permissionList(params) {
  return request({ url: `/backend/permission/list?${query(params)}`, method: 'get' })
}

export function permissionDetail(params) {
  return request({ url: `/backend/permission/detail?${query(params)}`, method: 'get' })
}

export function permissionCreate(data) {
  return request({ url: '/backend/permission/create', method: 'post', data })
}

export function permissionUpdate(data) {
  return request({ url: '/backend/permission/update', method: 'post', data })
}

export function permissionStatus(data) {
  return request({ url: '/backend/permission/status', method: 'post', data })
}
