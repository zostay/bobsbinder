<template>
  <div>
    <div class="d-flex align-center mb-4">
      <h1 class="text-h4">Contacts</h1>
      <v-spacer />
      <v-btn color="primary" prepend-icon="mdi-plus" @click="showForm = true">
        Add Contact
      </v-btn>
    </div>

    <v-progress-linear v-if="store.loading" indeterminate color="primary" />

    <v-card v-if="store.contacts.length === 0 && !store.loading">
      <v-card-text class="text-center pa-8">
        <v-icon size="64" color="grey">mdi-account-group-outline</v-icon>
        <p class="text-h6 mt-4">No contacts yet</p>
        <p class="text-body-2 text-grey">Add people who should be contacted — your attorney, pastor, family members, etc.</p>
      </v-card-text>
    </v-card>

    <v-card v-for="contact in store.contacts" :key="contact.id" class="mb-2">
      <v-card-title>{{ contact.name }}</v-card-title>
      <v-card-subtitle v-if="contact.role">{{ contact.role }}</v-card-subtitle>
      <v-card-text>
        <div v-if="contact.relationship">Relationship: {{ contact.relationship }}</div>
        <div v-if="contact.phone">Phone: {{ contact.phone }}</div>
        <div v-if="contact.email">Email: {{ contact.email }}</div>
        <div v-if="contact.address">Address: {{ contact.address }}</div>
        <div v-if="contact.notes" class="mt-2 text-grey">{{ contact.notes }}</div>
      </v-card-text>
      <v-card-actions>
        <v-btn size="small" variant="text" color="primary" @click="startEdit(contact)">Edit</v-btn>
        <v-btn size="small" variant="text" color="error" @click="store.deleteContact(contact.id)">Delete</v-btn>
      </v-card-actions>
    </v-card>

    <v-dialog v-model="showForm" max-width="600">
      <v-card>
        <v-card-title>{{ editing ? 'Edit Contact' : 'Add Contact' }}</v-card-title>
        <v-card-text>
          <v-text-field v-model="form.name" label="Name" required />
          <v-text-field v-model="form.role" label="Role (e.g. Attorney, Pastor)" />
          <v-text-field v-model="form.relationship" label="Relationship" />
          <v-text-field v-model="form.phone" label="Phone" />
          <v-text-field v-model="form.email" label="Email" />
          <v-textarea v-model="form.address" label="Address" rows="2" />
          <v-textarea v-model="form.notes" label="Notes" rows="2" />
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
import { useContactStore } from '../stores/contacts'
import type { Contact } from '../types'

const store = useContactStore()
const showForm = ref(false)
const editing = ref<number | null>(null)
const form = reactive({
  name: '',
  role: '',
  relationship: '',
  phone: '',
  email: '',
  address: '',
  notes: '',
})

function resetForm() {
  form.name = ''
  form.role = ''
  form.relationship = ''
  form.phone = ''
  form.email = ''
  form.address = ''
  form.notes = ''
  editing.value = null
}

function startEdit(contact: Contact) {
  form.name = contact.name
  form.role = contact.role || ''
  form.relationship = contact.relationship || ''
  form.phone = contact.phone || ''
  form.email = contact.email || ''
  form.address = contact.address || ''
  form.notes = contact.notes || ''
  editing.value = contact.id
  showForm.value = true
}

function closeForm() {
  showForm.value = false
  resetForm()
}

async function save() {
  if (editing.value) {
    await store.updateContact(editing.value, { ...form })
  } else {
    await store.createContact({ ...form })
  }
  closeForm()
}

onMounted(() => {
  store.fetchContacts()
})
</script>
