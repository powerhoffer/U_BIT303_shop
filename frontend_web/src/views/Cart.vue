<template>
  <div>
    <h1 class="page-title">Cart</h1>
    <div class="panel">
      <el-table v-loading="loading" :data="items" border>
        <el-table-column label="Product" min-width="260">
          <template slot-scope="{ row }">
            <div class="cart-product">
              <img :src="productImage({ id: row.goods_id, image_url: row.image_url })" :alt="row.goods_name" @error="imageError($event, row)">
              <div>
                <strong>{{ row.goods_name }}</strong>
                <p>{{ row.category_name }}</p>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="Price" width="120">
          <template slot-scope="{ row }"><span class="points">{{ row.points_price }} Credits</span></template>
        </el-table-column>
        <el-table-column label="Stock" prop="stock" width="100" />
        <el-table-column label="Quantity" width="170">
          <template slot-scope="{ row }">
            <el-input-number :value="row.count" :min="0" :max="999" size="small" @change="value => updateCount(row, value)" />
          </template>
        </el-table-column>
        <el-table-column label="Subtotal" width="130">
          <template slot-scope="{ row }"><span class="points">{{ row.total_points }} Credits</span></template>
        </el-table-column>
        <el-table-column label="Actions" width="120">
          <template slot-scope="{ row }">
            <el-button type="text" @click="remove(row)">Remove</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="cart-footer">
        <div class="total">Total <strong>{{ totalPoints }} Credits</strong></div>
        <el-button type="primary" :disabled="items.length === 0" :loading="checkoutLoading" @click="checkout">Redeem</el-button>
      </div>
    </div>
  </div>
</template>

<script>
import { cartList, cartUpdate, cartRemove } from '@/api/cart'
import { orderCreate } from '@/api/order'
import { productImage, imageError } from '@/utils/product-image'

export default {
  name: 'Cart',
  data() {
    return {
      loading: false,
      checkoutLoading: false,
      items: []
    }
  },
  computed: {
    totalPoints() {
      return this.items.reduce((sum, item) => sum + Number(item.total_points || 0), 0)
    }
  },
  created() {
    this.loadCart()
  },
  methods: {
    productImage,
    imageError,
    async loadCart() {
      this.loading = true
      try {
        const res = await cartList({ page: 1, size: 50 })
        this.items = res.data.list || []
      } finally {
        this.loading = false
      }
    },
    async updateCount(row, count) {
      await cartUpdate({ id: row.id, count })
      this.$message.success(count === 0 ? 'Item removed' : 'Cart updated')
      this.loadCart()
    },
    async remove(row) {
      await cartRemove({ id: row.id })
      this.$message.success('Item removed')
      this.loadCart()
    },
    async checkout() {
      this.checkoutLoading = true
      try {
        const res = await orderCreate({ remark: '' })
        this.$message.success('Order created')
        this.$router.push(`/orders/${res.data.order.id}`)
      } finally {
        this.checkoutLoading = false
      }
    }
  }
}
</script>

<style scoped>
.cart-product {
  display: flex;
  align-items: center;
  gap: 14px;
}

.cart-product img {
  width: 72px;
  height: 72px;
  object-fit: contain;
  background: #f8fafc;
  border-radius: 6px;
}

.cart-product p {
  margin: 6px 0 0;
  color: #7b8794;
}

.cart-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 24px;
  margin-top: 22px;
}

.total strong {
  margin-left: 8px;
  color: #f59e0b;
  font-size: 22px;
}
</style>
