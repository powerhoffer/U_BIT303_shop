import request from '@/utils/request'
import { query } from '@/utils/query'

export function pointsBalance() {
  return request({
    url: '/backend/points/balance',
    method: 'get'
  })
}

export function pointsRecords(params) {
  return request({
    url: `/backend/points/records?${query(params)}`,
    method: 'get'
  })
}
