import request from '@/utils/request'
import { query } from './query'

export function login(data) {
  return request({ url: '/backend/admin/login', method: 'post', data })
}

export function logout() {
  return request({ url: '/backend/admin/logout', method: 'post' })
}

export function getInfo() {
  return request({ url: '/backend/admin/info', method: 'get' })
}

export function updatePassword(data) {
  return request({ url: '/backend/admin/password', method: 'post', data })
}

export function adminList(params) {
  return request({ url: `/backend/admin/manage/list?${query(params)}`, method: 'get' })
}

export function adminDetail(params) {
  return request({ url: `/backend/admin/manage/detail?${query(params)}`, method: 'get' })
}

export function adminCreate(data) {
  return request({ url: '/backend/admin/manage/create', method: 'post', data })
}

export function adminUpdate(data) {
  return request({ url: '/backend/admin/manage/update', method: 'post', data })
}

export function adminStatus(data) {
  return request({ url: '/backend/admin/manage/status', method: 'post', data })
}

export function adminResetPassword(data) {
  return request({ url: '/backend/admin/manage/reset-password', method: 'post', data })
}

export function adminRoles(data) {
  return request({ url: '/backend/admin/manage/roles', method: 'post', data })
}
