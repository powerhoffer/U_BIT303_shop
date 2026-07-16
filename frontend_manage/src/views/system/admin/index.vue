<template>
  <div class="page-container">
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="queryForm" size="mini">
        <el-form-item label="Username">
          <el-input v-model.trim="queryForm.username" clearable placeholder="Admin username" />
        </el-form-item>
        <el-form-item label="Name">
          <el-input v-model.trim="queryForm.real_name" clearable placeholder="Admin name" />
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
        <span>Administrators</span>
        <el-button type="primary" size="mini" icon="el-icon-plus" @click="openCreate">New Admin</el-button>
      </div>
      <el-table v-loading="loading" :data="list" border style="width: 100%">
        <el-table-column align="center" prop="id" label="ID" width="70" />
        <el-table-column prop="username" label="Username" min-width="140" />
        <el-table-column prop="real_name" label="Name" min-width="140" />
        <el-table-column prop="phone" label="Phone" min-width="140" />
        <el-table-column prop="email" label="Email" min-width="180" />
        <el-table-column align="center" label="Type" width="120">
          <template slot-scope="scope">
            <el-tag :type="scope.row.is_super === 1 ? 'danger' : 'info'">
              {{ scope.row.is_super === 1 ? 'Super Admin' : 'Admin' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="Roles" min-width="160">
          <template slot-scope="scope">
            <span v-if="!scope.row.role_ids || !scope.row.role_ids.length" class="muted">None</span>
            <el-tag v-for="roleId in scope.row.role_ids" v-else :key="roleId" size="mini" class="role-tag">
              {{ roleName(roleId) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="center" label="Status" width="100">
          <template slot-scope="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'info'">
              {{ scope.row.status === 1 ? 'Active' : 'Disabled' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="center" label="Actions" width="390" fixed="right">
          <template slot-scope="scope">
            <el-button size="mini" icon="el-icon-edit" @click="openEdit(scope.row)">Edit</el-button>
            <el-button size="mini" icon="el-icon-user" @click="openRoles(scope.row)">Roles</el-button>
            <el-button size="mini" :type="scope.row.status === 1 ? 'warning' : 'success'" @click="toggleStatus(scope.row)">
              {{ scope.row.status === 1 ? 'Disable' : 'Enable' }}
            </el-button>
            <el-button size="mini" type="danger" @click="openReset(scope.row)">Reset Password</el-button>
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

    <el-dialog :title="isEdit ? 'Edit Admin' : 'New Admin'" :visible.sync="formVisible" width="560px">
      <el-form ref="adminForm" :model="form" :rules="formRules" label-width="100px">
        <el-form-item label="Username" prop="username">
          <el-input v-model.trim="form.username" :disabled="isEdit" maxlength="64" />
        </el-form-item>
        <el-form-item v-if="!isEdit" label="Password" prop="password">
          <el-input v-model="form.password" show-password />
        </el-form-item>
        <el-form-item label="Name">
          <el-input v-model.trim="form.real_name" maxlength="64" />
        </el-form-item>
        <el-form-item label="Phone">
          <el-input v-model.trim="form.phone" maxlength="20" />
        </el-form-item>
        <el-form-item label="Email" prop="email">
          <el-input v-model.trim="form.email" maxlength="128" />
        </el-form-item>
        <el-form-item label="Super Admin">
          <el-switch v-model="form.is_super" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item v-if="!isEdit" label="Roles">
          <el-select v-model="form.role_ids" multiple clearable placeholder="Select roles" style="width: 100%">
            <el-option v-for="role in roleOptions" :key="role.id" :label="role.name" :value="role.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer">
        <el-button @click="formVisible = false">Cancel</el-button>
        <el-button type="primary" :loading="submitLoading" @click="submitAdmin">Save</el-button>
      </div>
    </el-dialog>

    <el-dialog title="Assign Roles" :visible.sync="rolesVisible" width="500px">
      <el-form label-width="90px">
        <el-form-item label="Admin">
          <span>{{ selectedAdmin.username || '-' }}</span>
        </el-form-item>
        <el-form-item label="Roles">
          <el-select v-model="selectedRoleIds" multiple clearable placeholder="Select roles" style="width: 100%">
            <el-option v-for="role in roleOptions" :key="role.id" :label="role.name" :value="role.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer">
        <el-button @click="rolesVisible = false">Cancel</el-button>
        <el-button type="primary" :loading="submitLoading" @click="submitRoles">Save</el-button>
      </div>
    </el-dialog>

    <el-dialog title="Reset Password" :visible.sync="resetVisible" width="420px">
      <el-form ref="resetForm" :model="resetForm" :rules="resetRules" label-width="110px">
        <el-form-item label="New Password" prop="password">
          <el-input v-model="resetForm.password" show-password />
        </el-form-item>
      </el-form>
      <div slot="footer">
        <el-button @click="resetVisible = false">Cancel</el-button>
        <el-button type="primary" :loading="submitLoading" @click="submitReset">Save</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import {
  adminList,
  adminDetail,
  adminCreate,
  adminUpdate,
  adminStatus,
  adminResetPassword,
  adminRoles
} from '@/api/admin'
import { roleList } from '@/api/role'

const emptyForm = () => ({
  id: undefined,
  username: '',
  password: '',
  real_name: '',
  phone: '',
  email: '',
  is_super: 0,
  role_ids: []
})

export default {
  name: 'AdminList',
  data() {
    return {
      loading: false,
      submitLoading: false,
      list: [],
      roleOptions: [],
      total: 0,
      page: 1,
      size: 10,
      queryForm: { username: '', real_name: '', status: '' },
      formVisible: false,
      rolesVisible: false,
      resetVisible: false,
      isEdit: false,
      form: emptyForm(),
      selectedAdmin: {},
      selectedRoleIds: [],
      resetForm: { id: undefined, password: '' },
      formRules: {
        username: [{ required: true, min: 3, max: 64, message: 'Username must be 3 to 64 characters', trigger: 'blur' }],
        password: [{ required: true, min: 6, max: 64, message: 'Password must be 6 to 64 characters', trigger: 'blur' }],
        email: [{ type: 'email', message: 'Please enter a valid email address', trigger: 'blur' }]
      },
      resetRules: {
        password: [{ required: true, min: 6, max: 64, message: 'Password must be 6 to 64 characters', trigger: 'blur' }]
      }
    }
  },
  created() {
    this.loadRoles()
    this.getList()
  },
  methods: {
    async loadRoles() {
      const res = await roleList({ page: 1, size: 50, status: 1 })
      this.roleOptions = res.data.list || []
    },
    async getList() {
      this.loading = true
      try {
        const res = await adminList({
          page: this.page,
          size: this.size,
          username: this.queryForm.username,
          real_name: this.queryForm.real_name,
          status: this.queryForm.status === '' ? -1 : this.queryForm.status
        })
        this.list = res.data.list || []
        this.total = res.data.total || 0
      } finally {
        this.loading = false
      }
    },
    roleName(id) {
      const role = this.roleOptions.find(item => item.id === id)
      return role ? role.name : `Role #${id}`
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
    handleSizeChange(value) {
      this.size = value
      this.page = 1
      this.getList()
    },
    openCreate() {
      this.isEdit = false
      this.form = emptyForm()
      this.formVisible = true
      this.$nextTick(() => this.$refs.adminForm && this.$refs.adminForm.clearValidate())
    },
    async openEdit(row) {
      const res = await adminDetail({ id: row.id })
      this.isEdit = true
      this.form = { ...emptyForm(), ...res.data.admin, password: '' }
      this.formVisible = true
      this.$nextTick(() => this.$refs.adminForm && this.$refs.adminForm.clearValidate())
    },
    submitAdmin() {
      this.$refs.adminForm.validate(async valid => {
        if (!valid) return
        this.submitLoading = true
        try {
          if (this.isEdit) {
            await adminUpdate({
              id: this.form.id,
              real_name: this.form.real_name,
              phone: this.form.phone,
              email: this.form.email,
              is_super: this.form.is_super,
              role_ids: this.form.role_ids || []
            })
          } else {
            await adminCreate(this.form)
          }
          this.$message.success('Admin saved successfully')
          this.formVisible = false
          this.getList()
        } finally {
          this.submitLoading = false
        }
      })
    },
    async openRoles(row) {
      const res = await adminDetail({ id: row.id })
      this.selectedAdmin = res.data.admin || row
      this.selectedRoleIds = [...(this.selectedAdmin.role_ids || [])]
      this.rolesVisible = true
    },
    async submitRoles() {
      this.submitLoading = true
      try {
        await adminRoles({ id: this.selectedAdmin.id, role_ids: this.selectedRoleIds })
        this.$message.success('Roles assigned successfully')
        this.rolesVisible = false
        this.getList()
      } finally {
        this.submitLoading = false
      }
    },
    toggleStatus(row) {
      const nextStatus = row.status === 1 ? 0 : 1
      this.$confirm(`Confirm ${nextStatus === 1 ? 'enabling' : 'disabling'} admin ${row.username}?`, 'Warning', {
        type: 'warning'
      }).then(async() => {
        await adminStatus({ id: row.id, status: nextStatus })
        this.$message.success('Status updated successfully')
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
          await adminResetPassword(this.resetForm)
          this.$message.success('Password reset successfully')
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
.page-container { padding: 24px; }
.filter-card { margin-bottom: 16px; }
.card-header { display: flex; align-items: center; justify-content: space-between; }
.pagination { margin-top: 16px; text-align: right; }
.role-tag { margin: 2px 4px 2px 0; }
.muted { color: #909399; }
</style>
