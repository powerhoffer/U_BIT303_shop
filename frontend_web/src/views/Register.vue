<template>
  <div class="account-page">
    <div class="account-panel wide">
      <h1>Create Account</h1>
      <p>Register an employee account to redeem goods with Credits.</p>
      <el-form ref="form" :model="form" :rules="rules" label-position="top" @submit.native.prevent>
        <div class="form-grid">
          <el-form-item label="Username" prop="username">
            <el-input v-model.trim="form.username" autocomplete="username" maxlength="64" />
          </el-form-item>
          <el-form-item label="Password" prop="password">
            <el-input v-model="form.password" type="password" autocomplete="new-password" show-password />
          </el-form-item>
          <el-form-item label="Name">
            <el-input v-model.trim="form.real_name" maxlength="64" />
          </el-form-item>
          <el-form-item label="Phone">
            <el-input v-model.trim="form.phone" maxlength="20" />
          </el-form-item>
          <el-form-item label="Email" prop="email" class="full-row">
            <el-input v-model.trim="form.email" maxlength="128" />
          </el-form-item>
        </div>
        <el-button class="account-submit" type="primary" :loading="loading" @click="submit">Create Account</el-button>
        <div class="account-link">
          Already have an account? <router-link to="/login">Sign in</router-link>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script>
import { register } from '@/api/auth'

export default {
  name: 'Register',
  data() {
    return {
      loading: false,
      form: { username: '', password: '', real_name: '', phone: '', email: '' },
      rules: {
        username: [{ required: true, min: 3, max: 64, message: 'Username must be 3 to 64 characters', trigger: 'blur' }],
        password: [{ required: true, min: 6, max: 64, message: 'Password must be 6 to 64 characters', trigger: 'blur' }],
        email: [{ type: 'email', message: 'Please enter a valid email address', trigger: 'blur' }]
      }
    }
  },
  methods: {
    submit() {
      this.$refs.form.validate(async valid => {
        if (!valid) return
        this.loading = true
        try {
          await register(this.form)
          this.$message.success('Account created successfully. Please sign in.')
          this.$router.replace({ path: '/login', query: { username: this.form.username } })
        } finally {
          this.loading = false
        }
      })
    }
  }
}
</script>
