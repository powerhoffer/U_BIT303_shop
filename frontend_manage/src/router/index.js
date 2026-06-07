import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

import Layout from '@/layout'

export const constantRoutes = [
  {
    path: '/login',
    component: () => import('@/views/login/index'),
    hidden: true
  },
  {
    path: '/redirect',
    component: Layout,
    hidden: true,
    children: [
      {
        path: '/redirect/:path(.*)',
        component: () => import('@/views/redirect/index')
      }
    ]
  },
  {
    path: '/404',
    component: () => import('@/views/error-page/404'),
    hidden: true
  }
]

export const asyncRoutes = [
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        component: () => import('@/views/dashboard/index'),
        name: 'Dashboard',
        meta: { title: '工作台', icon: 'dashboard', affix: true }
      }
    ]
  },
  {
    path: '/employee',
    component: Layout,
    redirect: '/employee/list',
    meta: { title: '员工管理', icon: 'people' },
    children: [
      {
        path: 'list',
        component: () => import('@/views/employee/index'),
        name: 'EmployeeList',
        meta: { title: '员工列表', icon: 'people' }
      }
    ]
  },
  {
    path: '/points',
    component: Layout,
    redirect: '/points/my',
    meta: { title: '积分管理', icon: 'money' },
    children: [
      {
        path: 'my',
        component: () => import('@/views/points/MyPoints'),
        name: 'MyPoints',
        meta: { title: '我的积分', icon: 'money' }
      },
      {
        path: 'manage',
        component: () => import('@/views/points/ManagePoints'),
        name: 'ManagePoints',
        meta: { title: '积分操作', icon: 'edit' }
      }
    ]
  },
  {
    path: '/category',
    component: Layout,
    redirect: '/category/list',
    meta: { title: '商品分类', icon: 'tree-table' },
    children: [
      {
        path: 'list',
        component: () => import('@/views/category/index'),
        name: 'CategoryList',
        meta: { title: '分类列表', icon: 'tree-table' }
      }
    ]
  },
  {
    path: '/profile',
    component: Layout,
    redirect: '/profile/password',
    hidden: true,
    children: [
      {
        path: 'password',
        component: () => import('@/views/profile/password'),
        name: 'ProfilePassword',
        meta: { title: '修改密码', icon: 'password' }
      }
    ]
  },
  { path: '*', redirect: '/404', hidden: true }
]

const createRouter = () => new Router({
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRoutes
})

const router = createRouter()

export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher
}

export default router
