<template>
  <v-dialog :model-value="modelValue" max-width="600" @update:model-value="$emit('update:modelValue', $event)">
    <v-card>
      <v-card-title>{{ editData ? 'Edit Reference Document' : 'Add Reference Document' }}</v-card-title>
      <v-card-text>
        <v-alert type="info" variant="tonal" density="compact" class="mb-3">
          Fields marked with <v-icon size="small">mdi-lock</v-icon> are confidential
          and will not be sent to the printing service. Only place passwords, PINs,
          access codes, or account numbers in locked fields.
        </v-alert>
        <v-text-field v-model="form.title" label="Title" required />
        <v-select
          v-model="form.category_id"
          label="Category"
          :items="categoryItems"
          item-title="name"
          item-value="id"
        />
        <v-textarea v-model="form.content" label="Description" rows="3" />
        <v-select
          v-model="form.location_id"
          label="Location"
          :items="locationItems"
          item-title="name"
          item-value="id"
          clearable
        />
        <v-file-input
          v-if="!editData"
          v-model="pendingFile"
          label="Attach File (optional)"
          prepend-icon="mdi-paperclip"
        />

        <template v-if="editData && files.length > 0">
          <h4 class="text-subtitle-2 mt-3 mb-1">Attached Files</h4>
          <v-list density="compact">
            <v-list-item v-for="f in files" :key="f.id">
              <v-list-item-title>{{ f.filename }}</v-list-item-title>
              <v-list-item-subtitle>{{ formatFileSize(f.file_size) }}</v-list-item-subtitle>
              <template v-slot:append>
                <v-btn icon size="small" variant="text" @click="handleDownload(f)">
                  <v-icon>mdi-download</v-icon>
                </v-btn>
                <v-btn icon size="small" variant="text" color="error" @click="handleDeleteFile(f.id)">
                  <v-icon>mdi-delete</v-icon>
                </v-btn>
              </template>
            </v-list-item>
          </v-list>
        </template>

        <v-file-input
          v-if="editData"
          v-model="pendingFile"
          label="Upload Another File"
          prepend-icon="mdi-paperclip"
          class="mt-2"
        />

        <v-select
          v-model="form.status"
          label="Status"
          :items="['draft', 'complete']"
        />
        <v-textarea v-model="form.secure_notes" label="Confidential Notes"
          rows="2" prepend-inner-icon="mdi-lock"
          hint="Will NOT appear on the printed cover letter."
          persistent-hint />
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn @click="close">Cancel</v-btn>
        <v-btn color="primary" :loading="saving" @click="save">Save</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { useDocumentStore } from '../stores/documents'
import { useDocumentCategoryStore } from '../stores/documentCategories'
import { useLocationStore } from '../stores/locations'
import type { Document, DocumentFile } from '../types'

const props = defineProps<{
  modelValue: boolean
  editData?: Document | null
  initialCategoryId?: number
  partyId?: number
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  saved: []
}>()

const documentStore = useDocumentStore()
const categoryStore = useDocumentCategoryStore()
const locationStore = useLocationStore()

const saving = ref(false)
const pendingFile = ref<File[] | null>(null)
const files = ref<DocumentFile[]>([])

const categoryItems = computed(() => categoryStore.categories)
const locationItems = computed(() => locationStore.locations)

const form = reactive({
  title: '',
  category_id: null as number | null,
  content: '',
  location_id: null as number | null,
  status: 'draft' as 'draft' | 'complete',
  secure_notes: '',
})

function resetForm() {
  form.title = ''
  form.category_id = props.initialCategoryId ?? null
  form.content = ''
  form.location_id = null
  form.status = 'draft'
  form.secure_notes = ''
  pendingFile.value = null
  files.value = []
}

watch(() => props.modelValue, async (open) => {
  if (open) {
    if (props.editData) {
      form.title = props.editData.title
      form.category_id = props.editData.category_id
      form.content = props.editData.content || ''
      form.location_id = props.editData.location_id ?? null
      form.status = props.editData.status
      form.secure_notes = props.editData.secure_notes || ''
      pendingFile.value = null
      files.value = await documentStore.fetchFiles(props.editData.id)
    } else {
      resetForm()
    }
  }
})

function close() {
  emit('update:modelValue', false)
  resetForm()
}

function formatFileSize(bytes: number) {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

async function handleDownload(f: DocumentFile) {
  await documentStore.downloadFile(f.document_id, f.id, f.filename)
}

async function handleDeleteFile(fileId: number) {
  if (!props.editData) return
  await documentStore.deleteFile(props.editData.id, fileId)
  files.value = files.value.filter((f) => f.id !== fileId)
}

async function save() {
  saving.value = true
  try {
    const payload: Record<string, any> = {
      ...form,
      doc_type: 'reference',
    }
    if (props.partyId) payload.party_id = props.partyId

    let docId: number
    if (props.editData) {
      await documentStore.updateDocument(props.editData.id, payload)
      docId = props.editData.id
    } else {
      const created = await documentStore.createDocument(payload)
      docId = created.id
    }

    if (pendingFile.value && pendingFile.value.length > 0) {
      await documentStore.uploadFile(docId, pendingFile.value[0])
    }

    close()
    emit('saved')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  categoryStore.fetchCategories()
  locationStore.fetchLocations()
})
</script>
