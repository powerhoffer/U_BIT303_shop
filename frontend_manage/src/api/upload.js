import request from '@/utils/request'

export function uploadGoodsImage(file) {
  const data = new FormData()
  data.append('file', file)
  return request({
    url: '/backend/upload/goods-image',
    method: 'post',
    data
  })
}
