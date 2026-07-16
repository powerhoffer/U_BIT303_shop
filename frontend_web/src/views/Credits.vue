<template>
  <div>
    <h1 class="page-title">My Credits</h1>
    <div class="credit-summary">
      <span>Available Credits</span>
      <strong>{{ balance }} Credits</strong>
    </div>
    <div class="panel">
      <el-table v-loading="loading" :data="records" border>
        <el-table-column label="Type" width="120">
          <template slot-scope="{ row }">{{ row.change_type === 1 ? 'Credit' : 'Debit' }}</template>
        </el-table-column>
        <el-table-column label="Credits" width="120">
          <template slot-scope="{ row }"><span class="points">{{ row.points }} Credits</span></template>
        </el-table-column>
        <el-table-column label="Before" prop="before_balance" width="120" />
        <el-table-column label="After" prop="after_balance" width="120" />
        <el-table-column label="Remark" prop="remark" min-width="220" />
        <el-table-column label="Created At" prop="created_at" min-width="180" />
      </el-table>
      <el-pagination
        v-if="total > 0"
        class="pagination"
        background
        layout="prev, pager, next"
        :current-page.sync="page"
        :page-size="size"
        :total="total"
        @current-change="loadRecords"
      />
    </div>
  </div>
</template>

<script>
import { pointsBalance, pointsRecords } from '@/api/points'

export default {
  name: 'Credits',
  data() {
    return {
      loading: false,
      balance: 0,
      records: [],
      total: 0,
      page: 1,
      size: 10
    }
  },
  created() {
    this.loadBalance()
    this.loadRecords()
  },
  methods: {
    async loadBalance() {
      const res = await pointsBalance()
      this.balance = res.data.balance || 0
    },
    async loadRecords() {
      this.loading = true
      try {
        const res = await pointsRecords({ page: this.page, size: this.size })
        this.records = res.data.list || []
        this.total = res.data.total || 0
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>
.credit-summary {
  display: flex;
  align-items: center;
  justify-content: space-between;
  max-width: 420px;
  margin-bottom: 22px;
  padding: 24px;
  background: #ffffff;
  border: 1px solid #e5eaf2;
  border-radius: 8px;
}

.credit-summary span {
  color: #52606d;
}

.credit-summary strong {
  color: #f59e0b;
  font-size: 30px;
}

.pagination {
  margin-top: 22px;
  text-align: center;
}
</style>
