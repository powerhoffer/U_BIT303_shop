<template>
  <div class="page-container">
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="queryForm" size="mini">
        <el-form-item label="Role Name">
          <el-input v-model.trim="queryForm.name" clearable placeholder="Role name" />
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
        <span>Roles</span>
        <el-button type="primary" size="mini" icon="el-icon-plus" @click="openCreate">New Role</el-button>
      </div>
      <el-table v-loading="loading" :data="list" border style="width: 100%">
        <el-table-column align="center" prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="Role Name" min-width="180" />
        <el-table-column prop="description" label="Description" min-width="240" show-overflow-tooltip />
        <el-table-column align="center" label="Permissions" width="130">
          <template slot-scope="scope">{{ (scope.row.permission_ids || []).length }}</template>
        </el-table-column>
        <el-table-column align="center" label="Status" width="110">
          <template slot-scope="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'info'">
              {{ scope.row.status === 1 ? 'Active' : 'Disabled' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="center" label="Actions" width="310" fixed="right">
          <template slot-scope="scope">
            <el-button size="mini" icon="el-icon-edit" @click="openEdit(scope.row)">Edit</el-button>
            <el-button size="mini" icon="el-icon-lock" @click="openPermissions(scope.row)">Permissions</el-button>
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
        :page-sizes="[10, 20, 50]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="getList"
      />
    </el-card>

    <el-dialog :title="isEdit ? 'Edit Role' : 'New Role'" :visible.sync="formVisible" width="520px">
      <el-form ref="roleForm" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="Role Name" prop="name">
          <el-input v-model.trim="form.name" maxlength="64" show-word-limit />
        </el-form-item>
        <el-form-item label="Description">
          <el-input v-model.trim="form.description" type="textarea" :rows="4" maxlength="255" show-word-limit />
        </el-form-item>
      </el-form>
      <div slot="footer">
        <el-button @click="formVisible = false">Cancel</el-button>
        <el-button type="primary" :loading="submitLoading" @click="submitRole">Save</el-button>
      </div>
    </el-dialog>

    <el-dialog title="Assign Permissions" :visible.sync="permissionsVisible" width="760px">
      <div class="selected-role">Role: <strong>{{ selectedRole.name || '-' }}</strong></div>
      <el-checkbox-group v-model="selectedPermissionIds" class="permission-groups">
        <div v-for="group in permissionGroups" :key="group.name" class="permission-group">
          <div class="group-title">{{ group.name }}</div>
          <el-checkbox v-for="permission in group.items" :key="permission.id" :label="permission.id">
            <el-tag size="mini" :type="methodTag(permission.method)">{{ permission.method }}</el-tag>
            <span class="permission-label">{{ permission.path }} · {{ permission.name }}</span>
          </el-checkbox>
        </div>
      </el-checkbox-group>
      <div slot="footer">
        <el-button @click="permissionsVisible = false">Cancel</el-button>
        <el-button type="primary" :loading="submitLoading" @click="submitPermissions">Save</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { roleList, roleDetail, roleCreate, roleUpdate, roleStatus, rolePermissions } from '@/api/role'
import { permissionList } from '@/api/permission'

const emptyForm = () => ({ id: undefined, name: '', description: '' })

export default {
  name: 'RoleList',
  data() {
    return {
      loading: false,
      submitLoading: false,
      list: [],
      permissions: [],
      total: 0,
      page: 1,
      size: 10,
      queryForm: { name: '', status: '' },
      formVisible: false,
      permissionsVisible: false,
      isEdit: false,
      form: emptyForm(),
      selectedRole: {},
      selectedPermissionIds: [],
      rules: {
        name: [{ required: true, min: 2, max: 64, message: 'Role name must be 2 to 64 characters', trigger: 'blur' }]
      }
    }
  },
  computed: {
    permissionGroups() {
      const groups = {}
      this.permissions.forEach(permission => {
        const name = permission.group_name || 'Other'
        if (!groups[name]) groups[name] = []
        groups[name].push(permission)
      })
      return Object.keys(groups).sort().map(name => ({ name, items: groups[name] }))
    }
  },
  created() {
    this.loadPermissions()
    this.getList()
  },
  methods: {
    async loadPermissions() {
      const res = await permissionList({ page: 1, size: 100, status: 1 })
      this.permissions = res.data.list || []
    },
    async getList() {
      this.loading = true
      try {
        const res = await roleList({
          page: this.page,
          size: this.size,
          name: this.queryForm.name,
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
      this.queryForm = { name: '', status: '' }
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
      this.$nextTick(() => this.$refs.roleForm && this.$refs.roleForm.clearValidate())
    },
    async openEdit(row) {
      const res = await roleDetail({ id: row.id })
      this.isEdit = true
      this.form = { ...emptyForm(), ...res.data.role }
      this.formVisible = true
      this.$nextTick(() => this.$refs.roleForm && this.$refs.roleForm.clearValidate())
    },
    submitRole() {
      this.$refs.roleForm.validate(async valid => {
        if (!valid) return
        this.submitLoading = true
        try {
          if (this.isEdit) await roleUpdate(this.form)
          else await roleCreate(this.form)
          this.$message.success('Role saved successfully')
          this.formVisible = false
          this.getList()
        } finally {
          this.submitLoading = false
        }
      })
    },
    async openPermissions(row) {
      const res = await roleDetail({ id: row.id })
      this.selectedRole = res.data.role || row
      this.selectedPermissionIds = [...(this.selectedRole.permission_ids || [])]
      this.permissionsVisible = true
    },
    async submitPermissions() {
      this.submitLoading = true
      try {
        await rolePermissions({ id: this.selectedRole.id, permission_ids: this.selectedPermissionIds })
        this.$message.success('Permissions assigned successfully')
        this.permissionsVisible = false
        this.getList()
      } finally {
        this.submitLoading = false
      }
    },
    toggleStatus(row) {
      const nextStatus = row.status === 1 ? 0 : 1
      this.$confirm(`Confirm ${nextStatus === 1 ? 'enabling' : 'disabling'} role ${row.name}?`, 'Warning', {
        type: 'warning'
      }).then(async() => {
        await roleStatus({ id: row.id, status: nextStatus })
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
.selected-role { margin-bottom: 18px; color: #606266; }
.permission-groups { max-height: 480px; overflow-y: auto; }
.permission-group { padding: 14px 0; border-top: 1px solid #ebeef5; }
.group-title { margin-bottom: 10px; font-weight: 600; color: #303133; }
.permission-group .el-checkbox { display: flex; align-items: center; margin: 8px 0; }
.permission-label { margin-left: 8px; overflow-wrap: anywhere; }
</style>
