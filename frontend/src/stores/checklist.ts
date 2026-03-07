import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'
import type { PartyChecklist } from '../types'

export const useChecklistStore = defineStore('checklist', () => {
  const checklists = ref<PartyChecklist[]>([])
  const loading = ref(false)

  async function fetchAll() {
    loading.value = true
    try {
      const { data } = await api.get<PartyChecklist[]>('/checklist')
      checklists.value = data
    } finally {
      loading.value = false
    }
  }

  async function updateStatus(partyId: number, categoryId: number, status: string) {
    await api.put(`/parties/${partyId}/checklist/${categoryId}`, { status })
    await fetchAll()
  }

  return { checklists, loading, fetchAll, updateStatus }
})
