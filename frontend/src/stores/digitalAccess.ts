import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'
import type { DigitalAccess } from '../types'

export const useDigitalAccessStore = defineStore('digitalAccess', () => {
  const items = ref<DigitalAccess[]>([])
  const loading = ref(false)

  async function fetchItems() {
    loading.value = true
    try {
      const { data } = await api.get<DigitalAccess[]>('/digital-access')
      items.value = data
    } finally {
      loading.value = false
    }
  }

  async function createItem(item: Partial<DigitalAccess>) {
    const { data } = await api.post<DigitalAccess>('/digital-access', item)
    items.value.push(data)
    return data
  }

  async function updateItem(id: number, item: Partial<DigitalAccess>) {
    await api.put(`/digital-access/${id}`, item)
    await fetchItems()
  }

  async function deleteItem(id: number) {
    await api.delete(`/digital-access/${id}`)
    items.value = items.value.filter((i) => i.id !== id)
  }

  return { items, loading, fetchItems, createItem, updateItem, deleteItem }
})
