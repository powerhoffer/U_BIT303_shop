<template>
  <div class="dashboard-container">
    <el-row :gutter="16">
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card shadow="never" class="stat-card">
          <div class="stat-title">当前员工</div>
          <div class="stat-value">{{ employeeName }}</div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card shadow="never" class="stat-card">
          <div class="stat-title">我的积分</div>
          <div class="stat-value">{{ balance }}</div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card shadow="never" class="stat-card">
          <div class="stat-title">管理模块</div>
          <div class="stat-value">3</div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :lg="6">
        <el-card shadow="never" class="stat-card">
          <div class="stat-title">后端服务</div>
          <div class="stat-value">在线</div>
        </el-card>
      </el-col>
    </el-row>

    <el-card shadow="never" class="quick-card">
      <div slot="header">快捷入口</div>
      <el-button type="primary" icon="el-icon-user" @click="$router.push('/employee/list')">员工管理</el-button>
      <el-button type="success" icon="el-icon-coin" @click="$router.push('/points/manage')">积分操作</el-button>
      <el-button type="warning" icon="el-icon-menu" @click="$router.push('/category/list')">商品分类</el-button>
    </el-card>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { pointsBalance } from '@/api/points'

export default {
  name: 'Dashboard',
  data() {
    return {
      balance: 0
    }
  },
  computed: {
    ...mapGetters(['name']),
    employeeName() {
      return this.name || '员工'
    }
  },
  created() {
    this.loadBalance()
  },
  methods: {
    async loadBalance() {
      const res = await pointsBalance()
      this.balance = res.data.balance || 0
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
