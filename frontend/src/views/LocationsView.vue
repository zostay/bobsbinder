<template>
  <div>
    <div class="d-flex align-center mb-4">
      <h1 class="text-h4">Locations</h1>
      <v-spacer />
      <v-btn color="primary" prepend-icon="mdi-plus" @click="showForm = true">
        Add Location
      </v-btn>
    </div>

    <v-progress-linear v-if="store.loading" indeterminate color="primary" />

    <v-card v-if="store.locations.length === 0 && !store.loading">
      <v-card-text class="text-center pa-8">
        <v-icon size="64" color="grey">mdi-map-marker-outline</v-icon>
        <p class="text-h6 mt-4">No locations yet</p>
        <p class="text-body-2 text-grey">Add physical and digital locations where important items are kept.</p>
      </v-card-text>
    </v-card>

    <v-card v-for="location in store.locations" :key="location.id" class="mb-2">
      <v-card-title>{{ location.name }}</v-card-title>
      <v-card-subtitle>
        <v-chip size="small" :color="location.type === 'digital' ? 'info' : 'success'">
          {{ location.type }}
        </v-chip>
      </v-card-subtitle>
      <v-card-text>
        <div v-if="location.description">{{ location.description }}</div>
        <div v-if="location.address">Address: {{ location.address }}</div>
        <div v-if="location.access_instructions" class="mt-2 text-grey">Access: {{ location.access_instructions }}</div>
      </v-card-text>
      <v-card-actions>
        <v-btn size="small" variant="text" color="primary" @click="startEdit(location)">Edit</v-btn>
        <v-btn size="small" variant="text" color="error" @click="store.deleteLocation(location.id)">Delete</v-btn>
      </v-card-actions>
    </v-card>

    <LocationFormDialog
      v-model="showForm"
      :edit-data="editingLocation"
      @saved="onSaved"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useLocationStore } from '../stores/locations'
import LocationFormDialog from '../components/LocationFormDialog.vue'
import type { Location } from '../types'

const store = useLocationStore()
const showForm = ref(false)
const editingLocation = ref<Location | null>(null)

function startEdit(location: Location) {
  editingLocation.value = location
  showForm.value = true
}

function onSaved() {
  editingLocation.value = null
  store.fetchLocations()
}

onMounted(() => {
  store.fetchLocations()
})
</script>
