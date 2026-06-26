export function query(params) {
  const search = new URLSearchParams()
  Object.keys(params || {}).forEach(key => {
    const value = params[key]
    if (value !== undefined && value !== null && value !== '') {
      search.append(key, value)
    }
  })
  return search.toString()
}
