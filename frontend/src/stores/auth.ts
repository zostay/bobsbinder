import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../services/api'
import type { User, AuthResponse } from '../types'
import router from '../router'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'))
  const user = ref<User | null>(null)

  const isAuthenticated = computed(() => !!token.value)

  async function login(email: string, password: string) {
    const { data } = await api.post<AuthResponse>('/auth/login', { email, password })
    token.value = data.token
    user.value = data.user
    localStorage.setItem('token', data.token)
    router.push({ name: 'home' })
  }

  async function register(email: string, password: string, name: string) {
    const { data } = await api.post<AuthResponse>('/auth/register', { email, password, name })
    token.value = data.token
    user.value = data.user
    localStorage.setItem('token', data.token)
    router.push({ name: 'home' })
  }

  function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem('token')
    router.push({ name: 'login' })
  }

  return { token, user, isAuthenticated, login, register, logout }
})
