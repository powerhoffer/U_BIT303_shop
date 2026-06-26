import redBullImage from '@/assets/products/red-bull.png'
import cocaColaImage from '@/assets/products/coca-cola.png'

const samples = [
  redBullImage,
  cocaColaImage
]

export function productImage(goods) {
  if (goods && goods.image_url) {
    return goods.image_url
  }
  const id = goods && goods.id ? Number(goods.id) : 1
  return samples[Math.abs(id - 1) % samples.length]
}

export function imageError(event, goods) {
  event.target.src = productImage({ id: goods && goods.id ? goods.id + 1 : 2 })
}
