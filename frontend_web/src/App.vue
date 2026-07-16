<template>
  <div id="app">
    <header class="site-header">
      <div class="brand" @click="$router.push('/products')">YUTANK Shop</div>
      <nav class="nav-links">
        <router-link to="/products">Products</router-link>
        <router-link to="/cart">Cart</router-link>
        <router-link to="/orders">My Orders</router-link>
        <router-link to="/credits">My Credits</router-link>
      </nav>
      <div class="user-actions">
        <template v-if="user">
          <span class="user-name">{{ user.real_name || user.username }}</span>
          <el-button size="small" icon="el-icon-lock" @click="$router.push('/profile/password')">Account Security</el-button>
          <el-button size="small" @click="handleLogout">Logout</el-button>
        </template>
        <template v-else>
          <el-button size="small" @click="$router.push('/register')">Register</el-button>
          <el-button size="small" type="primary" @click="$router.push('/login')">Login</el-button>
        </template>
      </div>
    </header>
    <main class="site-main">
      <router-view @login-change="loadUser" />
    </main>
  </div>
</template>

<script>
import { logout, getInfo } from '@/api/auth'
import { getToken, getUser, removeToken, removeUser, setUser } from '@/utils/auth'

export default {
  name: 'App',
  data() {
    return {
      user: getUser()
    }
  },
  created() {
    this.loadUser()
  },
  methods: {
    async loadUser() {
      if (!getToken()) {
        this.user = null
        return
      }
      try {
        const res = await getInfo()
        this.user = res.data && res.data.employee ? res.data.employee : getUser()
        setUser(this.user)
      } catch (e) {
        this.user = getUser()
      }
    },
    async handleLogout() {
      try {
        await logout()
      } catch (e) {
        // Token cleanup below keeps the storefront responsive even if logout fails.
      }
      removeToken()
      removeUser()
      this.user = null
      this.$router.push('/products')
    }
  }
}
</script>
