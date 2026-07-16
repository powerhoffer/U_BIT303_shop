<template>
  <div class="page-container">
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="queryForm" size="mini">
        <el-form-item label="Name">
          <el-input v-model.trim="queryForm.name" clearable placeholder="Goods name" />
        </el-form-item>
        <el-form-item label="Category">
          <el-select v-model="queryForm.category_id" clearable placeholder="All" style="width: 180px">
            <el-option
              v-for="item in categories"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="Status">
          <el-select v-model="queryForm.status" clearable placeholder="All" style="width: 120px">
            <el-option label="On Shelf" :value="1" />
            <el-option label="Off Shelf" :value="0" />
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
        <span>Goods Management</span>
        <el-button type="primary" size="mini" icon="el-icon-plus" @click="openCreate">New Goods</el-button>
      </div>

      <el-table v-loading="loading" :data="list" border style="width: 100%">
        <el-table-column align="center" prop="id" label="ID" width="80" />
        <el-table-column align="center" label="Image" width="110">
          <template slot-scope="scope">
            <el-image
              v-if="scope.row.image_url"
              class="goods-image"
              :src="scope.row.image_url"
              fit="cover"
              :preview-src-list="[scope.row.image_url]"
            />
            <span v-else class="image-empty">No Image</span>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="Goods Name" min-width="180" />
        <el-table-column prop="category_name" label="Category" min-width="150" />
        <el-table-column align="center" prop="points_price" label="Credits Price" width="130" />
        <el-table-column align="center" prop="stock" label="Stock" width="100" />
        <el-table-column align="center" label="Status" width="110">
          <template slot-scope="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'info'">
              {{ scope.row.status === 1 ? 'On Shelf' : 'Off Shelf' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="center" label="Actions" width="300" fixed="right">
          <template slot-scope="scope">
            <el-button size="mini" icon="el-icon-view" @click="openDetail(scope.row)">Detail</el-button>
            <el-button size="mini" icon="el-icon-edit" @click="openEdit(scope.row)">Edit</el-button>
            <el-button
              size="mini"
              :type="scope.row.status === 1 ? 'warning' : 'success'"
              @click="toggleStatus(scope.row)"
            >
              {{ scope.row.status === 1 ? 'Off Shelf' : 'On Shelf' }}
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

    <el-dialog :title="isEdit ? 'Edit Goods' : 'New Goods'" :visible.sync="dialogVisible" width="620px">
      <el-form ref="goodsForm" :model="form" :rules="rules" label-width="120px">
        <el-form-item label="Category" prop="category_id">
          <el-select v-model="form.category_id" placeholder="Select category" style="width: 100%">
            <el-option
              v-for="item in enabledCategories"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="Name" prop="name">
          <el-input v-model.trim="form.name" maxlength="128" show-word-limit />
        </el-form-item>
        <el-form-item label="Image URL" prop="image_url">
          <el-input v-model.trim="form.image_url" maxlength="255" show-word-limit />
          <div class="image-tools">
            <el-upload
              action="#"
              accept="image/jpeg,image/png,image/gif,image/webp"
              :show-file-list="false"
              :http-request="handleImageUpload"
              :before-upload="beforeImageUpload"
            >
              <el-button size="mini" icon="el-icon-upload" :loading="uploadLoading">Upload Image</el-button>
            </el-upload>
            <el-image v-if="form.image_url" class="form-image-preview" :src="form.image_url" fit="cover" />
          </div>
        </el-form-item>
        <el-form-item label="Credits Price" prop="points_price">
          <el-input-number v-model="form.points_price" :min="1" :precision="0" />
        </el-form-item>
        <el-form-item label="Stock" prop="stock">
          <el-input-number v-model="form.stock" :min="0" :precision="0" />
        </el-form-item>
        <el-form-item label="Description" prop="description">
          <el-input v-model.trim="form.description" type="textarea" :rows="4" maxlength="1000" show-word-limit />
        </el-form-item>
      </el-form>
      <div slot="footer">
        <el-button @click="dialogVisible = false">Cancel</el-button>
        <el-button type="primary" :loading="submitLoading" @click="submitGoods">Confirm</el-button>
      </div>
    </el-dialog>

    <el-dialog title="Goods Detail" :visible.sync="detailVisible" width="620px">
      <div v-if="detail.id" class="detail-list">
        <div class="detail-row">
          <span class="detail-label">ID</span>
          <span class="detail-value">{{ detail.id }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Name</span>
          <span class="detail-value">{{ detail.name }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Category</span>
          <span class="detail-value">{{ detail.category_name }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Credits Price</span>
          <span class="detail-value">{{ detail.points_price }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Stock</span>
          <span class="detail-value">{{ detail.stock }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Status</span>
          <span class="detail-value">{{ detail.status === 1 ? 'On Shelf' : 'Off Shelf' }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Image URL</span>
          <span class="detail-value">{{ detail.image_url || '-' }}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Description</span>
          <span class="detail-value">{{ detail.description || '-' }}</span>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { categoryList } from '@/api/category'
import { uploadGoodsImage } from '@/api/upload'
import {
  goodsList,
  goodsDetail,
  goodsCreate,
  goodsUpdate,
  goodsStatus
} from '@/api/goods'

const emptyForm = () => ({
  id: undefined,
  category_id: undefined,
  name: '',
  image_url: '',
  points_price: 1,
  stock: 0,
  description: ''
})

export default {
  name: 'GoodsList',
  data() {
    return {
      loading: false,
      submitLoading: false,
      uploadLoading: false,
      list: [],
      total: 0,
      page: 1,
      size: 10,
      categories: [],
      queryForm: {
        name: '',
        category_id: '',
        status: ''
      },
      dialogVisible: false,
      detailVisible: false,
      isEdit: false,
      form: emptyForm(),
      detail: {},
      rules: {
        category_id: [{ required: true, message: 'Please select a category', trigger: 'change' }],
        name: [{ required: true, message: 'Please enter goods name', trigger: 'blur' }],
        points_price: [{ required: true, message: 'Credits price must be greater than 0', trigger: 'blur' }]
      }
    }
  },
  computed: {
    enabledCategories() {
      return this.categories.filter(item => item.status === 1)
    }
  },
  created() {
    this.loadCategories()
    this.getList()
  },
  methods: {
    async loadCategories() {
      const res = await categoryList()
      this.categories = res.data.list || []
    },
    async getList() {
      this.loading = true
      try {
        const params = {
          page: this.page,
          size: this.size,
          name: this.queryForm.name,
          category_id: this.queryForm.category_id,
          status: this.queryForm.status === '' ? -1 : this.queryForm.status
        }
        const res = await goodsList(params)
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
      this.queryForm = { name: '', category_id: '', status: '' }
      this.page = 1
      this.getList()
    },
    handleSizeChange(val) {
      this.size = val
      this.page = 1
      this.getList()
    },
    beforeImageUpload(file) {
      if (!file.type.startsWith('image/')) {
        this.$message.error('Please select an image file')
        return false
      }
      if (file.size > 5 * 1024 * 1024) {
        this.$message.error('Image size must not exceed 5 MB')
        return false
      }
      return true
    },
    async handleImageUpload(options) {
      this.uploadLoading = true
      try {
        const res = await uploadGoodsImage(options.file)
        this.form.image_url = res.data.url
        this.$message.success('Image uploaded successfully')
        if (options.onSuccess) options.onSuccess(res)
      } catch (error) {
        if (options.onError) options.onError(error)
      } finally {
        this.uploadLoading = false
      }
    },
    openCreate() {
      this.isEdit = false
      this.form = emptyForm()
      this.dialogVisible = true
      this.$nextTick(() => this.$refs.goodsForm && this.$refs.goodsForm.clearValidate())
    },
    openEdit(row) {
      this.isEdit = true
      this.form = {
        ...emptyForm(),
        id: row.id,
        category_id: row.category_id,
        name: row.name,
        image_url: row.image_url,
        points_price: row.points_price,
        stock: row.stock,
        description: row.description || ''
      }
      this.dialogVisible = true
      this.$nextTick(() => this.$refs.goodsForm && this.$refs.goodsForm.clearValidate())
    },
    async openDetail(row) {
      const res = await goodsDetail({ id: row.id })
      this.detail = res.data.goods || {}
      this.detailVisible = true
    },
    submitGoods() {
      this.$refs.goodsForm.validate(async valid => {
        if (!valid) return
        this.submitLoading = true
        try {
          const payload = { ...this.form }
          if (this.isEdit) {
            await goodsUpdate(payload)
          } else {
            delete payload.id
            await goodsCreate(payload)
          }
          this.$message.success('Saved successfully')
          this.dialogVisible = false
          this.getList()
        } finally {
          this.submitLoading = false
        }
      })
    },
    toggleStatus(row) {
      const nextStatus = row.status === 1 ? 0 : 1
      const action = nextStatus === 1 ? 'put on shelf' : 'take off shelf'
      this.$confirm(`Confirm to ${action} ${row.name}?`, 'Warning', {
        type: 'warning'
      }).then(async() => {
        await goodsStatus({ id: row.id, status: String(nextStatus) })
        this.$message.success('Status updated successfully')
        this.getList()
      }).catch(() => {})
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
.goods-image {
  width: 64px;
  height: 64px;
  border-radius: 4px;
}
.image-empty {
  color: #909399;
  font-size: 12px;
}
.image-tools {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 10px;
}
.form-image-preview {
  width: 56px;
  height: 56px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
}
.detail-list {
  border: 1px solid #ebeef5;
  border-bottom: 0;
}
.detail-row {
  display: flex;
  min-height: 42px;
  border-bottom: 1px solid #ebeef5;
}
.detail-label {
  flex: 0 0 130px;
  padding: 12px;
  color: #606266;
  background: #f5f7fa;
}
.detail-value {
  flex: 1;
  padding: 12px;
  color: #303133;
  word-break: break-all;
}
.pagination {
  margin-top: 16px;
  text-align: right;
}
</style>
