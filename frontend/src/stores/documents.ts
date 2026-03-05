import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'
import type { Document } from '../types'

export const useDocumentStore = defineStore('documents', () => {
  const documents = ref<Document[]>([])
  const loading = ref(false)

  async function fetchDocuments() {
    loading.value = true
    try {
      const { data } = await api.get<Document[]>('/documents')
      documents.value = data
    } finally {
      loading.value = false
    }
  }

  async function createDocument(doc: Partial<Document>) {
    const { data } = await api.post<Document>('/documents', doc)
    documents.value.unshift(data)
    return data
  }

  async function updateDocument(id: number, doc: Partial<Document>) {
    await api.put(`/documents/${id}`, doc)
    await fetchDocuments()
  }

  async function deleteDocument(id: number) {
    await api.delete(`/documents/${id}`)
    documents.value = documents.value.filter((d) => d.id !== id)
  }

  return { documents, loading, fetchDocuments, createDocument, updateDocument, deleteDocument }
})
