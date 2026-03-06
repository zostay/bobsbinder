<template>
  <v-dialog :model-value="modelValue" max-width="600" @update:model-value="$emit('update:modelValue', $event)">
    <v-card>
      <v-card-title>{{ editData ? 'Edit Location' : 'Add Location' }}</v-card-title>
      <v-card-text>
        <v-text-field v-model="form.name" label="Name" required />
        <v-select v-model="form.type" label="Type" :items="['physical', 'digital']" required />
        <v-textarea v-model="form.description" label="Description" rows="2" />
        <v-text-field v-model="form.address" label="Address" />
        <v-textarea v-model="form.access_instructions" label="Access Instructions" rows="2" />
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
import { useLocationStore } from '../stores/locations'
import type { Location } from '../types'

const props = defineProps<{
  modelValue: boolean
  editData?: Location | null
  initialType?: 'physical' | 'digital'
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  saved: []
}>()

const store = useLocationStore()

const form = reactive({
  name: '',
  type: 'physical' as 'physical' | 'digital',
  description: '',
  address: '',
  access_instructions: '',
})

function resetForm() {
  form.name = ''
  form.type = props.initialType || 'physical'
  form.description = ''
  form.address = ''
  form.access_instructions = ''
}

watch(() => props.modelValue, (open) => {
  if (open) {
    if (props.editData) {
      form.name = props.editData.name
      form.type = props.editData.type
      form.description = props.editData.description || ''
      form.address = props.editData.address || ''
      form.access_instructions = props.editData.access_instructions || ''
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
    await store.updateLocation(props.editData.id, { ...form })
  } else {
    await store.createLocation({ ...form })
  }
  close()
  emit('saved')
}
</script>
