import request from '@/utils/request'
import { query } from '@/utils/query'

export function cartAdd(data) {
  return request({
    url: '/frontend/cart/add',
    method: 'post',
    data
  })
}

export function cartList(params) {
  return request({
    url: `/frontend/cart/list?${query(params)}`,
    method: 'get'
  })
}

export function cartUpdate(data) {
  return request({
    url: '/frontend/cart/update',
    method: 'post',
    data
  })
}

export function cartRemove(data) {
  return request({
    url: '/frontend/cart/remove',
    method: 'post',
    data
  })
}
