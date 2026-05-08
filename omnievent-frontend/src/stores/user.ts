import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, LoginRequest, RegisterRequest } from '@/models/user'
import { userState } from '@/lib/userstate'
import { userService } from '@/lib/services'

export const useUserStore = defineStore('user', () => {
  const token = ref<string | null>(userState.getToken())
  const user = ref<User | null>(null)

  const isLoggedIn = computed(() => !!token.value)

  async function login(username: string, password: string) {
    const request: LoginRequest = { username, password }
    const response = await userService.login(request)

    token.value = response.token
    user.value = response.user || null
    userState.setToken(response.token)

    return response
  }

  async function register(username: string, email: string, nickname: string, password: string) {
    const request: RegisterRequest = { username, email, nickname, password }
    return await userService.register(request)
  }

  async function fetchProfile() {
    if (!token.value) return null

    const profile = await userService.getProfile()
    user.value = profile
    return profile
  }

  async function updateProfile(nickname: string, email: string) {
    const updated = await userService.updateProfile({ nickname, email })
    user.value = updated
    return updated
  }

  function logout() {
    token.value = null
    user.value = null
    userState.clearToken()
  }

  return {
    token,
    user,
    isLoggedIn,
    login,
    register,
    fetchProfile,
    updateProfile,
    logout
  }
})
