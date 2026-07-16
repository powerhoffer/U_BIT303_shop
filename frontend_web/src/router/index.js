import Vue from 'vue'
import Router from 'vue-router'
import { getToken } from '@/utils/auth'

Vue.use(Router)

const router = new Router({
  mode: 'history',
  routes: [
    { path: '/', redirect: '/products' },
    { path: '/login', name: 'Login', component: () => import('@/views/Login.vue') },
    { path: '/register', name: 'Register', component: () => import('@/views/Register.vue') },
    { path: '/products', name: 'Products', component: () => import('@/views/Products.vue') },
    { path: '/products/:id', name: 'ProductDetail', component: () => import('@/views/ProductDetail.vue') },
    { path: '/cart', name: 'Cart', component: () => import('@/views/Cart.vue'), meta: { auth: true } },
    { path: '/orders', name: 'Orders', component: () => import('@/views/Orders.vue'), meta: { auth: true } },
    { path: '/orders/:id', name: 'OrderDetail', component: () => import('@/views/OrderDetail.vue'), meta: { auth: true } },
    { path: '/credits', name: 'Credits', component: () => import('@/views/Credits.vue'), meta: { auth: true } },
    { path: '/profile/password', name: 'AccountSecurity', component: () => import('@/views/AccountSecurity.vue'), meta: { auth: true } },
    { path: '*', redirect: '/products' }
  ],
  scrollBehavior() {
    return { x: 0, y: 0 }
  }
})

router.beforeEach((to, from, next) => {
  if (to.meta.auth && !getToken()) {
    next(`/login?redirect=${encodeURIComponent(to.fullPath)}`)
    return
  }
  next()
})

export default router
