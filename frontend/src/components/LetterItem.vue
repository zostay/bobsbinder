<template>
  <div class="d-flex align-center py-1" :class="{ 'text-grey': item.suppressed }">
    <v-icon size="small" class="mr-2 cursor-grab">mdi-drag-vertical</v-icon>

    <span v-if="item.item_type === 'bulleted'" class="mr-2">&#x2022;</span>
    <span v-else-if="item.item_type === 'numbered'" class="mr-2">{{ index + 1 }}.</span>

    <template v-if="editingContent">
      <v-text-field
        v-model="editText"
        density="compact"
        variant="outlined"
        hide-details
        class="flex-grow-1"
        @keyup.enter="saveEdit"
        @keyup.escape="cancelEdit"
      />
      <v-btn icon size="small" color="success" @click="saveEdit" class="ml-1">
        <v-icon size="small">mdi-check</v-icon>
      </v-btn>
      <v-btn icon size="small" @click="cancelEdit" class="ml-1">
        <v-icon size="small">mdi-close</v-icon>
      </v-btn>
    </template>
    <template v-else>
      <span class="flex-grow-1" :class="{ 'text-decoration-line-through': item.suppressed }">
        {{ item.content }}
      </span>

      <v-chip v-if="item.provenance !== 'auto'" size="x-small" class="ml-2" variant="outlined">
        {{ item.provenance === 'auto_edited' ? 'edited' : 'manual' }}
      </v-chip>

      <v-btn
        v-if="!item.suppressed"
        icon size="small" variant="text" @click="startEdit" class="ml-1"
      >
        <v-icon size="small">mdi-pencil</v-icon>
      </v-btn>

      <v-btn
        v-if="item.suppressed"
        icon size="small" variant="text" color="success" @click="$emit('unsuppress', item.id)" class="ml-1"
      >
        <v-icon size="small">mdi-eye</v-icon>
      </v-btn>
      <v-btn
        v-else
        icon size="small" variant="text" color="error" @click="$emit('delete', item.id)" class="ml-1"
      >
        <v-icon size="small">mdi-delete</v-icon>
      </v-btn>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { SurvivorLetterItem } from '../types'

const props = defineProps<{
  item: SurvivorLetterItem
  index: number
}>()

const emit = defineEmits<{
  edit: [id: number, content: string]
  delete: [id: number]
  unsuppress: [id: number]
}>()

const editingContent = ref(false)
const editText = ref('')

function startEdit() {
  editText.value = props.item.content
  editingContent.value = true
}

function saveEdit() {
  emit('edit', props.item.id, editText.value)
  editingContent.value = false
}

function cancelEdit() {
  editingContent.value = false
}
</script>
