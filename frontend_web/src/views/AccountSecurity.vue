<template>
  <div>
    <h1 class="page-title">Account Security</h1>
    <div class="panel security-panel">
      <h2>Change Password</h2>
      <p class="muted">You will need to sign in again after changing your password.</p>
      <el-form ref="form" :model="form" :rules="rules" label-position="top" @submit.native.prevent>
        <el-form-item label="Current Password" prop="old_password">
          <el-input v-model="form.old_password" type="password" autocomplete="current-password" show-password />
        </el-form-item>
        <el-form-item label="New Password" prop="new_password">
          <el-input v-model="form.new_password" type="password" autocomplete="new-password" show-password />
        </el-form-item>
        <el-form-item label="Confirm New Password" prop="confirm_password">
          <el-input v-model="form.confirm_password" type="password" autocomplete="new-password" show-password />
        </el-form-item>
        <el-button type="primary" :loading="loading" @click="submit">Change Password</el-button>
      </el-form>
    </div>
  </div>
</template>

<script>
import { updatePassword } from '@/api/auth'
import { removeToken, removeUser } from '@/utils/auth'

export default {
  name: 'AccountSecurity',
  data() {
    const confirmPassword = (rule, value, callback) => {
      if (value !== this.form.new_password) callback(new Error('The new passwords do not match'))
      else callback()
    }
    return {
      loading: false,
      form: { old_password: '', new_password: '', confirm_password: '' },
      rules: {
        old_password: [{ required: true, message: 'Current password is required', trigger: 'blur' }],
        new_password: [{ required: true, min: 6, max: 64, message: 'New password must be 6 to 64 characters', trigger: 'blur' }],
        confirm_password: [{ required: true, validator: confirmPassword, trigger: 'blur' }]
      }
    }
  },
  methods: {
    submit() {
      this.$refs.form.validate(async valid => {
        if (!valid) return
        this.loading = true
        try {
          await updatePassword({ old_password: this.form.old_password, new_password: this.form.new_password })
          removeToken()
          removeUser()
          this.$emit('login-change')
          this.$message.success('Password changed. Please sign in again.')
          this.$router.replace('/login')
        } finally {
          this.loading = false
        }
      })
    }
  }
}
</script>

<style scoped>
.security-panel {
  width: 520px;
  max-width: 100%;
  box-sizing: border-box;
}
.security-panel h2 {
  margin: 0 0 8px;
  font-size: 22px;
}
.security-panel p {
  margin: 0 0 24px;
}
</style>
