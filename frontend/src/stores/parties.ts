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

  return { parties, loading, selfParty, fetchParties }
})
