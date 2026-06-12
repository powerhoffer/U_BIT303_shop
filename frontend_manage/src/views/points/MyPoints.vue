<template>
  <div class="page-container">
    <el-row :gutter="16">
      <el-col :xs="24" :sm="10" :md="8">
        <el-card shadow="never" class="balance-card">
          <div class="balance-label">Current Credit Balance</div>
          <div class="balance-value">{{ balance }}</div>
          <el-button size="mini" icon="el-icon-refresh" @click="loadData">Refresh</el-button>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="14" :md="16">
        <el-card shadow="never">
          <div slot="header">My Credit Records</div>
          <el-table v-loading="loading" :data="records" border>
            <el-table-column align="center" prop="id" label="ID" width="80" />
            <el-table-column align="center" label="Type" width="100">
              <template slot-scope="scope">
                <el-tag :type="scope.row.change_type === 1 ? 'success' : 'danger'">
                  {{ scope.row.change_type === 1 ? 'Add' : 'Deduct' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column align="center" prop="points" label="Credits" width="120" />
            <el-table-column align="center" prop="before_balance" label="Before" width="120" />
            <el-table-column align="center" prop="after_balance" label="After" width="120" />
            <el-table-column prop="remark" label="Remark" min-width="180" />
            <el-table-column prop="created_at" label="Time" min-width="180" />
          </el-table>
          <el-pagination
            class="pagination"
            :current-page.sync="page"
            :page-size="size"
            :page-sizes="[10, 20, 50]"
            :total="total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="getRecords"
          />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { pointsBalance, pointsRecords } from '@/api/points'

export default {
  name: 'MyPoints',
  data() {
    return {
      balance: 0,
      loading: false,
      records: [],
      total: 0,
      page: 1,
      size: 10
    }
  },
  created() {
    this.loadData()
  },
  methods: {
    loadData() {
      this.getBalance()
      this.getRecords()
    },
    async getBalance() {
      const res = await pointsBalance()
      this.balance = res.data.balance || 0
    },
    async getRecords() {
      this.loading = true
      try {
        const res = await pointsRecords({ page: this.page, size: this.size })
        this.records = res.data.list || []
        this.total = res.data.total || 0
      } finally {
        this.loading = false
      }
    },
    handleSizeChange(val) {
      this.size = val
      this.page = 1
      this.getRecords()
    }
  }
}
</script>

<style scoped>
.page-container {
  padding: 24px;
}
.balance-card {
  margin-bottom: 16px;
}
.balance-label {
  color: #909399;
}
.balance-value {
  margin: 16px 0;
  font-size: 42px;
  font-weight: 600;
  color: #303133;
}
.pagination {
  margin-top: 16px;
  text-align: right;
}
</style>
