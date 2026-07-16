<template>
  <div class="page-container">
    <el-card shadow="never" class="password-card">
      <div slot="header">Change Password</div>
      <el-form ref="form" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="Current Password" prop="old_password">
          <el-input v-model="form.old_password" show-password />
        </el-form-item>
        <el-form-item label="New Password" prop="new_password">
          <el-input v-model="form.new_password" show-password />
        </el-form-item>
        <el-form-item label="Confirm Password" prop="confirm_password">
          <el-input v-model="form.confirm_password" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="submit">Save</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
import { updatePassword } from '@/api/admin'

export default {
  name: 'ProfilePassword',
  data() {
    const validateConfirm = (rule, value, callback) => {
      if (value !== this.form.new_password) {
        callback(new Error('The two new passwords do not match'))
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
        old_password: [{ required: true, message: 'Please enter current password', trigger: 'blur' }],
        new_password: [{ required: true, min: 6, message: 'New password must be at least 6 characters', trigger: 'blur' }],
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
          this.$message.success('Password changed. Please log in again')
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
