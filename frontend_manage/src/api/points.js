import request from '@/utils/request'
import { query } from './query'

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

export function pointsAdd(data) {
  return request({
    url: '/backend/points/manage/add',
    method: 'post',
    data
  })
}

export function pointsBatchAdd(data) {
  return request({
    url: '/backend/points/manage/batch-add',
    method: 'post',
    data
  })
}

export function pointsDeduct(data) {
  return request({
    url: '/backend/points/manage/deduct',
    method: 'post',
    data
  })
}

export function pointsManageRecords(params) {
  return request({
    url: `/backend/points/manage/records?${query(params)}`,
    method: 'get'
  })
}
