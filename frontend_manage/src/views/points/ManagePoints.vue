<template>
  <div class="page-container">
    <el-card shadow="never" class="operation-card">
      <div slot="header" class="card-header">
        <span>Credit Operations</span>
        <el-button type="primary" size="mini" icon="el-icon-user" @click="openBatchDialog">Batch Allocate</el-button>
      </div>
      <el-form ref="form" :inline="true" :model="form" :rules="rules" size="mini">
        <el-form-item label="Employee" prop="employee_id">
          <el-select
            v-model="form.employee_id"
            filterable
            remote
            reserve-keyword
            placeholder="Search username or name"
            :remote-method="searchEmployees"
            :loading="employeeLoading"
            style="width: 220px"
            @change="loadEmployeeRecords"
          >
            <el-option
              v-for="item in employeeOptions"
              :key="item.id"
              :label="`${item.username} ${item.real_name || ''}`"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="Credits" prop="points">
          <el-input-number v-model="form.points" :min="1" :max="4294967295" :precision="0" />
        </el-form-item>
        <el-form-item label="Remark" prop="remark">
          <el-input v-model.trim="form.remark" maxlength="255" placeholder="Operation remark" style="width: 260px" />
        </el-form-item>
        <el-form-item>
          <el-button type="success" :loading="submitLoading" @click="submit('add')">Add</el-button>
          <el-button type="danger" :loading="submitLoading" @click="submit('deduct')">Deduct</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card shadow="never">
      <div slot="header" class="card-header">
        <span>Employee Credit Records</span>
        <el-button size="mini" icon="el-icon-refresh" :disabled="!form.employee_id" @click="loadEmployeeRecords">Refresh</el-button>
      </div>
      <el-table v-loading="recordsLoading" :data="records" border>
        <el-table-column align="center" prop="id" label="ID" width="80" />
        <el-table-column align="center" prop="employee_id" label="Employee ID" width="110" />
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
        @current-change="loadEmployeeRecords"
      />
    </el-card>

    <el-dialog
      title="Batch Allocate Credits"
      :visible.sync="batchVisible"
      width="900px"
      top="5vh"
      custom-class="batch-points-dialog"
      @closed="resetBatchDialog"
    >
      <el-form :inline="true" :model="batchQuery" size="mini" class="batch-filter">
        <el-form-item label="Username">
          <el-input v-model.trim="batchQuery.username" clearable placeholder="Employee username" @keyup.enter.native="handleBatchSearch" />
        </el-form-item>
        <el-form-item label="Name">
          <el-input v-model.trim="batchQuery.real_name" clearable placeholder="Employee name" @keyup.enter.native="handleBatchSearch" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="el-icon-search" @click="handleBatchSearch">Search</el-button>
          <el-button icon="el-icon-refresh-left" @click="resetBatchSearch">Reset</el-button>
        </el-form-item>
      </el-form>

      <div class="batch-toolbar">
        <span>Selected <strong>{{ selectedCount }}</strong> / 200 employees</span>
        <el-button size="mini" icon="el-icon-delete" :disabled="selectedCount === 0" @click="clearBatchSelection">Clear Selection</el-button>
      </div>

      <el-table
        ref="batchTable"
        v-loading="batchLoading"
        :data="batchEmployees"
        :row-key="row => row.id"
        border
        max-height="320"
        @selection-change="handleBatchSelectionChange"
      >
        <el-table-column type="selection" width="50" align="center" />
        <el-table-column align="center" prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="Username" min-width="170" />
        <el-table-column prop="real_name" label="Name" min-width="160" />
        <el-table-column prop="email" label="Email" min-width="220" />
      </el-table>
      <el-pagination
        class="pagination"
        small
        :current-page.sync="batchPage"
        :page-size="batchSize"
        :page-sizes="[10, 20, 50]"
        :total="batchTotal"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleBatchSizeChange"
        @current-change="loadBatchEmployees"
      />

      <el-divider />
      <el-form ref="batchForm" :model="batchForm" :rules="batchRules" label-width="110px" class="batch-form">
        <div class="batch-form-grid">
          <el-form-item label="Credits Each" prop="points">
            <el-input-number v-model="batchForm.points" :min="1" :max="4294967295" :precision="0" />
          </el-form-item>
          <el-form-item label="Total Credits">
            <span class="total-points">{{ formatNumber(batchTotalPoints) }}</span>
          </el-form-item>
        </div>
        <el-form-item label="Remark" prop="remark">
          <el-input v-model.trim="batchForm.remark" maxlength="255" show-word-limit placeholder="Batch allocation remark" />
        </el-form-item>
      </el-form>

      <div slot="footer">
        <el-button @click="batchVisible = false">Cancel</el-button>
        <el-button type="primary" :loading="batchSubmitLoading" @click="submitBatchAllocation">Allocate</el-button>
      </div>
    </el-dialog>

    <el-dialog
      title="Batch Allocation Result"
      :visible.sync="batchResultVisible"
      width="720px"
      custom-class="batch-result-dialog"
    >
      <div class="result-summary">
        <span>Processed <strong>{{ batchResult.processed_count }}</strong> employees</span>
        <span>Total <strong>{{ formatNumber(batchResult.total_points) }}</strong> credits</span>
      </div>
      <el-table :data="batchResult.list" border max-height="420">
        <el-table-column align="center" prop="employee_id" label="Employee ID" width="120" />
        <el-table-column prop="username" label="Username" min-width="150" />
        <el-table-column prop="real_name" label="Name" min-width="150" />
        <el-table-column align="center" prop="balance" label="New Balance" width="130" />
      </el-table>
      <div slot="footer">
        <el-button type="primary" @click="batchResultVisible = false">Done</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { employeeList } from '@/api/employee'
import { pointsAdd, pointsBatchAdd, pointsDeduct, pointsManageRecords } from '@/api/points'

const emptyBatchResult = () => ({
  processed_count: 0,
  total_points: 0,
  list: []
})

export default {
  name: 'ManagePoints',
  data() {
    return {
      employeeLoading: false,
      submitLoading: false,
      recordsLoading: false,
      employeeOptions: [],
      records: [],
      total: 0,
      page: 1,
      size: 10,
      form: {
        employee_id: undefined,
        points: 1,
        remark: ''
      },
      rules: {
        employee_id: [{ required: true, message: 'Please select an employee', trigger: 'change' }],
        points: [{ required: true, message: 'Please enter credits', trigger: 'blur' }]
      },
      batchVisible: false,
      batchLoading: false,
      batchSubmitLoading: false,
      batchEmployees: [],
      batchTotal: 0,
      batchPage: 1,
      batchSize: 20,
      batchQuery: {
        username: '',
        real_name: ''
      },
      batchForm: {
        points: 1,
        remark: ''
      },
      batchRules: {
        points: [{ required: true, message: 'Please enter credits', trigger: 'blur' }]
      },
      selectedEmployeeMap: {},
      syncingBatchSelection: false,
      batchResultVisible: false,
      batchResult: emptyBatchResult()
    }
  },
  computed: {
    selectedEmployees() {
      return Object.keys(this.selectedEmployeeMap).map(id => this.selectedEmployeeMap[id])
    },
    selectedCount() {
      return this.selectedEmployees.length
    },
    batchTotalPoints() {
      return this.selectedCount * Number(this.batchForm.points || 0)
    }
  },
  created() {
    this.searchEmployees('')
  },
  methods: {
    async searchEmployees(keyword) {
      this.employeeLoading = true
      try {
        const params = { page: 1, size: 20, status: 1 }
        if (keyword) {
          params.username = keyword
        }
        const res = await employeeList(params)
        this.employeeOptions = res.data.list || []
      } finally {
        this.employeeLoading = false
      }
    },
    submit(type) {
      this.$refs.form.validate(async valid => {
        if (!valid) return
        this.submitLoading = true
        try {
          const payload = {
            employee_id: this.form.employee_id,
            points: this.form.points,
            remark: this.form.remark
          }
          if (type === 'add') {
            await pointsAdd(payload)
          } else {
            await pointsDeduct(payload)
          }
          this.$message.success('Credit operation completed')
          this.form.points = 1
          this.form.remark = ''
          this.page = 1
          this.loadEmployeeRecords()
        } finally {
          this.submitLoading = false
        }
      })
    },
    async loadEmployeeRecords() {
      if (!this.form.employee_id) {
        this.records = []
        this.total = 0
        return
      }
      this.recordsLoading = true
      try {
        const res = await pointsManageRecords({
          employee_id: this.form.employee_id,
          page: this.page,
          size: this.size
        })
        this.records = res.data.list || []
        this.total = res.data.total || 0
      } finally {
        this.recordsLoading = false
      }
    },
    handleSizeChange(val) {
      this.size = val
      this.page = 1
      this.loadEmployeeRecords()
    },
    openBatchDialog() {
      this.batchQuery = { username: '', real_name: '' }
      this.batchForm = { points: 1, remark: '' }
      this.batchPage = 1
      this.batchSize = 20
      this.selectedEmployeeMap = {}
      this.batchResult = emptyBatchResult()
      this.batchVisible = true
      this.$nextTick(() => {
        if (this.$refs.batchForm) this.$refs.batchForm.clearValidate()
        this.loadBatchEmployees()
      })
    },
    async loadBatchEmployees() {
      this.batchLoading = true
      try {
        const res = await employeeList({
          page: this.batchPage,
          size: this.batchSize,
          username: this.batchQuery.username,
          real_name: this.batchQuery.real_name,
          status: 1
        })
        this.batchEmployees = res.data.list || []
        this.batchTotal = res.data.total || 0
        this.$nextTick(this.syncBatchPageSelection)
      } finally {
        this.batchLoading = false
      }
    },
    syncBatchPageSelection() {
      if (!this.$refs.batchTable) return
      this.syncingBatchSelection = true
      this.$refs.batchTable.clearSelection()
      this.batchEmployees.forEach(row => {
        if (this.selectedEmployeeMap[row.id]) {
          this.$refs.batchTable.toggleRowSelection(row, true)
        }
      })
      this.$nextTick(() => {
        this.syncingBatchSelection = false
      })
    },
    handleBatchSelectionChange(selection) {
      if (this.syncingBatchSelection) return
      const currentPageIds = new Set(this.batchEmployees.map(row => row.id))
      const selectedOnPage = selection.filter(row => currentPageIds.has(row.id))
      const selectedOutsidePageCount = Object.keys(this.selectedEmployeeMap)
        .filter(id => !currentPageIds.has(Number(id))).length
      const availableSlots = Math.max(0, 200 - selectedOutsidePageCount)
      const acceptedRows = selectedOnPage.slice(0, availableSlots)
      const acceptedIds = new Set(acceptedRows.map(row => row.id))

      this.batchEmployees.forEach(row => {
        if (acceptedIds.has(row.id)) {
          this.$set(this.selectedEmployeeMap, row.id, row)
        } else if (this.selectedEmployeeMap[row.id]) {
          this.$delete(this.selectedEmployeeMap, row.id)
        }
      })

      if (selectedOnPage.length > acceptedRows.length) {
        this.$message.warning('A batch can contain at most 200 employees')
        this.$nextTick(this.syncBatchPageSelection)
      }
    },
    clearBatchSelection() {
      this.selectedEmployeeMap = {}
      this.syncBatchPageSelection()
    },
    handleBatchSearch() {
      this.batchPage = 1
      this.loadBatchEmployees()
    },
    resetBatchSearch() {
      this.batchQuery = { username: '', real_name: '' }
      this.batchPage = 1
      this.loadBatchEmployees()
    },
    handleBatchSizeChange(val) {
      this.batchSize = val
      this.batchPage = 1
      this.loadBatchEmployees()
    },
    submitBatchAllocation() {
      if (this.selectedCount === 0) {
        this.$message.warning('Please select at least one employee')
        return
      }
      this.$refs.batchForm.validate(valid => {
        if (!valid) return
        const employeeIds = this.selectedEmployees.map(item => item.id).sort((a, b) => a - b)
        const employeeLookup = { ...this.selectedEmployeeMap }
        const points = Number(this.batchForm.points)
        const totalPoints = employeeIds.length * points
        this.$confirm(
          `Allocate ${this.formatNumber(points)} credits to ${employeeIds.length} employees (${this.formatNumber(totalPoints)} credits total)?`,
          'Confirm Batch Allocation',
          {
            type: 'warning',
            confirmButtonText: 'Allocate'
          }
        ).then(async() => {
          this.batchSubmitLoading = true
          try {
            const res = await pointsBatchAdd({
              employee_ids: employeeIds,
              points,
              remark: this.batchForm.remark
            })
            const data = res.data || emptyBatchResult()
            this.batchResult = {
              processed_count: data.processed_count || 0,
              total_points: data.total_points || 0,
              list: (data.list || []).map(item => ({
                ...item,
                username: employeeLookup[item.employee_id] ? employeeLookup[item.employee_id].username : '',
                real_name: employeeLookup[item.employee_id] ? employeeLookup[item.employee_id].real_name : ''
              }))
            }
            const shouldRefreshRecords = employeeIds.includes(Number(this.form.employee_id))
            this.batchVisible = false
            this.batchResultVisible = true
            this.$message.success(`Allocated credits to ${this.batchResult.processed_count} employees`)
            if (shouldRefreshRecords) this.loadEmployeeRecords()
          } finally {
            this.batchSubmitLoading = false
          }
        }).catch(() => {})
      })
    },
    resetBatchDialog() {
      this.batchEmployees = []
      this.batchTotal = 0
      this.batchPage = 1
      this.selectedEmployeeMap = {}
      this.syncingBatchSelection = false
      this.batchForm = { points: 1, remark: '' }
    },
    formatNumber(value) {
      return Number(value || 0).toLocaleString('en-US')
    }
  }
}
</script>

<style scoped>
.page-container {
  padding: 24px;
}
.operation-card {
  margin-bottom: 16px;
}
.card-header,
.batch-toolbar,
.result-summary {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}
.batch-toolbar {
  min-height: 36px;
  margin-bottom: 12px;
  color: #606266;
}
.batch-toolbar strong,
.result-summary strong {
  color: #303133;
}
.batch-form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}
.total-points {
  font-size: 18px;
  font-weight: 600;
  color: #409eff;
}
.result-summary {
  margin-bottom: 16px;
  padding: 12px 16px;
  background: #f5f7fa;
  border: 1px solid #e4e7ed;
}
.pagination {
  margin-top: 16px;
  text-align: right;
}
@media (max-width: 760px) {
  .page-container {
    padding: 16px;
  }
  .batch-form-grid {
    grid-template-columns: minmax(0, 1fr);
    gap: 0;
  }
  .batch-toolbar,
  .result-summary {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>

<style>
@media (max-width: 960px) {
  .batch-points-dialog,
  .batch-result-dialog {
    width: calc(100% - 32px) !important;
    margin-top: 3vh !important;
  }
  .batch-points-dialog .pagination {
    overflow-x: auto;
    white-space: nowrap;
  }
}
</style>
