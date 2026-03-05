import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'
import type { FullSurvivorLetter, SurvivorLetterItem } from '../types'

export const useSurvivorLetterStore = defineStore('survivorLetter', () => {
  const letter = ref<FullSurvivorLetter | null>(null)
  const loading = ref(false)

  async function fetchLetter() {
    loading.value = true
    try {
      const { data } = await api.get<FullSurvivorLetter>('/survivor-letter')
      letter.value = data
    } finally {
      loading.value = false
    }
  }

  async function updateBoilerplate(updates: {
    greeting: string
    intro: string
    closing: string
    signature: string
  }) {
    await api.put('/survivor-letter', updates)
    await fetchLetter()
  }

  async function updateSection(
    sectionId: number,
    updates: { title?: string; visible?: boolean; sort_order?: number },
  ) {
    await api.put(`/survivor-letter/sections/${sectionId}`, updates)
    await fetchLetter()
  }

  async function reorderSections(sectionOrders: { id: number; sort_order: number }[]) {
    await api.put('/survivor-letter/sections/reorder', { section_orders: sectionOrders })
    await fetchLetter()
  }

  async function addItem(sectionId: number, content: string, itemType: string = 'numbered') {
    const { data } = await api.post<SurvivorLetterItem>(
      `/survivor-letter/sections/${sectionId}/items`,
      { content, item_type: itemType },
    )
    await fetchLetter()
    return data
  }

  async function editItem(itemId: number, content: string) {
    await api.put(`/survivor-letter/items/${itemId}`, { content })
    await fetchLetter()
  }

  async function reorderItems(itemOrders: { id: number; sort_order: number }[]) {
    await api.put('/survivor-letter/items/reorder', { item_orders: itemOrders })
    await fetchLetter()
  }

  async function deleteItem(itemId: number) {
    await api.delete(`/survivor-letter/items/${itemId}`)
    await fetchLetter()
  }

  async function unsuppressItem(itemId: number) {
    await api.post(`/survivor-letter/items/${itemId}/unsuppress`)
    await fetchLetter()
  }

  return {
    letter,
    loading,
    fetchLetter,
    updateBoilerplate,
    updateSection,
    reorderSections,
    addItem,
    editItem,
    reorderItems,
    deleteItem,
    unsuppressItem,
  }
})
