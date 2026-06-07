<template>
  <div class="page-container">
    <el-card shadow="never" class="operation-card">
      <div slot="header">积分操作</div>
      <el-form ref="form" :inline="true" :model="form" :rules="rules" size="mini">
        <el-form-item label="员工" prop="employee_id">
          <el-select
            v-model="form.employee_id"
            filterable
            remote
            reserve-keyword
            placeholder="输入账号或姓名搜索"
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
        <el-form-item label="积分" prop="points">
          <el-input-number v-model="form.points" :min="1" :precision="0" />
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model.trim="form.remark" placeholder="操作备注" style="width: 260px" />
        </el-form-item>
        <el-form-item>
          <el-button type="success" :loading="submitLoading" @click="submit('add')">增加</el-button>
          <el-button type="danger" :loading="submitLoading" @click="submit('deduct')">扣除</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card shadow="never">
      <div slot="header" class="card-header">
        <span>员工积分流水</span>
        <el-button size="mini" icon="el-icon-refresh" :disabled="!form.employee_id" @click="loadEmployeeRecords">刷新</el-button>
      </div>
      <el-table v-loading="recordsLoading" :data="records" border>
        <el-table-column align="center" prop="id" label="ID" width="80" />
        <el-table-column align="center" prop="employee_id" label="员工ID" width="100" />
        <el-table-column align="center" label="类型" width="100">
          <template slot-scope="scope">
            <el-tag :type="scope.row.change_type === 1 ? 'success' : 'danger'">
              {{ scope.row.change_type === 1 ? '增加' : '扣除' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="center" prop="points" label="变动积分" width="120" />
        <el-table-column align="center" prop="before_balance" label="变动前" width="120" />
        <el-table-column align="center" prop="after_balance" label="变动后" width="120" />
        <el-table-column prop="remark" label="备注" min-width="180" />
        <el-table-column prop="created_at" label="时间" min-width="180" />
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
  </div>
</template>

<script>
import { employeeList } from '@/api/employee'
import { pointsAdd, pointsDeduct, pointsManageRecords } from '@/api/points'

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
        employee_id: [{ required: true, message: '请选择员工', trigger: 'change' }],
        points: [{ required: true, message: '请输入积分', trigger: 'blur' }]
      }
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
          this.$message.success('积分操作成功')
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
