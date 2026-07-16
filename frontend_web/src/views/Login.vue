<template>
  <div class="login-page">
    <div class="login-panel">
      <h1>YUTANK Shop</h1>
      <p>Sign in to use cart, credits, and orders.</p>
      <el-form ref="form" :model="form" :rules="rules" label-position="top" @submit.native.prevent>
        <el-form-item label="Username" prop="username">
          <el-input v-model="form.username" autocomplete="username" />
        </el-form-item>
        <el-form-item label="Password" prop="password">
          <el-input v-model="form.password" type="password" autocomplete="current-password" show-password />
        </el-form-item>
        <el-checkbox v-model="form.remember">Remember me</el-checkbox>
        <el-button class="login-button" type="primary" :loading="loading" @click="submit">Login</el-button>
        <div class="register-link">New to YUTANK Shop? <router-link to="/register">Create an account</router-link></div>
      </el-form>
    </div>
  </div>
</template>

<script>
import { login } from '@/api/auth'
import { setToken, setUser } from '@/utils/auth'

export default {
  name: 'Login',
  data() {
    return {
      loading: false,
      form: {
        username: '',
        password: '',
        remember: false
      },
      rules: {
        username: [{ required: true, message: 'Username is required', trigger: 'blur' }],
        password: [{ required: true, message: 'Password is required', trigger: 'blur' }]
      }
    }
  },
  created() {
    if (this.$route.query.username) this.form.username = this.$route.query.username
  },
  methods: {
    submit() {
      this.$refs.form.validate(async valid => {
        if (!valid) return
        this.loading = true
        try {
          const res = await login(this.form)
          setToken(res.data.token)
          setUser(res.data.employee)
          this.$emit('login-change')
          const redirect = this.$route.query.redirect || '/products'
          this.$router.replace(redirect)
        } finally {
          this.loading = false
        }
      })
    }
  }
}
</script>

<style scoped>
.login-page {
  display: flex;
  justify-content: center;
  padding-top: 72px;
}

.login-panel {
  width: 420px;
  max-width: 100%;
  background: #ffffff;
  border: 1px solid #e5eaf2;
  border-radius: 8px;
  padding: 34px;
}

.login-panel h1 {
  margin: 0 0 8px;
  font-size: 32px;
}

.login-panel p {
  margin: 0 0 28px;
  color: #7b8794;
}

.login-button {
  width: 100%;
  margin-top: 24px;
}

.register-link {
  margin-top: 18px;
  text-align: center;
  color: #7b8794;
}

.register-link a {
  color: #2563eb;
  font-weight: 600;
}
</style>
