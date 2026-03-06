<template>
  <v-dialog :model-value="modelValue" max-width="600" @update:model-value="$emit('update:modelValue', $event)">
    <v-card>
      <v-card-title>{{ editData ? 'Edit Contact' : 'Add Contact' }}</v-card-title>
      <v-card-text>
        <v-alert type="info" variant="tonal" density="compact" class="mb-3">
          Fields marked with <v-icon size="small">mdi-lock</v-icon> are confidential
          and will not be sent to the printing service. Only place passwords, PINs,
          access codes, or account numbers in locked fields.
        </v-alert>
        <v-text-field v-model="form.name" label="Name" required />
        <v-checkbox v-model="form.is_primary" label="Primary contact" hide-details class="mb-2" />
        <v-text-field v-model="form.role" label="Role (e.g. Attorney, Pastor)" />
        <v-text-field v-model="form.relationship" label="Relationship" />
        <v-text-field v-model="form.phone" label="Phone" />
        <v-text-field v-model="form.email" label="Email" />
        <v-textarea v-model="form.address" label="Address" rows="2" />
        <v-textarea v-model="form.notes" label="Notes" rows="2" />
        <v-textarea v-model="form.secure_notes" label="Confidential Notes"
          rows="2" prepend-inner-icon="mdi-lock"
          hint="Will NOT appear on the printed cover letter."
          persistent-hint />
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn @click="close">Cancel</v-btn>
        <v-btn color="primary" @click="save">Save</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { reactive, watch } from 'vue'
import { useContactStore } from '../stores/contacts'
import type { Contact } from '../types'

const props = defineProps<{
  modelValue: boolean
  editData?: Contact | null
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  saved: []
}>()

const store = useContactStore()

const form = reactive({
  name: '',
  is_primary: false,
  role: '',
  relationship: '',
  phone: '',
  email: '',
  address: '',
  notes: '',
  secure_notes: '',
})

function resetForm() {
  form.name = ''
  form.is_primary = false
  form.role = ''
  form.relationship = ''
  form.phone = ''
  form.email = ''
  form.address = ''
  form.notes = ''
  form.secure_notes = ''
}

watch(() => props.modelValue, (open) => {
  if (open) {
    if (props.editData) {
      form.name = props.editData.name
      form.is_primary = props.editData.is_primary || false
      form.role = props.editData.role || ''
      form.relationship = props.editData.relationship || ''
      form.phone = props.editData.phone || ''
      form.email = props.editData.email || ''
      form.address = props.editData.address || ''
      form.notes = props.editData.notes || ''
      form.secure_notes = props.editData.secure_notes || ''
    } else {
      resetForm()
    }
  }
})

function close() {
  emit('update:modelValue', false)
  resetForm()
}

async function save() {
  if (props.editData) {
    await store.updateContact(props.editData.id, { ...form })
  } else {
    await store.createContact({ ...form })
  }
  close()
  emit('saved')
}
</script>
