const apiBaseUrl = (process.env.VUE_APP_BASE_API || 'http://127.0.0.1:8000').replace(/\/$/, '')

export function productImage(goods) {
  const imageUrl = goods && goods.image_url ? goods.image_url.trim() : ''
  if (!imageUrl) {
    return ''
  }
  if (/^(https?:)?\/\//i.test(imageUrl) || imageUrl.startsWith('data:')) {
    return imageUrl
  }
  return `${apiBaseUrl}/${imageUrl.replace(/^\/+/, '')}`
}

export function imageError(event) {
  event.target.onerror = null
  event.target.style.visibility = 'hidden'
}
