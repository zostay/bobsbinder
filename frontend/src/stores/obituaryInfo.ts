import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'
import type { PartyObituaryInfo } from '../types'

export const useObituaryInfoStore = defineStore('obituaryInfo', () => {
  const items = ref<PartyObituaryInfo[]>([])
  const loading = ref(false)

  async function fetchItems(partyId: number) {
    loading.value = true
    try {
      const { data } = await api.get<PartyObituaryInfo[]>(`/parties/${partyId}/obituary-info`)
      items.value = data
    } finally {
      loading.value = false
    }
  }

  async function createItem(partyId: number, item: Partial<PartyObituaryInfo>) {
    const { data } = await api.post<PartyObituaryInfo>(`/parties/${partyId}/obituary-info`, item)
    items.value.push(data)
    return data
  }

  async function updateItem(partyId: number, id: number, item: Partial<PartyObituaryInfo>) {
    await api.put(`/parties/${partyId}/obituary-info/${id}`, item)
    await fetchItems(partyId)
  }

  async function deleteItem(partyId: number, id: number) {
    await api.delete(`/parties/${partyId}/obituary-info/${id}`)
    items.value = items.value.filter((i) => i.id !== id)
  }

  return { items, loading, fetchItems, createItem, updateItem, deleteItem }
})
