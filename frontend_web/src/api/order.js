import request from '@/utils/request'
import { query } from '@/utils/query'

export function orderCreate(data) {
  return request({
    url: '/frontend/order/create',
    method: 'post',
    data
  })
}

export function orderList(params) {
  return request({
    url: `/frontend/order/list?${query(params)}`,
    method: 'get'
  })
}

export function orderDetail(params) {
  return request({
    url: `/frontend/order/detail?${query(params)}`,
    method: 'get'
  })
}

export function orderCancel(data) {
  return request({
    url: '/frontend/order/cancel',
    method: 'post',
    data
  })
}
