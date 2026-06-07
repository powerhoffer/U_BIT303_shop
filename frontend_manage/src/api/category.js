import request from '@/utils/request'

export function categoryList() {
  return request({
    url: '/backend/category/list',
    method: 'get'
  })
}
