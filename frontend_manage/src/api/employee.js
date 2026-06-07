import request from '@/utils/request'
import { query } from './query'

export function login(data) {
  return request({
    url: '/backend/employee/login',
    method: 'post',
    data
  })
}

export function logout() {
  return request({
    url: '/backend/employee/logout',
    method: 'post'
  })
}

export function getInfo() {
  return request({
    url: '/backend/employee/info',
    method: 'get'
  })
}

export function updatePassword(data) {
  return request({
    url: '/backend/employee/password',
    method: 'post',
    data
  })
}

export function employeeList(params) {
  return request({
    url: `/backend/employee/manage/list?${query(params)}`,
    method: 'get'
  })
}

export function employeeCreate(data) {
  return request({
    url: '/backend/employee/manage/create',
    method: 'post',
    data
  })
}

export function employeeUpdate(data) {
  return request({
    url: '/backend/employee/manage/update',
    method: 'post',
    data
  })
}

export function employeeStatus(data) {
  return request({
    url: '/backend/employee/manage/status',
    method: 'post',
    data
  })
}

export function employeeResetPassword(data) {
  return request({
    url: '/backend/employee/manage/reset-password',
    method: 'post',
    data
  })
}
