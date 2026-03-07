<template>
  <div class="print-template">
    <h2>Contacts</h2>
    <div v-for="contact in contacts" :key="contact.id" class="mb-3">
      <p>
        <strong>{{ contact.name }}</strong>
        <span v-if="contact.is_primary" style="font-weight: bold; color: #333;"> (Primary)</span>
      </p>
      <p v-if="contact.role" class="ml-4">Role: {{ contact.role }}</p>
      <p v-if="contact.relationship" class="ml-4">Relationship: {{ contact.relationship }}</p>
      <p v-if="contact.phone" class="ml-4">Phone: {{ contact.phone }}</p>
      <p v-if="contact.email" class="ml-4">Email: {{ contact.email }}</p>
      <p v-if="contact.address" class="ml-4 pre-wrap">Address: {{ contact.address }}</p>
      <p v-if="contact.notes" class="ml-4 pre-wrap">Notes: {{ contact.notes }}</p>
    </div>
    <p v-if="!contacts.length">No contacts recorded.</p>
  </div>
</template>

<script setup lang="ts">
import type { Contact } from '../types'

defineProps<{
  contacts: Contact[]
}>()
</script>
