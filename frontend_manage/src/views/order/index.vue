<template>
  <div class="page-container">
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="queryForm" size="mini">
        <el-form-item label="Order No">
          <el-input v-model.trim="queryForm.order_no" clearable placeholder="Order no" />
        </el-form-item>
        <el-form-item label="Employee ID">
          <el-input v-model.trim="queryForm.employee_id" clearable placeholder="Employee ID" />
        </el-form-item>
        <el-form-item label="Status">
          <el-select v-model="queryForm.status" clearable placeholder="All" style="width: 140px">
            <el-option label="Pending" :value="1" />
            <el-option label="Completed" :value="2" />
            <el-option label="Cancelled" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="el-icon-search" @click="handleSearch">Search</el-button>
          <el-button icon="el-icon-refresh-left" @click="handleReset">Reset</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card shadow="never">
      <div slot="header" class="card-header">
        <span>Order Management</span>
      </div>

      <el-table v-loading="loading" :data="list" border style="width: 100%">
        <el-table-column align="center" prop="id" label="ID" width="80" />
        <el-table-column prop="order_no" label="Order No" min-width="180" />
        <el-table-column align="center" prop="employee_id" label="Employee ID" width="120" />
        <el-table-column align="center" prop="total_points" label="Total Credits" width="130" />
        <el-table-column align="center" label="Status" width="120">
          <template slot-scope="scope">
            <el-tag :type="statusType(scope.row.status)">
              {{ statusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="Remark" min-width="160">
          <template slot-scope="scope">
            {{ scope.row.remark || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="Created At" min-width="170" />
        <el-table-column align="center" label="Actions" width="260" fixed="right">
          <template slot-scope="scope">
            <el-button size="mini" icon="el-icon-view" @click="openDetail(scope.row)">Detail</el-button>
            <el-button
              v-if="scope.row.status === 1"
              size="mini"
              type="success"
              :loading="actionLoadingId === scope.row.id && actionType === 'complete'"
              @click="completeOrder(scope.row)"
            >
              Complete
            </el-button>
            <el-button
              v-if="scope.row.status === 1"
              size="mini"
              type="warning"
              :loading="actionLoadingId === scope.row.id && actionType === 'cancel'"
              @click="cancelOrder(scope.row)"
            >
              Cancel
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="pagination"
        :current-page.sync="page"
        :page-size="size"
        :page-sizes="[10, 20, 50]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="getList"
      />
    </el-card>

    <el-dialog title="Order Detail" :visible.sync="detailVisible" width="860px">
      <div v-loading="detailLoading">
        <div v-if="detail.id" class="detail-list">
          <div class="detail-row">
            <span class="detail-label">ID</span>
            <span class="detail-value">{{ detail.id }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">Order No</span>
            <span class="detail-value">{{ detail.order_no }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">Employee ID</span>
            <span class="detail-value">{{ detail.employee_id }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">Total Credits</span>
            <span class="detail-value">{{ detail.total_points }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">Status</span>
            <span class="detail-value">
              <el-tag :type="statusType(detail.status)">{{ statusText(detail.status) }}</el-tag>
            </span>
          </div>
          <div class="detail-row">
            <span class="detail-label">Remark</span>
            <span class="detail-value">{{ detail.remark || '-' }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">Created At</span>
            <span class="detail-value">{{ detail.created_at }}</span>
          </div>
        </div>

        <div class="detail-section-title">Goods Items</div>
        <el-table :data="detail.items || []" border>
          <el-table-column align="center" prop="goods_id" label="Goods ID" width="90" />
          <el-table-column align="center" label="Image" width="100">
            <template slot-scope="scope">
              <el-image
                v-if="scope.row.goods_image_url"
                class="goods-image"
                :src="scope.row.goods_image_url"
                fit="cover"
                :preview-src-list="[scope.row.goods_image_url]"
              />
              <span v-else class="image-empty">No Image</span>
            </template>
          </el-table-column>
          <el-table-column prop="goods_name" label="Goods Name" min-width="170" />
          <el-table-column align="center" prop="points_price" label="Credits Price" width="120" />
          <el-table-column align="center" prop="count" label="Count" width="90" />
          <el-table-column align="center" prop="total_points" label="Total Credits" width="120" />
          <el-table-column prop="created_at" label="Created At" min-width="170" />
        </el-table>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import {
  orderList,
  orderDetail,
  orderComplete,
  orderCancel
} from '@/api/order'

const STATUS_MAP = {
  1: 'Pending',
  2: 'Completed',
  3: 'Cancelled'
}

const STATUS_TYPE_MAP = {
  1: 'warning',
  2: 'success',
  3: 'info'
}

export default {
  name: 'OrderList',
  data() {
    return {
      loading: false,
      detailLoading: false,
      actionLoadingId: undefined,
      actionType: '',
      list: [],
      total: 0,
      page: 1,
      size: 10,
      queryForm: {
        order_no: '',
        employee_id: '',
        status: ''
      },
      detailVisible: false,
      detail: {}
    }
  },
  created() {
    this.getList()
  },
  methods: {
    statusText(status) {
      return STATUS_MAP[status] || 'Unknown'
    },
    statusType(status) {
      return STATUS_TYPE_MAP[status] || ''
    },
    async getList() {
      this.loading = true
      try {
        const params = {
          page: this.page,
          size: this.size,
          order_no: this.queryForm.order_no,
          employee_id: this.queryForm.employee_id,
          status: this.queryForm.status === '' ? -1 : this.queryForm.status
        }
        const res = await orderList(params)
        this.list = res.data.list || []
        this.total = res.data.total || 0
      } finally {
        this.loading = false
      }
    },
    handleSearch() {
      this.page = 1
      this.getList()
    },
    handleReset() {
      this.queryForm = { order_no: '', employee_id: '', status: '' }
      this.page = 1
      this.getList()
    },
    handleSizeChange(val) {
      this.size = val
      this.page = 1
      this.getList()
    },
    async openDetail(row) {
      this.detailVisible = true
      this.detailLoading = true
      this.detail = {}
      try {
        const res = await orderDetail({ id: row.id })
        this.detail = res.data.order || {}
      } finally {
        this.detailLoading = false
      }
    },
    completeOrder(row) {
      this.$confirm(`Confirm to complete order ${row.order_no}?`, 'Warning', {
        type: 'warning'
      }).then(async() => {
        await this.runOrderAction(row, 'complete')
      }).catch(() => {})
    },
    cancelOrder(row) {
      this.$confirm(`Confirm to cancel order ${row.order_no}?`, 'Warning', {
        type: 'warning'
      }).then(async() => {
        await this.runOrderAction(row, 'cancel')
      }).catch(() => {})
    },
    async runOrderAction(row, type) {
      this.actionLoadingId = row.id
      this.actionType = type
      try {
        if (type === 'complete') {
          await orderComplete({ id: row.id })
          this.$message.success('Order completed successfully')
        } else {
          await orderCancel({ id: row.id })
          this.$message.success('Order cancelled successfully')
        }
        this.getList()
      } finally {
        this.actionLoadingId = undefined
        this.actionType = ''
      }
    }
  }
}
</script>

<style scoped>
.page-container {
  padding: 24px;
}
.filter-card {
  margin-bottom: 16px;
}
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.goods-image {
  width: 56px;
  height: 56px;
  border-radius: 4px;
}
.image-empty {
  color: #909399;
  font-size: 12px;
}
.detail-list {
  border: 1px solid #ebeef5;
  border-bottom: 0;
}
.detail-row {
  display: flex;
  min-height: 42px;
  border-bottom: 1px solid #ebeef5;
}
.detail-label {
  flex: 0 0 140px;
  padding: 12px;
  color: #606266;
  background: #f5f7fa;
}
.detail-value {
  flex: 1;
  padding: 12px;
  color: #303133;
  word-break: break-all;
}
.detail-section-title {
  margin: 18px 0 10px;
  font-weight: 600;
  color: #303133;
}
.pagination {
  margin-top: 16px;
  text-align: right;
}
</style>
