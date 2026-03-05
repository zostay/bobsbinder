<template>
  <div>
    <div class="d-flex align-center mb-4">
      <h1 class="text-h4">Documents</h1>
      <v-spacer />
      <v-btn color="primary" prepend-icon="mdi-plus">
        Add Document
      </v-btn>
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
      <v-card-title>{{ doc.title }}</v-card-title>
      <v-card-subtitle>
        <v-chip size="small" :color="doc.status === 'complete' ? 'success' : 'warning'">
          {{ doc.status }}
        </v-chip>
      </v-card-subtitle>
      <v-card-text v-if="doc.content">{{ doc.content }}</v-card-text>
      <v-card-actions>
        <v-btn size="small" variant="text" color="primary">Edit</v-btn>
        <v-btn size="small" variant="text" color="error" @click="store.deleteDocument(doc.id)">Delete</v-btn>
      </v-card-actions>
    </v-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useDocumentStore } from '../stores/documents'

const store = useDocumentStore()

onMounted(() => {
  store.fetchDocuments()
})
</script>
