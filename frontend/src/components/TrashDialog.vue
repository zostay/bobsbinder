<template>
  <v-dialog :model-value="modelValue" max-width="600" @update:model-value="$emit('update:modelValue', $event)">
    <v-card>
      <v-card-title>Trash</v-card-title>
      <v-card-text>
        <v-progress-linear v-if="loading" indeterminate color="primary" class="mb-3" />

        <div v-if="!loading && store.trashedDocuments.length === 0" class="text-center pa-4">
          <v-icon size="48" color="grey">mdi-delete-empty</v-icon>
          <p class="text-body-2 mt-2 text-grey">Trash is empty</p>
        </div>

        <v-list v-if="store.trashedDocuments.length > 0">
          <v-list-item v-for="doc in store.trashedDocuments" :key="doc.id">
            <v-list-item-title>{{ doc.title }}</v-list-item-title>
            <v-list-item-subtitle>
              Deleted {{ formatDate(doc.deleted_at!) }}
            </v-list-item-subtitle>

            <template v-slot:append>
              <v-btn size="small" variant="text" color="primary" @click="restore(doc.id)" :loading="restoring === doc.id">
                Restore
              </v-btn>
              <v-btn size="small" variant="text" color="error" @click="confirmPermanentDelete(doc)">
                Delete Forever
              </v-btn>
            </template>
          </v-list-item>
        </v-list>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn @click="close">Close</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog v-model="showConfirm" max-width="400">
    <v-card>
      <v-card-title>Delete Forever</v-card-title>
      <v-card-text>
        Are you sure you want to permanently delete "{{ deletingDoc?.title }}"? This cannot be undone.
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn @click="showConfirm = false">Cancel</v-btn>
        <v-btn color="error" @click="permanentDelete">Delete Forever</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useDocumentStore } from '../stores/documents'
import type { Document } from '../types'

const props = defineProps<{
  modelValue: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  restored: []
}>()

const store = useDocumentStore()
const loading = ref(false)
const restoring = ref<number | null>(null)
const showConfirm = ref(false)
const deletingDoc = ref<Document | null>(null)

watch(() => props.modelValue, async (open) => {
  if (open) {
    loading.value = true
    try {
      await store.fetchTrash()
    } finally {
      loading.value = false
    }
  }
})

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString()
}

async function restore(id: number) {
  restoring.value = id
  try {
    await store.restoreDocument(id)
    emit('restored')
  } finally {
    restoring.value = null
  }
}

function confirmPermanentDelete(doc: Document) {
  deletingDoc.value = doc
  showConfirm.value = true
}

async function permanentDelete() {
  if (!deletingDoc.value) return
  await store.permanentDeleteDocument(deletingDoc.value.id)
  showConfirm.value = false
  deletingDoc.value = null
}

function close() {
  emit('update:modelValue', false)
}
</script>
