import request from '@/utils/request'
import { query } from './query'

export function orderList(params) {
  return request({
    url: `/backend/order/list?${query(params)}`,
    method: 'get'
  })
}

export function orderDetail(params) {
  return request({
    url: `/backend/order/detail?${query(params)}`,
    method: 'get'
  })
}

export function orderComplete(data) {
  return request({
    url: '/backend/order/complete',
    method: 'post',
    data
  })
}

export function orderCancel(data) {
  return request({
    url: '/backend/order/cancel',
    method: 'post',
    data
  })
}
