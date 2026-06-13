import request from '@/utils/request'
import { query } from './query'

export function goodsList(params) {
  return request({
    url: `/backend/goods/list?${query(params)}`,
    method: 'get'
  })
}

export function goodsDetail(params) {
  return request({
    url: `/backend/goods/detail?${query(params)}`,
    method: 'get'
  })
}

export function goodsCreate(data) {
  return request({
    url: '/backend/goods/create',
    method: 'post',
    data
  })
}

export function goodsUpdate(data) {
  return request({
    url: '/backend/goods/update',
    method: 'post',
    data
  })
}

export function goodsStatus(data) {
  return request({
    url: '/backend/goods/status',
    method: 'post',
    data
  })
}
