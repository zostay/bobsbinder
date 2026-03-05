import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'
import type { Contact } from '../types'

export const useContactStore = defineStore('contacts', () => {
  const contacts = ref<Contact[]>([])
  const loading = ref(false)

  async function fetchContacts() {
    loading.value = true
    try {
      const { data } = await api.get<Contact[]>('/contacts')
      contacts.value = data
    } finally {
      loading.value = false
    }
  }

  async function createContact(contact: Partial<Contact>) {
    const { data } = await api.post<Contact>('/contacts', contact)
    contacts.value.push(data)
    return data
  }

  async function updateContact(id: number, contact: Partial<Contact>) {
    await api.put(`/contacts/${id}`, contact)
    await fetchContacts()
  }

  async function deleteContact(id: number) {
    await api.delete(`/contacts/${id}`)
    contacts.value = contacts.value.filter((c) => c.id !== id)
  }

  return { contacts, loading, fetchContacts, createContact, updateContact, deleteContact }
})
