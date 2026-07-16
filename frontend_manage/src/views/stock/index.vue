<template>
  <div class="page-container">
    <el-card shadow="never" class="operation-card">
      <div slot="header">Stock Adjustment</div>
      <el-form ref="adjustForm" :inline="true" :model="adjustForm" :rules="adjustRules" size="mini">
        <el-form-item label="Goods" prop="goods_id">
          <el-select
            v-model="adjustForm.goods_id"
            filterable
            remote
            reserve-keyword
            placeholder="Search goods"
            :remote-method="searchGoods"
            :loading="goodsLoading"
            style="width: 260px"
          >
            <el-option
              v-for="item in goodsOptions"
              :key="item.id"
              :label="`${item.name} (Stock: ${item.stock})`"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="Action" prop="action">
          <el-radio-group v-model="adjustForm.action">
            <el-radio-button label="increase">Increase</el-radio-button>
            <el-radio-button label="decrease">Decrease</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="Quantity" prop="quantity">
          <el-input-number v-model="adjustForm.quantity" :min="1" :precision="0" />
        </el-form-item>
        <el-form-item label="Remark">
          <el-input v-model.trim="adjustForm.remark" maxlength="255" placeholder="Adjustment reason" style="width: 240px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="submitLoading" @click="submitAdjustment">Apply</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="queryForm" size="mini">
        <el-form-item label="Goods ID">
          <el-input-number v-model="queryForm.goods_id" :min="1" :precision="0" controls-position="right" />
        </el-form-item>
        <el-form-item label="Change Type">
          <el-select v-model="queryForm.change_type" clearable placeholder="All" style="width: 220px">
            <el-option v-for="item in changeTypes" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="Created Date">
          <el-date-picker
            v-model="queryForm.dateRange"
            type="daterange"
            value-format="yyyy-MM-dd"
            start-placeholder="Start date"
            end-placeholder="End date"
            range-separator="to"
            clearable
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="el-icon-search" @click="handleSearch">Search</el-button>
          <el-button icon="el-icon-refresh-left" @click="handleReset">Reset</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card shadow="never">
      <div slot="header" class="card-header">
        <span>Stock Records</span>
        <el-button size="mini" icon="el-icon-refresh" @click="getRecords">Refresh</el-button>
      </div>
      <el-table v-loading="loading" :data="records" border style="width: 100%">
        <el-table-column align="center" prop="id" label="ID" width="80" />
        <el-table-column align="center" prop="goods_id" label="Goods ID" width="100" />
        <el-table-column prop="goods_name" label="Goods Name" min-width="180" />
        <el-table-column label="Change Type" min-width="190">
          <template slot-scope="scope">
            <el-tag :type="changeTypeTag(scope.row.change_type)">{{ changeTypeLabel(scope.row.change_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="center" prop="change_quantity" label="Quantity" width="100" />
        <el-table-column align="center" prop="before_stock" label="Before" width="90" />
        <el-table-column align="center" prop="after_stock" label="After" width="90" />
        <el-table-column prop="biz_type" label="Business Type" min-width="150" />
        <el-table-column align="center" prop="operator_id" label="Operator ID" width="110" />
        <el-table-column prop="remark" label="Remark" min-width="180" show-overflow-tooltip />
        <el-table-column prop="created_at" label="Created At" min-width="180" />
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
  </div>
</template>

<script>
import { goodsList } from '@/api/goods'
import { stockAdjust, stockRecords } from '@/api/stock'

const changeTypes = [
  { value: 1, label: 'Initial Stock' },
  { value: 2, label: 'Admin Increase' },
  { value: 3, label: 'Admin Decrease' },
  { value: 4, label: 'Order Deduction' },
  { value: 5, label: 'Order Cancellation Restore' }
]

export default {
  name: 'StockList',
  data() {
    return {
      loading: false,
      goodsLoading: false,
      submitLoading: false,
      goodsOptions: [],
      records: [],
      total: 0,
      page: 1,
      size: 10,
      changeTypes,
      adjustForm: {
        goods_id: undefined,
        action: 'increase',
        quantity: 1,
        remark: ''
      },
      queryForm: {
        goods_id: undefined,
        change_type: '',
        dateRange: []
      },
      adjustRules: {
        goods_id: [{ required: true, message: 'Please select goods', trigger: 'change' }],
        action: [{ required: true, message: 'Please select an action', trigger: 'change' }],
        quantity: [{ required: true, message: 'Please enter a quantity', trigger: 'blur' }]
      }
    }
  },
  created() {
    this.searchGoods('')
    this.getRecords()
  },
  methods: {
    async searchGoods(keyword) {
      this.goodsLoading = true
      try {
        const res = await goodsList({ page: 1, size: 50, name: keyword, status: -1 })
        this.goodsOptions = res.data.list || []
      } finally {
        this.goodsLoading = false
      }
    },
    submitAdjustment() {
      this.$refs.adjustForm.validate(async valid => {
        if (!valid) return
        this.submitLoading = true
        try {
          const res = await stockAdjust(this.adjustForm)
          this.$message.success(`Stock updated to ${res.data.stock}`)
          this.adjustForm.quantity = 1
          this.adjustForm.remark = ''
          await Promise.all([this.searchGoods(''), this.getRecords()])
        } finally {
          this.submitLoading = false
        }
      })
    },
    async getRecords() {
      this.loading = true
      try {
        const res = await stockRecords({
          page: this.page,
          size: this.size,
          goods_id: this.queryForm.goods_id,
          change_type: this.queryForm.change_type,
          start_time: this.queryForm.dateRange && this.queryForm.dateRange[0],
          end_time: this.queryForm.dateRange && this.queryForm.dateRange[1]
        })
        this.records = res.data.list || []
        this.total = res.data.total || 0
      } finally {
        this.loading = false
      }
    },
    changeTypeLabel(value) {
      const item = changeTypes.find(type => type.value === value)
      return item ? item.label : 'Unknown'
    },
    changeTypeTag(value) {
      if (value === 2 || value === 5) return 'success'
      if (value === 3 || value === 4) return 'warning'
      return 'info'
    },
    handleSearch() {
      this.page = 1
      this.getRecords()
    },
    handleReset() {
      this.queryForm = { goods_id: undefined, change_type: '', dateRange: [] }
      this.page = 1
      this.getRecords()
    },
    handleSizeChange(value) {
      this.size = value
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
.operation-card,
.filter-card {
  margin-bottom: 16px;
}
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.pagination {
  margin-top: 16px;
  text-align: right;
}
</style>
