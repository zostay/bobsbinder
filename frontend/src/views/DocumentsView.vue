<template>
  <div>
    <div class="d-flex align-center mb-4">
      <h1 class="text-h4">Documents</h1>
      <v-spacer />
      <v-menu>
        <template v-slot:activator="{ props }">
          <v-btn color="primary" prepend-icon="mdi-plus" append-icon="mdi-menu-down" v-bind="props">
            Add Document
          </v-btn>
        </template>
        <v-list>
          <v-list-item @click="showRefForm = true">
            <template v-slot:prepend>
              <v-icon>mdi-link-variant</v-icon>
            </template>
            <v-list-item-title>Reference Document</v-list-item-title>
          </v-list-item>
          <v-list-item @click="router.push({ name: 'document-new' })">
            <template v-slot:prepend>
              <v-icon>mdi-file-document-edit-outline</v-icon>
            </template>
            <v-list-item-title>Written Document</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
      <v-menu>
        <template v-slot:activator="{ props }">
          <v-btn icon variant="text" v-bind="props" class="ml-2">
            <v-icon>mdi-dots-vertical</v-icon>
          </v-btn>
        </template>
        <v-list>
          <v-list-item @click="showTrash = true">
            <template v-slot:prepend>
              <v-icon>mdi-delete-outline</v-icon>
            </template>
            <v-list-item-title>View Trashed Documents</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
    </div>

    <v-progress-linear v-if="store.loading" indeterminate color="primary" />

    <v-card v-if="store.documents.length === 0 && !store.loading">
      <v-card-text class="text-center pa-8">
        <v-icon size="64" color="grey">mdi-file-document-outline</v-icon>
        <p class="text-h6 mt-4">No documents yet</p>
        <p class="text-body-2 text-grey">Start adding your important documents to get organized.</p>
      </v-card-text>
    </v-card>

    <v-card v-for="doc in store.documents" :key="doc.id" class="mb-2">
      <v-card-title class="d-flex align-center">
        <v-icon size="small" class="mr-2">
          {{ doc.doc_type === 'typed' ? 'mdi-file-document-edit-outline' : 'mdi-link-variant' }}
        </v-icon>
        <router-link :to="{ name: 'document-detail', params: { id: doc.id } }" class="text-decoration-none">
          {{ doc.title }}
        </router-link>
      </v-card-title>
      <v-card-subtitle>
        <v-chip size="small" :color="doc.status === 'complete' ? 'success' : 'warning'" class="mr-2">
          {{ doc.status }}
        </v-chip>
        <span v-if="getCategoryName(doc.category_id)" class="text-caption">{{ getCategoryName(doc.category_id) }}</span>
      </v-card-subtitle>
      <v-card-text v-if="doc.content && doc.doc_type === 'reference'">{{ doc.content }}</v-card-text>
      <v-card-actions>
        <v-btn size="small" variant="text" prepend-icon="mdi-eye" :to="{ name: 'document-detail', params: { id: doc.id } }">View</v-btn>
        <v-btn size="small" variant="text" color="primary" prepend-icon="mdi-pencil" @click="startEdit(doc)">Edit</v-btn>
        <v-spacer />
        <v-menu>
          <template v-slot:activator="{ props }">
            <v-btn icon size="small" variant="text" v-bind="props">
              <v-icon>mdi-dots-vertical</v-icon>
            </v-btn>
          </template>
          <v-list>
            <v-list-item @click="store.deleteDocument(doc.id)">
              <template v-slot:prepend>
                <v-icon color="error">mdi-delete</v-icon>
              </template>
              <v-list-item-title class="text-error">Send to Trash</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>
      </v-card-actions>
    </v-card>

    <ReferenceDocFormDialog
      v-model="showRefForm"
      :edit-data="editingRefDoc"
      @saved="onSaved"
    />

    <TrashDialog
      v-model="showTrash"
      @restored="store.fetchDocuments()"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useDocumentStore } from '../stores/documents'
import { useDocumentCategoryStore } from '../stores/documentCategories'
import ReferenceDocFormDialog from '../components/ReferenceDocFormDialog.vue'
import TrashDialog from '../components/TrashDialog.vue'
import type { Document } from '../types'

const route = useRoute()
const router = useRouter()
const store = useDocumentStore()
const categoryStore = useDocumentCategoryStore()

const showRefForm = ref(false)
const showTrash = ref(false)
const editingRefDoc = ref<Document | null>(null)

function getCategoryName(categoryId: number) {
  const cat = categoryStore.categories.find((c) => c.id === categoryId)
  return cat?.name || ''
}

function startEdit(doc: Document) {
  if (doc.doc_type === 'typed') {
    router.push({ name: 'document-edit', params: { id: doc.id } })
  } else {
    editingRefDoc.value = doc
    showRefForm.value = true
  }
}

function onSaved() {
  editingRefDoc.value = null
  store.fetchDocuments()
}

watch(() => route.query.edit, async (editId) => {
  if (editId) {
    const doc = await store.fetchDocument(Number(editId))
    if (doc.doc_type === 'reference') {
      editingRefDoc.value = doc
      showRefForm.value = true
    }
    router.replace({ name: 'documents' })
  }
})

onMounted(() => {
  store.fetchDocuments()
  categoryStore.fetchCategories()
})
</script>
