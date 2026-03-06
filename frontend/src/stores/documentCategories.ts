import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'
import type { DocumentCategory } from '../types'

export const useDocumentCategoryStore = defineStore('documentCategories', () => {
  const categories = ref<DocumentCategory[]>([])
  const loading = ref(false)

  async function fetchCategories() {
    if (categories.value.length > 0) return
    loading.value = true
    try {
      const { data } = await api.get<DocumentCategory[]>('/document-categories')
      categories.value = data
    } finally {
      loading.value = false
    }
  }

  return { categories, loading, fetchCategories }
})
