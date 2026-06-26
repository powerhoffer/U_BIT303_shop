import request from '@/utils/request'

export function categoryList() {
  return request({
    url: '/frontend/category/list',
    method: 'get'
  })
}
