import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'
import type { InsurancePolicy } from '../types'

export const useInsurancePolicyStore = defineStore('insurancePolicies', () => {
  const policies = ref<InsurancePolicy[]>([])
  const loading = ref(false)

  async function fetchPolicies() {
    loading.value = true
    try {
      const { data } = await api.get<InsurancePolicy[]>('/insurance-policies')
      policies.value = data
    } finally {
      loading.value = false
    }
  }

  async function createPolicy(policy: Partial<InsurancePolicy>) {
    const { data } = await api.post<InsurancePolicy>('/insurance-policies', policy)
    policies.value.push(data)
    return data
  }

  async function updatePolicy(id: number, policy: Partial<InsurancePolicy>) {
    await api.put(`/insurance-policies/${id}`, policy)
    await fetchPolicies()
  }

  async function deletePolicy(id: number) {
    await api.delete(`/insurance-policies/${id}`)
    policies.value = policies.value.filter((p) => p.id !== id)
  }

  return { policies, loading, fetchPolicies, createPolicy, updatePolicy, deletePolicy }
})
