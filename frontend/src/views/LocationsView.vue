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

    <v-dialog v-model="showForm" max-width="600">
      <v-card>
        <v-card-title>{{ editing ? 'Edit Location' : 'Add Location' }}</v-card-title>
        <v-card-text>
          <v-text-field v-model="form.name" label="Name" required />
          <v-select v-model="form.type" label="Type" :items="['physical', 'digital']" required />
          <v-textarea v-model="form.description" label="Description" rows="2" />
          <v-text-field v-model="form.address" label="Address" />
          <v-textarea v-model="form.access_instructions" label="Access Instructions" rows="2" />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn @click="closeForm">Cancel</v-btn>
          <v-btn color="primary" @click="save">Save</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { useLocationStore } from '../stores/locations'
import type { Location } from '../types'

const store = useLocationStore()
const showForm = ref(false)
const editing = ref<number | null>(null)
const form = reactive({
  name: '',
  type: 'physical' as 'physical' | 'digital',
  description: '',
  address: '',
  access_instructions: '',
})

function resetForm() {
  form.name = ''
  form.type = 'physical'
  form.description = ''
  form.address = ''
  form.access_instructions = ''
  editing.value = null
}

function startEdit(location: Location) {
  form.name = location.name
  form.type = location.type
  form.description = location.description || ''
  form.address = location.address || ''
  form.access_instructions = location.access_instructions || ''
  editing.value = location.id
  showForm.value = true
}

function closeForm() {
  showForm.value = false
  resetForm()
}

async function save() {
  if (editing.value) {
    await store.updateLocation(editing.value, { ...form })
  } else {
    await store.createLocation({ ...form })
  }
  closeForm()
}

onMounted(() => {
  store.fetchLocations()
})
</script>
