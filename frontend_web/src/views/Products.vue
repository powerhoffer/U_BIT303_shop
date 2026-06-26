<template>
  <div>
    <h1 class="page-title">Products</h1>
    <div class="toolbar">
      <el-select v-model="filters.category_id" clearable placeholder="All categories" @change="search">
        <el-option v-for="item in categories" :key="item.id" :label="item.name" :value="item.id" />
      </el-select>
      <el-input
        v-model="filters.name"
        class="grow"
        clearable
        placeholder="Search products"
        @keyup.enter.native="search"
        @clear="search"
      />
      <el-button type="primary" @click="search">Search</el-button>
    </div>

    <div v-loading="loading" class="product-grid">
      <article v-for="item in goods" :key="item.id" class="product-card" @click="openDetail(item.id)">
        <div class="product-image">
          <img :src="productImage(item)" :alt="item.name" @error="imageError($event, item)">
        </div>
        <div class="product-body">
          <h2>{{ item.name }}</h2>
          <p>{{ item.category_name || 'General' }}</p>
          <div class="product-meta">
            <span class="points">{{ item.points_price }} pts</span>
            <span class="muted">Stock {{ item.stock }}</span>
          </div>
        </div>
      </article>
    </div>

    <div v-if="!loading && goods.length === 0" class="empty-panel">No products found.</div>

    <el-pagination
      v-if="total > 0"
      class="pagination"
      background
      layout="prev, pager, next"
      :current-page.sync="filters.page"
      :page-size="filters.size"
      :total="total"
      @current-change="loadGoods"
    />
  </div>
</template>

<script>
import { categoryList } from '@/api/category'
import { goodsList } from '@/api/goods'
import { productImage, imageError } from '@/utils/product-image'

export default {
  name: 'Products',
  data() {
    return {
      loading: false,
      categories: [],
      goods: [],
      total: 0,
      filters: {
        page: 1,
        size: 10,
        category_id: '',
        name: ''
      }
    }
  },
  created() {
    this.loadCategories()
    this.loadGoods()
  },
  methods: {
    productImage,
    imageError,
    async loadCategories() {
      const res = await categoryList()
      this.categories = res.data.list || []
    },
    async loadGoods() {
      this.loading = true
      try {
        const res = await goodsList(this.filters)
        this.goods = res.data.list || []
        this.total = res.data.total || 0
      } finally {
        this.loading = false
      }
    },
    search() {
      this.filters.page = 1
      this.loadGoods()
    },
    openDetail(id) {
      this.$router.push(`/products/${id}`)
    }
  }
}
</script>

<style scoped>
.product-grid {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 18px;
  min-height: 260px;
}

.product-card {
  background: #ffffff;
  border: 1px solid #e5eaf2;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.18s ease, box-shadow 0.18s ease;
}

.product-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 16px 30px rgba(31, 41, 51, 0.1);
}

.product-image {
  height: 190px;
  background: #f8fafc;
  display: flex;
  align-items: center;
  justify-content: center;
  border-bottom: 1px solid #eef2f7;
}

.product-image img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  padding: 14px;
}

.product-body {
  padding: 15px;
}

.product-body h2 {
  margin: 0 0 8px;
  min-height: 44px;
  font-size: 16px;
  line-height: 22px;
}

.product-body p {
  margin: 0 0 14px;
  color: #7b8794;
  font-size: 13px;
}

.product-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.pagination {
  margin-top: 26px;
  text-align: center;
}

@media (max-width: 1180px) {
  .product-grid {
    grid-template-columns: repeat(4, minmax(0, 1fr));
  }
}

@media (max-width: 960px) {
  .product-grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

@media (max-width: 640px) {
  .product-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 420px) {
  .product-grid {
    grid-template-columns: 1fr;
  }
}
</style>
