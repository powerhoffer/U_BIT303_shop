<template>
  <div v-loading="loading">
    <el-button class="back-button" icon="el-icon-arrow-left" @click="$router.push('/products')">Products</el-button>
    <section v-if="goods" class="detail-panel">
      <div class="detail-image">
        <img :src="productImage(goods)" :alt="goods.name" @error="imageError($event, goods)">
      </div>
      <div class="detail-info">
        <p class="category">{{ goods.category_name }}</p>
        <h1>{{ goods.name }}</h1>
        <p class="description">{{ goods.description || 'No description available.' }}</p>
        <div class="price-row">
          <span>Price</span>
          <strong>{{ goods.points_price }} Credits</strong>
        </div>
        <div class="stock-row">Stock {{ goods.stock }}</div>
        <div class="actions">
          <el-input-number v-model="count" :min="1" :max="999" />
          <el-button type="primary" @click="addCart">Add to Cart</el-button>
        </div>
      </div>
    </section>
  </div>
</template>

<script>
import { goodsDetail } from '@/api/goods'
import { cartAdd } from '@/api/cart'
import { getToken } from '@/utils/auth'
import { productImage, imageError } from '@/utils/product-image'

export default {
  name: 'ProductDetail',
  data() {
    return {
      loading: false,
      goods: null,
      count: 1
    }
  },
  created() {
    this.loadDetail()
  },
  methods: {
    productImage,
    imageError,
    async loadDetail() {
      this.loading = true
      try {
        const res = await goodsDetail({ id: this.$route.params.id })
        this.goods = res.data.goods
      } finally {
        this.loading = false
      }
    },
    async addCart() {
      if (!getToken()) {
        this.$router.push(`/login?redirect=${encodeURIComponent(this.$route.fullPath)}`)
        return
      }
      await cartAdd({ goods_id: this.goods.id, count: this.count })
      this.$message.success('Added to cart')
    }
  }
}
</script>

<style scoped>
.back-button {
  margin-bottom: 18px;
}

.detail-panel {
  display: grid;
  grid-template-columns: minmax(320px, 46%) 1fr;
  gap: 34px;
  background: #ffffff;
  border: 1px solid #e5eaf2;
  border-radius: 8px;
  padding: 30px;
}

.detail-image {
  min-height: 420px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f8fafc;
  border-radius: 8px;
}

.detail-image img {
  width: 100%;
  max-height: 420px;
  object-fit: contain;
  padding: 24px;
}

.category {
  margin: 0 0 12px;
  color: #2563eb;
  font-weight: 700;
}

.detail-info h1 {
  margin: 0 0 20px;
  font-size: 34px;
  line-height: 1.2;
}

.description {
  min-height: 120px;
  margin: 0 0 26px;
  color: #52606d;
  font-size: 16px;
  line-height: 1.8;
}

.price-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  max-width: 340px;
  padding: 18px 0;
  border-top: 1px solid #e5eaf2;
  border-bottom: 1px solid #e5eaf2;
}

.price-row strong {
  color: #f59e0b;
  font-size: 28px;
}

.stock-row {
  margin: 18px 0 28px;
  color: #7b8794;
}

.actions {
  display: flex;
  align-items: center;
  gap: 14px;
}

@media (max-width: 820px) {
  .detail-panel {
    grid-template-columns: 1fr;
  }
}
</style>
