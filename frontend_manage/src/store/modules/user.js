import { login, logout, getInfo } from '@/api/admin'
import { getToken, setToken, removeToken } from '@/utils/auth'
import { resetRouter } from '@/router'

const defaultAvatar = 'https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif'

const state = {
  token: getToken(),
  name: '',
  avatar: defaultAvatar,
  introduction: '',
  roles: [],
  isSuper: 0,
  roleIds: []
}

const mutations = {
  SET_TOKEN: (state, token) => {
    state.token = token
  },
  SET_NAME: (state, name) => {
    state.name = name
  },
  SET_AVATAR: (state, avatar) => {
    state.avatar = avatar
  },
  SET_INTRODUCTION: (state, introduction) => {
    state.introduction = introduction
  },
  SET_ROLES: (state, roles) => {
    state.roles = roles
  },
  SET_IS_SUPER: (state, isSuper) => {
    state.isSuper = isSuper
  },
  SET_ROLE_IDS: (state, roleIds) => {
    state.roleIds = roleIds
  }
}

const actions = {
  login({ commit }, userInfo) {
    const { username, password, remember } = userInfo
    return login({
      username: username.trim(),
      password,
      remember
    }).then(response => {
      const token = response.data.token
      const admin = response.data.admin || {}
      commit('SET_TOKEN', token)
      commit('SET_NAME', admin.real_name || admin.username || 'Administrator')
      commit('SET_AVATAR', defaultAvatar)
      commit('SET_IS_SUPER', admin.is_super || 0)
      commit('SET_ROLE_IDS', admin.role_ids || [])
      setToken(token)
      localStorage.setItem('admin', JSON.stringify(admin))
      return response
    })
  },

  getInfo({ commit }) {
    return getInfo().then(response => {
      const admin = response.data.admin || {}
      const name = admin.real_name || admin.username || 'Administrator'
      const roles = ['*']
      commit('SET_NAME', name)
      commit('SET_AVATAR', defaultAvatar)
      commit('SET_INTRODUCTION', admin.is_super === 1 ? 'Super administrator' : 'Administrator')
      commit('SET_ROLES', roles)
      commit('SET_IS_SUPER', admin.is_super || 0)
      commit('SET_ROLE_IDS', admin.role_ids || [])
      localStorage.setItem('admin', JSON.stringify(admin))
      return { roles, name, admin }
    })
  },

  logout({ commit, dispatch }) {
    return logout().catch(() => {}).finally(() => {
      commit('SET_TOKEN', '')
      commit('SET_ROLES', [])
      commit('SET_IS_SUPER', 0)
      commit('SET_ROLE_IDS', [])
      removeToken()
      resetRouter()
      localStorage.removeItem('admin')
      dispatch('tagsView/delAllViews', null, { root: true })
    })
  },

  resetToken({ commit }) {
    commit('SET_TOKEN', '')
    commit('SET_ROLES', [])
    commit('SET_IS_SUPER', 0)
    commit('SET_ROLE_IDS', [])
    removeToken()
    localStorage.removeItem('admin')
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
