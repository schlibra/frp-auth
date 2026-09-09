function getToken() {
  return localStorage.getItem("token")
}
function setToken(token: string) {
  localStorage.setItem("token", token)
}
function removeToken() {
  localStorage.removeItem("token")
}
export {
  getToken,
  setToken,
  removeToken,
}
