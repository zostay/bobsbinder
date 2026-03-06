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

    <DigitalInfoFormDialog
      v-model="showForm"
      :edit-data="editingItem"
      initial-type="computer"
      @saved="onSaved"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useDigitalAccessStore } from '../stores/digitalAccess'
import DigitalInfoFormDialog from '../components/DigitalInfoFormDialog.vue'
import type { DigitalAccess } from '../types'

const store = useDigitalAccessStore()
const showForm = ref(false)
const editingItem = ref<DigitalAccess | null>(null)

const typeLabels: Record<string, string> = {
  computer: 'Computer',
  phone: 'Phone',
  password_manager: 'Password Manager',
}

function startEdit(item: DigitalAccess) {
  editingItem.value = item
  showForm.value = true
}

function onSaved() {
  editingItem.value = null
  store.fetchItems()
}

onMounted(() => {
  store.fetchItems()
})
</script>
