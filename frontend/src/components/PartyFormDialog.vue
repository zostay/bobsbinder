<template>
  <v-dialog :model-value="modelValue" max-width="500" @update:model-value="$emit('update:modelValue', $event)">
    <v-card>
      <v-card-title>{{ editData ? 'Edit Party' : 'Add Party' }}</v-card-title>
      <v-card-text>
        <v-text-field v-model="form.name" label="Name" required />
        <v-select
          v-model="form.relationship"
          label="Relationship"
          :items="relationshipOptions"
        />
        <v-textarea v-model="form.notes" label="Notes" rows="2" />
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
import { usePartyStore } from '../stores/parties'
import type { Party } from '../types'

const props = defineProps<{
  modelValue: boolean
  editData?: Party | null
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  saved: []
}>()

const store = usePartyStore()

const relationshipOptions = [
  { title: 'Self', value: 'self' },
  { title: 'Spouse', value: 'spouse' },
  { title: 'Dependent', value: 'dependent' },
  { title: 'Other', value: 'other' },
]

const form = reactive({
  name: '',
  relationship: 'other' as string,
  notes: '',
})

function resetForm() {
  form.name = ''
  form.relationship = 'other'
  form.notes = ''
}

watch(() => props.modelValue, (open) => {
  if (open) {
    if (props.editData) {
      form.name = props.editData.name
      form.relationship = props.editData.relationship
      form.notes = props.editData.notes || ''
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
    await store.updateParty(props.editData.id, { ...form })
  } else {
    await store.createParty({ ...form })
  }
  close()
  emit('saved')
}
</script>
