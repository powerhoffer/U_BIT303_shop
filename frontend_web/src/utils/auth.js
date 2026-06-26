const TokenKey = 'YUTANK_SHOP_WEB_TOKEN'
const UserKey = 'YUTANK_SHOP_WEB_USER'

export function getToken() {
  return localStorage.getItem(TokenKey)
}

export function setToken(token) {
  localStorage.setItem(TokenKey, token)
}

export function removeToken() {
  localStorage.removeItem(TokenKey)
}

export function getUser() {
  const raw = localStorage.getItem(UserKey)
  if (!raw) return null
  try {
    return JSON.parse(raw)
  } catch (e) {
    return null
  }
}

export function setUser(user) {
  localStorage.setItem(UserKey, JSON.stringify(user || {}))
}

export function removeUser() {
  localStorage.removeItem(UserKey)
}
