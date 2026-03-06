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
      <v-card-title :class="{ 'font-weight-bold': contact.is_primary }">
        {{ contact.name }}
        <v-chip v-if="contact.is_primary" size="small" color="primary" class="ml-2">Primary</v-chip>
      </v-card-title>
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

    <ContactFormDialog
      v-model="showForm"
      :edit-data="editingContact"
      @saved="onSaved"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useContactStore } from '../stores/contacts'
import ContactFormDialog from '../components/ContactFormDialog.vue'
import type { Contact } from '../types'

const store = useContactStore()
const showForm = ref(false)
const editingContact = ref<Contact | null>(null)

function startEdit(contact: Contact) {
  editingContact.value = contact
  showForm.value = true
}

function onSaved() {
  editingContact.value = null
  store.fetchContacts()
}

onMounted(() => {
  store.fetchContacts()
})
</script>
