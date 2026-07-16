<template>
  <div class="dashboard-container">
    <el-row :gutter="16">
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card shadow="never" class="stat-card">
          <div class="stat-title">Current Admin</div>
          <div class="stat-value">{{ adminName }}</div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card shadow="never" class="stat-card">
          <div class="stat-title">Account Type</div>
          <div class="stat-value">{{ isSuper === 1 ? 'Super Admin' : 'Administrator' }}</div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card shadow="never" class="stat-card">
          <div class="stat-title">Management Modules</div>
          <div class="stat-value">8</div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card shadow="never" class="stat-card">
          <div class="stat-title">Backend Service</div>
          <div class="stat-value">Online</div>
        </el-card>
      </el-col>
    </el-row>

    <el-card shadow="never" class="quick-card">
      <div slot="header">Quick Links</div>
      <el-button type="primary" icon="el-icon-user" @click="$router.push('/employee/list')">Employees</el-button>
      <el-button type="success" icon="el-icon-coin" @click="$router.push('/points/manage')">Credit Operations</el-button>
      <el-button type="warning" icon="el-icon-menu" @click="$router.push('/category/list')">Categories</el-button>
      <el-button icon="el-icon-box" @click="$router.push('/stock/list')">Stock</el-button>
    </el-card>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: 'Dashboard',
  computed: {
    ...mapGetters(['name', 'isSuper']),
    adminName() {
      return this.name || 'Administrator'
    }
  }
}
</script>

<style scoped>
.dashboard-container {
  padding: 24px;
}
.stat-card {
  margin-bottom: 16px;
}
.stat-title {
  color: #909399;
  font-size: 14px;
}
.stat-value {
  margin-top: 12px;
  font-size: 26px;
  font-weight: 600;
  color: #303133;
}
.quick-card {
  margin-top: 8px;
}
</style>
