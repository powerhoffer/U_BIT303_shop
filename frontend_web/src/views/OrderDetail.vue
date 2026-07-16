<template>
  <div v-loading="loading">
    <el-button class="back-button" icon="el-icon-arrow-left" @click="$router.push('/orders')">My Orders</el-button>
    <div v-if="order" class="panel">
      <div class="order-head">
        <div>
          <h1>{{ order.order_no }}</h1>
          <p class="muted">{{ order.created_at }}</p>
        </div>
        <div>
          <el-tag :type="statusType(order.status)">{{ statusText(order.status) }}</el-tag>
        </div>
      </div>
      <div class="summary-row">
        <span>Total</span>
        <strong>{{ order.total_points }} Credits</strong>
      </div>
      <el-table :data="order.items || []" border>
        <el-table-column label="Product" min-width="260">
          <template slot-scope="{ row }">
            <div class="order-product">
              <img :src="productImage({ id: row.goods_id, image_url: row.goods_image_url })" :alt="row.goods_name" @error="imageError($event, row)">
              <strong>{{ row.goods_name }}</strong>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="Price" width="120">
          <template slot-scope="{ row }">{{ row.points_price }} Credits</template>
        </el-table-column>
        <el-table-column label="Quantity" prop="count" width="110" />
        <el-table-column label="Subtotal" width="130">
          <template slot-scope="{ row }"><span class="points">{{ row.total_points }} Credits</span></template>
        </el-table-column>
      </el-table>
      <div v-if="order.status === 1" class="detail-actions">
        <el-button type="danger" @click="cancel">Cancel Order</el-button>
      </div>
    </div>
  </div>
</template>

<script>
import { orderDetail, orderCancel } from '@/api/order'
import { productImage, imageError } from '@/utils/product-image'

export default {
  name: 'OrderDetail',
  data() {
    return {
      loading: false,
      order: null
    }
  },
  created() {
    this.loadOrder()
  },
  methods: {
    productImage,
    imageError,
    statusText(status) {
      return ({ 1: 'Pending', 2: 'Completed', 3: 'Canceled' })[status] || 'Unknown'
    },
    statusType(status) {
      return ({ 1: 'warning', 2: 'success', 3: 'info' })[status] || ''
    },
    async loadOrder() {
      this.loading = true
      try {
        const res = await orderDetail({ id: this.$route.params.id })
        this.order = res.data.order
      } finally {
        this.loading = false
      }
    },
    async cancel() {
      await this.$confirm('Cancel this order?', 'Confirm', { type: 'warning' })
      await orderCancel({ id: this.order.id })
      this.$message.success('Order canceled')
      this.loadOrder()
    }
  }
}
</script>

<style scoped>
.back-button {
  margin-bottom: 18px;
}

.order-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 20px;
  margin-bottom: 22px;
}

.order-head h1 {
  margin: 0 0 8px;
  font-size: 24px;
}

.summary-row {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 14px;
  margin-bottom: 18px;
}

.summary-row strong {
  color: #f59e0b;
  font-size: 24px;
}

.order-product {
  display: flex;
  align-items: center;
  gap: 14px;
}

.order-product img {
  width: 64px;
  height: 64px;
  object-fit: contain;
  background: #f8fafc;
  border-radius: 6px;
}

.detail-actions {
  margin-top: 22px;
  text-align: right;
}
</style>
