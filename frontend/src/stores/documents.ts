import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'
import type { Document, DocumentFile } from '../types'

export const useDocumentStore = defineStore('documents', () => {
  const documents = ref<Document[]>([])
  const trashedDocuments = ref<Document[]>([])
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

  async function fetchDocument(id: number) {
    const { data } = await api.get<Document>(`/documents/${id}`)
    return data
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

  async function fetchTrash() {
    const { data } = await api.get<Document[]>('/documents/trash')
    trashedDocuments.value = data
  }

  async function restoreDocument(id: number) {
    await api.post(`/documents/${id}/restore`)
    trashedDocuments.value = trashedDocuments.value.filter((d) => d.id !== id)
  }

  async function permanentDeleteDocument(id: number) {
    await api.delete(`/documents/${id}/permanent`)
    trashedDocuments.value = trashedDocuments.value.filter((d) => d.id !== id)
  }

  async function fetchFiles(documentId: number) {
    const { data } = await api.get<DocumentFile[]>(`/documents/${documentId}/files`)
    return data
  }

  async function uploadFile(documentId: number, file: File) {
    const formData = new FormData()
    formData.append('file', file)
    const { data } = await api.post<DocumentFile>(`/documents/${documentId}/files`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
    return data
  }

  async function downloadFile(documentId: number, fileId: number, filename: string) {
    const { data } = await api.get(`/documents/${documentId}/files/${fileId}`, {
      responseType: 'blob',
    })
    const url = window.URL.createObjectURL(data as Blob)
    const a = document.createElement('a')
    a.href = url
    a.download = filename
    a.click()
    window.URL.revokeObjectURL(url)
  }

  async function deleteFile(documentId: number, fileId: number) {
    await api.delete(`/documents/${documentId}/files/${fileId}`)
  }

  return {
    documents, trashedDocuments, loading,
    fetchDocuments, fetchDocument, createDocument, updateDocument, deleteDocument,
    fetchTrash, restoreDocument, permanentDeleteDocument,
    fetchFiles, uploadFile, downloadFile, deleteFile,
  }
})
