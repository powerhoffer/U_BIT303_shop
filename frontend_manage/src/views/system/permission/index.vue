<template>
  <div class="page-container">
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="queryForm" size="mini">
        <el-form-item label="Name">
          <el-input v-model.trim="queryForm.name" clearable placeholder="Permission name" />
        </el-form-item>
        <el-form-item label="Group">
          <el-input v-model.trim="queryForm.group_name" clearable placeholder="Permission group" />
        </el-form-item>
        <el-form-item label="Method">
          <el-select v-model="queryForm.method" clearable placeholder="All" style="width: 120px">
            <el-option v-for="method in methods" :key="method" :label="method" :value="method" />
          </el-select>
        </el-form-item>
        <el-form-item label="Status">
          <el-select v-model="queryForm.status" clearable placeholder="All" style="width: 120px">
            <el-option label="Active" :value="1" />
            <el-option label="Disabled" :value="0" />
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
        <span>Permissions</span>
        <el-button type="primary" size="mini" icon="el-icon-plus" @click="openCreate">New Permission</el-button>
      </div>
      <el-table v-loading="loading" :data="list" border style="width: 100%">
        <el-table-column align="center" prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="Name" min-width="190" />
        <el-table-column prop="group_name" label="Group" min-width="170" />
        <el-table-column align="center" label="Method" width="100">
          <template slot-scope="scope">
            <el-tag :type="methodTag(scope.row.method)">{{ scope.row.method }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="path" label="API Path" min-width="260" />
        <el-table-column align="center" label="Status" width="110">
          <template slot-scope="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'info'">
              {{ scope.row.status === 1 ? 'Active' : 'Disabled' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="center" label="Actions" width="210" fixed="right">
          <template slot-scope="scope">
            <el-button size="mini" icon="el-icon-edit" @click="openEdit(scope.row)">Edit</el-button>
            <el-button size="mini" :type="scope.row.status === 1 ? 'warning' : 'success'" @click="toggleStatus(scope.row)">
              {{ scope.row.status === 1 ? 'Disable' : 'Enable' }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        class="pagination"
        :current-page.sync="page"
        :page-size="size"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="getList"
      />
    </el-card>

    <el-dialog :title="isEdit ? 'Edit Permission' : 'New Permission'" :visible.sync="formVisible" width="560px">
      <el-form ref="permissionForm" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="Name" prop="name">
          <el-input v-model.trim="form.name" maxlength="128" show-word-limit />
        </el-form-item>
        <el-form-item label="Group">
          <el-input v-model.trim="form.group_name" maxlength="64" show-word-limit />
        </el-form-item>
        <el-form-item label="Method" prop="method">
          <el-select v-model="form.method" placeholder="Select method" style="width: 100%">
            <el-option v-for="method in methods" :key="method" :label="method" :value="method" />
          </el-select>
        </el-form-item>
        <el-form-item label="API Path" prop="path">
          <el-input v-model.trim="form.path" maxlength="255" placeholder="/backend/module/action" />
        </el-form-item>
      </el-form>
      <div slot="footer">
        <el-button @click="formVisible = false">Cancel</el-button>
        <el-button type="primary" :loading="submitLoading" @click="submitPermission">Save</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import {
  permissionList,
  permissionDetail,
  permissionCreate,
  permissionUpdate,
  permissionStatus
} from '@/api/permission'

const methods = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH']
const emptyForm = () => ({ id: undefined, name: '', group_name: '', method: 'GET', path: '' })

export default {
  name: 'PermissionList',
  data() {
    return {
      loading: false,
      submitLoading: false,
      list: [],
      total: 0,
      page: 1,
      size: 10,
      methods,
      queryForm: { name: '', group_name: '', method: '', status: '' },
      formVisible: false,
      isEdit: false,
      form: emptyForm(),
      rules: {
        name: [{ required: true, min: 2, max: 128, message: 'Name must be 2 to 128 characters', trigger: 'blur' }],
        method: [{ required: true, message: 'Please select an HTTP method', trigger: 'change' }],
        path: [
          { required: true, message: 'API path is required', trigger: 'blur' },
          { pattern: /^\//, message: 'API path must start with /', trigger: 'blur' }
        ]
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
        const res = await permissionList({
          page: this.page,
          size: this.size,
          name: this.queryForm.name,
          group_name: this.queryForm.group_name,
          method: this.queryForm.method,
          status: this.queryForm.status === '' ? -1 : this.queryForm.status
        })
        this.list = res.data.list || []
        this.total = res.data.total || 0
      } finally {
        this.loading = false
      }
    },
    methodTag(method) {
      return { GET: 'success', POST: 'primary', PUT: 'warning', DELETE: 'danger', PATCH: 'info' }[method] || 'info'
    },
    handleSearch() {
      this.page = 1
      this.getList()
    },
    handleReset() {
      this.queryForm = { name: '', group_name: '', method: '', status: '' }
      this.page = 1
      this.getList()
    },
    handleSizeChange(value) {
      this.size = value
      this.page = 1
      this.getList()
    },
    openCreate() {
      this.isEdit = false
      this.form = emptyForm()
      this.formVisible = true
      this.$nextTick(() => this.$refs.permissionForm && this.$refs.permissionForm.clearValidate())
    },
    async openEdit(row) {
      const res = await permissionDetail({ id: row.id })
      this.isEdit = true
      this.form = { ...emptyForm(), ...res.data.permission }
      this.formVisible = true
      this.$nextTick(() => this.$refs.permissionForm && this.$refs.permissionForm.clearValidate())
    },
    submitPermission() {
      this.$refs.permissionForm.validate(async valid => {
        if (!valid) return
        this.submitLoading = true
        try {
          if (this.isEdit) await permissionUpdate(this.form)
          else await permissionCreate(this.form)
          this.$message.success('Permission saved successfully')
          this.formVisible = false
          this.getList()
        } finally {
          this.submitLoading = false
        }
      })
    },
    toggleStatus(row) {
      const nextStatus = row.status === 1 ? 0 : 1
      this.$confirm(`Confirm ${nextStatus === 1 ? 'enabling' : 'disabling'} permission ${row.name}?`, 'Warning', {
        type: 'warning'
      }).then(async() => {
        await permissionStatus({ id: row.id, status: nextStatus })
        this.$message.success('Status updated successfully')
        this.getList()
      }).catch(() => {})
    }
  }
}
</script>

<style scoped>
.page-container { padding: 24px; }
.filter-card { margin-bottom: 16px; }
.card-header { display: flex; align-items: center; justify-content: space-between; }
.pagination { margin-top: 16px; text-align: right; }
</style>
