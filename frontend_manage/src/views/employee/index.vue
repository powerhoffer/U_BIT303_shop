<template>
  <div class="page-container">
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="queryForm" size="mini">
        <el-form-item label="账号">
          <el-input v-model.trim="queryForm.username" clearable placeholder="员工账号" />
        </el-form-item>
        <el-form-item label="姓名">
          <el-input v-model.trim="queryForm.real_name" clearable placeholder="员工姓名" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryForm.status" clearable placeholder="全部" style="width: 120px">
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="el-icon-search" @click="handleSearch">搜索</el-button>
          <el-button icon="el-icon-refresh-left" @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card shadow="never">
      <div slot="header" class="card-header">
        <span>员工列表</span>
        <el-button type="primary" size="mini" icon="el-icon-plus" @click="openCreate">新增员工</el-button>
      </div>
      <el-table v-loading="loading" :data="list" border style="width: 100%">
        <el-table-column align="center" prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="账号" min-width="150" />
        <el-table-column prop="real_name" label="姓名" min-width="140" />
        <el-table-column prop="phone" label="手机号" min-width="150" />
        <el-table-column prop="email" label="邮箱" min-width="190" />
        <el-table-column align="center" label="状态" width="100">
          <template slot-scope="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'info'">
              {{ scope.row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="center" label="操作" width="260" fixed="right">
          <template slot-scope="scope">
            <el-button size="mini" icon="el-icon-edit" @click="openEdit(scope.row)">编辑</el-button>
            <el-button size="mini" :type="scope.row.status === 1 ? 'warning' : 'success'" @click="toggleStatus(scope.row)">
              {{ scope.row.status === 1 ? '禁用' : '启用' }}
            </el-button>
            <el-button size="mini" type="danger" @click="openReset(scope.row)">重置密码</el-button>
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

    <el-dialog :title="isEdit ? '编辑员工' : '新增员工'" :visible.sync="dialogVisible" width="520px">
      <el-form ref="employeeForm" :model="form" :rules="rules" label-width="90px">
        <el-form-item label="账号" prop="username">
          <el-input v-model.trim="form.username" :disabled="isEdit" />
        </el-form-item>
        <el-form-item v-if="!isEdit" label="密码" prop="password">
          <el-input v-model="form.password" show-password />
        </el-form-item>
        <el-form-item label="姓名" prop="real_name">
          <el-input v-model.trim="form.real_name" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model.trim="form.phone" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model.trim="form.email" />
        </el-form-item>
      </el-form>
      <div slot="footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="submitEmployee">确定</el-button>
      </div>
    </el-dialog>

    <el-dialog title="重置密码" :visible.sync="resetVisible" width="420px">
      <el-form ref="resetForm" :model="resetForm" :rules="resetRules" label-width="90px">
        <el-form-item label="新密码" prop="password">
          <el-input v-model="resetForm.password" show-password />
        </el-form-item>
      </el-form>
      <div slot="footer">
        <el-button @click="resetVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="submitReset">确定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import {
  employeeList,
  employeeCreate,
  employeeUpdate,
  employeeStatus,
  employeeResetPassword
} from '@/api/employee'

const emptyForm = () => ({
  id: undefined,
  username: '',
  password: '',
  real_name: '',
  phone: '',
  email: ''
})

export default {
  name: 'EmployeeList',
  data() {
    return {
      loading: false,
      submitLoading: false,
      list: [],
      total: 0,
      page: 1,
      size: 10,
      queryForm: {
        username: '',
        real_name: '',
        status: ''
      },
      dialogVisible: false,
      resetVisible: false,
      isEdit: false,
      form: emptyForm(),
      resetForm: {
        id: undefined,
        password: ''
      },
      rules: {
        username: [{ required: true, message: '请输入账号', trigger: 'blur' }],
        password: [{ required: true, min: 6, message: '密码至少6位', trigger: 'blur' }]
      },
      resetRules: {
        password: [{ required: true, min: 6, message: '密码至少6位', trigger: 'blur' }]
      }
    }
  },
  created() {
    this.getList()
  },
  methods: {
    async getList() {
      this.loading = true
      try {
        const params = {
          page: this.page,
          size: this.size,
          username: this.queryForm.username,
          real_name: this.queryForm.real_name,
          status: this.queryForm.status === '' ? -1 : this.queryForm.status
        }
        const res = await employeeList(params)
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
      this.queryForm = { username: '', real_name: '', status: '' }
      this.page = 1
      this.getList()
    },
    handleSizeChange(val) {
      this.size = val
      this.page = 1
      this.getList()
    },
    openCreate() {
      this.isEdit = false
      this.form = emptyForm()
      this.dialogVisible = true
      this.$nextTick(() => this.$refs.employeeForm && this.$refs.employeeForm.clearValidate())
    },
    openEdit(row) {
      this.isEdit = true
      this.form = { ...emptyForm(), ...row }
      this.dialogVisible = true
      this.$nextTick(() => this.$refs.employeeForm && this.$refs.employeeForm.clearValidate())
    },
    async submitEmployee() {
      this.$refs.employeeForm.validate(async valid => {
        if (!valid) return
        this.submitLoading = true
        try {
          if (this.isEdit) {
            await employeeUpdate({
              id: this.form.id,
              real_name: this.form.real_name,
              phone: this.form.phone,
              email: this.form.email
            })
          } else {
            await employeeCreate(this.form)
          }
          this.$message.success('保存成功')
          this.dialogVisible = false
          this.getList()
        } finally {
          this.submitLoading = false
        }
      })
    },
    toggleStatus(row) {
      const nextStatus = row.status === 1 ? 0 : 1
      this.$confirm(`确认${nextStatus === 1 ? '启用' : '禁用'}员工 ${row.username}？`, '提示', {
        type: 'warning'
      }).then(async() => {
        await employeeStatus({ id: row.id, status: String(nextStatus) })
        this.$message.success('状态已更新')
        this.getList()
      }).catch(() => {})
    },
    openReset(row) {
      this.resetForm = { id: row.id, password: '' }
      this.resetVisible = true
      this.$nextTick(() => this.$refs.resetForm && this.$refs.resetForm.clearValidate())
    },
    submitReset() {
      this.$refs.resetForm.validate(async valid => {
        if (!valid) return
        this.submitLoading = true
        try {
          await employeeResetPassword(this.resetForm)
          this.$message.success('密码已重置')
          this.resetVisible = false
        } finally {
          this.submitLoading = false
        }
      })
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
.pagination {
  margin-top: 16px;
  text-align: right;
}
</style>
