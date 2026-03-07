<template>
  <div class="print-template">
    <h2>Documents</h2>
    <template v-for="cat in categoriesWithDocs" :key="cat.id">
      <h3>{{ cat.name }}</h3>
      <div v-for="doc in docsForCategory(cat.id)" :key="doc.id" class="mb-2">
        <p><strong>{{ doc.title }}</strong></p>
        <p v-if="doc.content" class="pre-wrap ml-4">{{ doc.content }}</p>
      </div>
    </template>
    <p v-if="!documents.length">No documents yet.</p>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Document, DocumentCategory } from '../types'

const props = defineProps<{
  documents: Document[]
  categories: DocumentCategory[]
}>()

const categoriesWithDocs = computed(() =>
  props.categories.filter((c) => props.documents.some((d) => d.category_id === c.id)),
)

function docsForCategory(categoryId: number) {
  return props.documents.filter((d) => d.category_id === categoryId)
}
</script>
