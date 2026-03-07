<template>
  <div>
    <div class="d-flex align-center mb-4">
      <v-btn icon variant="text" @click="router.push({ name: 'documents' })">
        <v-icon>mdi-arrow-left</v-icon>
      </v-btn>
      <h1 class="text-h4 ml-2">{{ doc?.title }}</h1>
      <v-spacer />
      <v-btn color="primary" variant="text" prepend-icon="mdi-pencil" @click="editDoc">Edit</v-btn>
      <v-btn color="error" variant="text" prepend-icon="mdi-delete" @click="confirmDelete">Delete</v-btn>
    </div>

    <v-progress-linear v-if="loading" indeterminate color="primary" />

    <template v-if="doc && !loading">
      <v-row>
        <v-col cols="12" md="8">
          <v-card v-if="doc.doc_type === 'typed'" class="mb-4">
            <v-card-text class="markdown-body" v-html="renderedContent" />
          </v-card>

          <v-card v-else class="mb-4">
            <v-card-text>
              <p v-if="doc.content" class="text-body-1">{{ doc.content }}</p>
              <p v-else class="text-body-2 text-grey">No description provided.</p>
            </v-card-text>
          </v-card>

          <v-card v-if="doc.doc_type === 'reference' && files.length > 0" class="mb-4">
            <v-card-title class="text-subtitle-1">Attached Files</v-card-title>
            <v-list density="compact">
              <v-list-item v-for="f in files" :key="f.id">
                <v-list-item-title>{{ f.filename }}</v-list-item-title>
                <v-list-item-subtitle>{{ formatFileSize(f.file_size) }}</v-list-item-subtitle>
                <template v-slot:append>
                  <v-btn icon size="small" variant="text" @click="handleDownload(f)">
                    <v-icon>mdi-download</v-icon>
                  </v-btn>
                </template>
              </v-list-item>
            </v-list>
          </v-card>
        </v-col>

        <v-col cols="12" md="4">
          <v-card>
            <v-card-title class="text-subtitle-1">Details</v-card-title>
            <v-list density="compact">
              <v-list-item>
                <v-list-item-title>Type</v-list-item-title>
                <v-list-item-subtitle>
                  <v-icon size="small" class="mr-1">{{ doc.doc_type === 'typed' ? 'mdi-file-document-edit-outline' : 'mdi-link-variant' }}</v-icon>
                  {{ doc.doc_type === 'typed' ? 'Written Document' : 'Reference Document' }}
                </v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title>Category</v-list-item-title>
                <v-list-item-subtitle>{{ categoryName }}</v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title>Status</v-list-item-title>
                <v-list-item-subtitle>
                  <v-chip size="small" :color="doc.status === 'complete' ? 'success' : 'warning'">
                    {{ doc.status }}
                  </v-chip>
                </v-list-item-subtitle>
              </v-list-item>
              <v-list-item v-if="locationName">
                <v-list-item-title>Location</v-list-item-title>
                <v-list-item-subtitle>{{ locationName }}</v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title>Created</v-list-item-title>
                <v-list-item-subtitle>{{ formatDate(doc.created_at) }}</v-list-item-subtitle>
              </v-list-item>
              <v-list-item>
                <v-list-item-title>Updated</v-list-item-title>
                <v-list-item-subtitle>{{ formatDate(doc.updated_at) }}</v-list-item-subtitle>
              </v-list-item>
            </v-list>
          </v-card>
        </v-col>
      </v-row>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useDocumentStore } from '../stores/documents'
import { useDocumentCategoryStore } from '../stores/documentCategories'
import { useLocationStore } from '../stores/locations'
import { marked } from 'marked'
import type { Document, DocumentFile } from '../types'

const route = useRoute()
const router = useRouter()
const documentStore = useDocumentStore()
const categoryStore = useDocumentCategoryStore()
const locationStore = useLocationStore()

const doc = ref<Document | null>(null)
const files = ref<DocumentFile[]>([])
const loading = ref(true)

const renderedContent = computed(() => {
  if (!doc.value?.content) return ''
  return marked(doc.value.content) as string
})

const categoryName = computed(() => {
  if (!doc.value) return ''
  const cat = categoryStore.categories.find((c) => c.id === doc.value!.category_id)
  return cat?.name || 'Unknown'
})

const locationName = computed(() => {
  if (!doc.value?.location_id) return ''
  const loc = locationStore.locations.find((l) => l.id === doc.value!.location_id)
  return loc?.name || ''
})

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString()
}

function formatFileSize(bytes: number) {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

function editDoc() {
  if (!doc.value) return
  if (doc.value.doc_type === 'typed') {
    router.push({ name: 'document-edit', params: { id: doc.value.id } })
  } else {
    router.push({ name: 'documents', query: { edit: doc.value.id.toString() } })
  }
}

async function confirmDelete() {
  if (!doc.value) return
  await documentStore.deleteDocument(doc.value.id)
  router.push({ name: 'documents' })
}

async function handleDownload(f: DocumentFile) {
  await documentStore.downloadFile(f.document_id, f.id, f.filename)
}

onMounted(async () => {
  await Promise.all([
    categoryStore.fetchCategories(),
    locationStore.fetchLocations(),
  ])

  const id = Number(route.params.id)
  try {
    doc.value = await documentStore.fetchDocument(id)
    if (doc.value.doc_type === 'reference') {
      files.value = await documentStore.fetchFiles(id)
    }
  } finally {
    loading.value = false
  }
})
</script>
