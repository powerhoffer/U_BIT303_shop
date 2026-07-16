import request from '@/utils/request'
import { query } from './query'

export function stockAdjust(data) {
  return request({ url: '/backend/stock/adjust', method: 'post', data })
}

export function stockRecords(params) {
  return request({ url: `/backend/stock/records?${query(params)}`, method: 'get' })
}
