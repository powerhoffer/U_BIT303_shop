import request from '@/utils/request'
import { query } from '@/utils/query'

export function goodsList(params) {
  return request({
    url: `/frontend/goods/list?${query(params)}`,
    method: 'get'
  })
}

export function goodsDetail(params) {
  return request({
    url: `/frontend/goods/detail?${query(params)}`,
    method: 'get'
  })
}
