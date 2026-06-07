<template>
  <div class="page-container">
    <el-card shadow="never">
      <div slot="header" class="card-header">
        <span>商品分类</span>
        <el-button size="mini" icon="el-icon-refresh" @click="getList">刷新</el-button>
      </div>
      <el-table v-loading="loading" :data="list" border style="width: 100%">
        <el-table-column align="center" prop="id" label="ID" width="100" />
        <el-table-column prop="name" label="分类名称" min-width="180" />
        <el-table-column align="center" prop="sort" label="排序" width="120" />
        <el-table-column align="center" label="状态" width="120">
          <template slot-scope="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'info'">
              {{ scope.row.status === 1 ? '启用' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
      <div class="table-footer">共 {{ total }} 个分类</div>
    </el-card>
  </div>
</template>

<script>
import { categoryList } from '@/api/category'

export default {
  name: 'CategoryList',
  data() {
    return {
      loading: false,
      list: [],
      total: 0
    }
  },
  created() {
    this.getList()
  },
  methods: {
    async getList() {
      this.loading = true
      try {
        const res = await categoryList()
        this.list = res.data.list || []
        this.total = res.data.total || 0
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>
.page-container {
  padding: 24px;
}
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.table-footer {
  padding-top: 16px;
  color: #606266;
  text-align: right;
}
</style>
