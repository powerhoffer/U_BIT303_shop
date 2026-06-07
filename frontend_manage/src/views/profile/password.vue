<template>
  <div class="page-container">
    <el-card shadow="never" class="password-card">
      <div slot="header">修改密码</div>
      <el-form ref="form" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="旧密码" prop="old_password">
          <el-input v-model="form.old_password" show-password />
        </el-form-item>
        <el-form-item label="新密码" prop="new_password">
          <el-input v-model="form.new_password" show-password />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirm_password">
          <el-input v-model="form.confirm_password" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="submit">保存</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
import { updatePassword } from '@/api/employee'

export default {
  name: 'ProfilePassword',
  data() {
    const validateConfirm = (rule, value, callback) => {
      if (value !== this.form.new_password) {
        callback(new Error('两次输入的新密码不一致'))
      } else {
        callback()
      }
    }
    return {
      loading: false,
      form: {
        old_password: '',
        new_password: '',
        confirm_password: ''
      },
      rules: {
        old_password: [{ required: true, message: '请输入旧密码', trigger: 'blur' }],
        new_password: [{ required: true, min: 6, message: '新密码至少6位', trigger: 'blur' }],
        confirm_password: [{ required: true, validator: validateConfirm, trigger: 'blur' }]
      }
    }
  },
  methods: {
    submit() {
      this.$refs.form.validate(async valid => {
        if (!valid) return
        this.loading = true
        try {
          await updatePassword({
            old_password: this.form.old_password,
            new_password: this.form.new_password
          })
          this.$message.success('密码已修改，请重新登录')
          await this.$store.dispatch('user/logout')
          this.$router.push('/login')
        } finally {
          this.loading = false
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
.password-card {
  max-width: 560px;
}
</style>
