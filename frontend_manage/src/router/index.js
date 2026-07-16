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
        meta: { title: 'Dashboard', icon: 'dashboard', affix: true }
      }
    ]
  },
  {
    path: '/employee',
    component: Layout,
    redirect: '/employee/list',
    meta: { title: 'Employees', icon: 'people' },
    children: [
      {
        path: 'list',
        component: () => import('@/views/employee/index'),
        name: 'EmployeeList',
        meta: { title: 'Employee List', icon: 'people' }
      }
    ]
  },
  {
    path: '/points',
    component: Layout,
    redirect: '/points/manage',
    meta: { title: 'Credits', icon: 'money' },
    children: [
      {
        path: 'manage',
        component: () => import('@/views/points/ManagePoints'),
        name: 'ManagePoints',
        meta: { title: 'Credit Operations', icon: 'edit' }
      }
    ]
  },
  {
    path: '/category',
    component: Layout,
    redirect: '/category/list',
    meta: { title: 'Categories', icon: 'tree-table' },
    children: [
      {
        path: 'list',
        component: () => import('@/views/category/index'),
        name: 'CategoryList',
        meta: { title: 'Category List', icon: 'tree-table' }
      }
    ]
  },
  {
    path: '/goods',
    component: Layout,
    redirect: '/goods/list',
    meta: { title: 'Goods', icon: 'shopping' },
    children: [
      {
        path: 'list',
        component: () => import('@/views/goods/index'),
        name: 'GoodsList',
        meta: { title: 'Goods Management', icon: 'shopping' }
      }
    ]
  },
  {
    path: '/order',
    component: Layout,
    redirect: '/order/list',
    meta: { title: 'Orders', icon: 'clipboard' },
    children: [
      {
        path: 'list',
        component: () => import('@/views/order/index'),
        name: 'OrderList',
        meta: { title: 'Order Management', icon: 'clipboard' }
      }
    ]
  },
  {
    path: '/stock',
    component: Layout,
    redirect: '/stock/list',
    meta: { title: 'Stock', icon: 'table' },
    children: [
      {
        path: 'list',
        component: () => import('@/views/stock/index'),
        name: 'StockList',
        meta: { title: 'Stock Management', icon: 'table' }
      }
    ]
  },
  {
    path: '/system',
    component: Layout,
    redirect: '/system/admin/list',
    meta: { title: 'Access Control', icon: 'lock' },
    children: [
      {
        path: 'admin/list',
        component: () => import('@/views/system/admin/index'),
        name: 'AdminList',
        meta: { title: 'Admins', icon: 'user' }
      },
      {
        path: 'role/list',
        component: () => import('@/views/system/role/index'),
        name: 'RoleList',
        meta: { title: 'Roles', icon: 'peoples' }
      },
      {
        path: 'permission/list',
        component: () => import('@/views/system/permission/index'),
        name: 'PermissionList',
        meta: { title: 'Permissions', icon: 'lock' }
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
        meta: { title: 'Change Password', icon: 'password' }
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
