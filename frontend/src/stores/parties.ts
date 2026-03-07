import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../services/api'
import type { Party } from '../types'

export const usePartyStore = defineStore('parties', () => {
  const parties = ref<Party[]>([])
  const loading = ref(false)

  const selfParty = computed(() => parties.value.find((p) => p.relationship === 'self'))

  async function fetchParties() {
    loading.value = true
    try {
      const { data } = await api.get<Party[]>('/parties')
      parties.value = data
    } finally {
      loading.value = false
    }
  }

  async function createParty(party: Partial<Party>) {
    const { data } = await api.post<Party>('/parties', party)
    parties.value.push(data)
    return data
  }

  async function updateParty(id: number, party: Partial<Party>) {
    await api.put(`/parties/${id}`, party)
    await fetchParties()
  }

  async function deleteParty(id: number) {
    await api.delete(`/parties/${id}`)
    parties.value = parties.value.filter((p) => p.id !== id)
  }

  return { parties, loading, selfParty, fetchParties, createParty, updateParty, deleteParty }
})
