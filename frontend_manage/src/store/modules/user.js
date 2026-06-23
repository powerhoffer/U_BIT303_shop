import { login, logout, getInfo } from '@/api/employee'
import { getToken, setToken, removeToken } from '@/utils/auth'
import { resetRouter } from '@/router'

const defaultAvatar = 'https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif'

const state = {
  token: getToken(),
  name: '',
  avatar: defaultAvatar,
  introduction: '',
  roles: []
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
      const employee = response.data.employee || {}
      commit('SET_TOKEN', token)
      commit('SET_NAME', employee.real_name || employee.username || 'Employee')
      commit('SET_AVATAR', defaultAvatar)
      setToken(token)
      localStorage.setItem('employee', JSON.stringify(employee))
      return response
    })
  },

  getInfo({ commit }) {
    return getInfo().then(response => {
      const employee = response.data.employee || {}
      const name = employee.real_name || employee.username || 'Employee'
      const roles = ['*']
      commit('SET_NAME', name)
      commit('SET_AVATAR', defaultAvatar)
      commit('SET_INTRODUCTION', 'YUTANK employee')
      commit('SET_ROLES', roles)
      localStorage.setItem('employee', JSON.stringify(employee))
      return { roles, name, employee }
    })
  },

  logout({ commit, dispatch }) {
    return logout().catch(() => {}).finally(() => {
      commit('SET_TOKEN', '')
      commit('SET_ROLES', [])
      removeToken()
      resetRouter()
      localStorage.removeItem('employee')
      dispatch('tagsView/delAllViews', null, { root: true })
    })
  },

  resetToken({ commit }) {
    commit('SET_TOKEN', '')
    commit('SET_ROLES', [])
    removeToken()
    localStorage.removeItem('employee')
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
