import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'
import type { ServiceAccount } from '../types'

export const useServiceAccountStore = defineStore('serviceAccounts', () => {
  const accounts = ref<ServiceAccount[]>([])
  const loading = ref(false)

  async function fetchAccounts() {
    loading.value = true
    try {
      const { data } = await api.get<ServiceAccount[]>('/service-accounts')
      accounts.value = data
    } finally {
      loading.value = false
    }
  }

  async function createAccount(account: Partial<ServiceAccount>) {
    const { data } = await api.post<ServiceAccount>('/service-accounts', account)
    accounts.value.push(data)
    return data
  }

  async function updateAccount(id: number, account: Partial<ServiceAccount>) {
    await api.put(`/service-accounts/${id}`, account)
    await fetchAccounts()
  }

  async function deleteAccount(id: number) {
    await api.delete(`/service-accounts/${id}`)
    accounts.value = accounts.value.filter((a) => a.id !== id)
  }

  return { accounts, loading, fetchAccounts, createAccount, updateAccount, deleteAccount }
})
