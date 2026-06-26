import request from '@/utils/request'

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
