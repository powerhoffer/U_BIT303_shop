<template>
  <div>
    <h1 class="page-title">My Orders</h1>
    <div class="panel">
      <el-table v-loading="loading" :data="orders" border>
        <el-table-column label="Order No." prop="order_no" min-width="190" />
        <el-table-column label="Total" width="130">
          <template slot-scope="{ row }"><span class="points">{{ row.total_points }} Credits</span></template>
        </el-table-column>
        <el-table-column label="Status" width="130">
          <template slot-scope="{ row }">
            <el-tag :type="statusType(row.status)">{{ statusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="Created At" prop="created_at" min-width="180" />
        <el-table-column label="Actions" width="170">
          <template slot-scope="{ row }">
            <el-button type="text" @click="$router.push(`/orders/${row.id}`)">Detail</el-button>
            <el-button v-if="row.status === 1" type="text" @click="cancel(row)">Cancel</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-if="total > 0"
        class="pagination"
        background
        layout="prev, pager, next"
        :current-page.sync="page"
        :page-size="size"
        :total="total"
        @current-change="loadOrders"
      />
    </div>
  </div>
</template>

<script>
import { orderList, orderCancel } from '@/api/order'

export default {
  name: 'Orders',
  data() {
    return {
      loading: false,
      orders: [],
      total: 0,
      page: 1,
      size: 10
    }
  },
  created() {
    this.loadOrders()
  },
  methods: {
    statusText(status) {
      return ({ 1: 'Pending', 2: 'Completed', 3: 'Canceled' })[status] || 'Unknown'
    },
    statusType(status) {
      return ({ 1: 'warning', 2: 'success', 3: 'info' })[status] || ''
    },
    async loadOrders() {
      this.loading = true
      try {
        const res = await orderList({ page: this.page, size: this.size })
        this.orders = res.data.list || []
        this.total = res.data.total || 0
      } finally {
        this.loading = false
      }
    },
    async cancel(row) {
      await this.$confirm('Cancel this order?', 'Confirm', { type: 'warning' })
      await orderCancel({ id: row.id })
      this.$message.success('Order canceled')
      this.loadOrders()
    }
  }
}
</script>

<style scoped>
.pagination {
  margin-top: 22px;
  text-align: center;
}
</style>
