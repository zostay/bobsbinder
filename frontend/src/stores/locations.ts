import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'
import type { Location } from '../types'

export const useLocationStore = defineStore('locations', () => {
  const locations = ref<Location[]>([])
  const loading = ref(false)

  async function fetchLocations() {
    loading.value = true
    try {
      const { data } = await api.get<Location[]>('/locations')
      locations.value = data
    } finally {
      loading.value = false
    }
  }

  async function createLocation(location: Partial<Location>) {
    const { data } = await api.post<Location>('/locations', location)
    locations.value.push(data)
    return data
  }

  async function updateLocation(id: number, location: Partial<Location>) {
    await api.put(`/locations/${id}`, location)
    await fetchLocations()
  }

  async function deleteLocation(id: number) {
    await api.delete(`/locations/${id}`)
    locations.value = locations.value.filter((l) => l.id !== id)
  }

  return { locations, loading, fetchLocations, createLocation, updateLocation, deleteLocation }
})
