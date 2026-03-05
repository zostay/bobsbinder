<template>
  <div>
    <div class="d-flex align-center mb-4">
      <h1 class="text-h4">Digital Access</h1>
      <v-spacer />
      <v-btn color="primary" prepend-icon="mdi-plus" @click="showForm = true">
        Add Entry
      </v-btn>
    </div>

    <v-progress-linear v-if="store.loading" indeterminate color="primary" />

    <v-card v-if="store.items.length === 0 && !store.loading">
      <v-card-text class="text-center pa-8">
        <v-icon size="64" color="grey">mdi-laptop</v-icon>
        <p class="text-h6 mt-4">No digital access entries yet</p>
        <p class="text-body-2 text-grey">Add computer passwords, phone access codes, and password manager info.</p>
      </v-card-text>
    </v-card>

    <v-card v-for="item in store.items" :key="item.id" class="mb-2">
      <v-card-title>{{ item.name }}</v-card-title>
      <v-card-subtitle>
        <v-chip size="small" color="info">{{ typeLabels[item.type] || item.type }}</v-chip>
      </v-card-subtitle>
      <v-card-text>
        <div v-if="item.username">Username: {{ item.username }}</div>
        <div v-if="item.instructions">Instructions: {{ item.instructions }}</div>
      </v-card-text>
      <v-card-actions>
        <v-btn size="small" variant="text" color="primary" @click="startEdit(item)">Edit</v-btn>
        <v-btn size="small" variant="text" color="error" @click="store.deleteItem(item.id)">Delete</v-btn>
      </v-card-actions>
    </v-card>

    <v-dialog v-model="showForm" max-width="600">
      <v-card>
        <v-card-title>{{ editing ? 'Edit Entry' : 'Add Entry' }}</v-card-title>
        <v-card-text>
          <v-text-field v-model="form.name" label="Name" required />
          <v-select
            v-model="form.type"
            label="Type"
            :items="[
              { title: 'Computer', value: 'computer' },
              { title: 'Phone', value: 'phone' },
              { title: 'Password Manager', value: 'password_manager' },
            ]"
            required
          />
          <v-text-field v-model="form.username" label="Username" />
          <v-textarea v-model="form.instructions" label="Access Instructions" rows="3" />
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
import { useDigitalAccessStore } from '../stores/digitalAccess'
import type { DigitalAccess } from '../types'

const store = useDigitalAccessStore()
const showForm = ref(false)
const editing = ref<number | null>(null)

const typeLabels: Record<string, string> = {
  computer: 'Computer',
  phone: 'Phone',
  password_manager: 'Password Manager',
}

const form = reactive({
  name: '',
  type: 'computer' as 'computer' | 'phone' | 'password_manager',
  username: '',
  instructions: '',
})

function resetForm() {
  form.name = ''
  form.type = 'computer'
  form.username = ''
  form.instructions = ''
  editing.value = null
}

function startEdit(item: DigitalAccess) {
  form.name = item.name
  form.type = item.type
  form.username = item.username || ''
  form.instructions = item.instructions || ''
  editing.value = item.id
  showForm.value = true
}

function closeForm() {
  showForm.value = false
  resetForm()
}

async function save() {
  if (editing.value) {
    await store.updateItem(editing.value, { ...form })
  } else {
    await store.createItem({ ...form })
  }
  closeForm()
}

onMounted(() => {
  store.fetchItems()
})
</script>
