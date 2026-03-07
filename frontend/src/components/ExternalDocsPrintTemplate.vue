<template>
  <div class="print-template">
    <h2>External Documents &amp; Locations</h2>
    <div v-for="doc in externalDocs" :key="doc.id" class="mb-3">
      <p>
        <strong>{{ doc.title }}</strong>
        <span v-if="locationFor(doc)" class="text-medium-emphasis">
          &mdash; located at {{ locationFor(doc)!.name }}
          <span v-if="locationFor(doc)!.address"> ({{ locationFor(doc)!.address }})</span>
        </span>
      </p>
    </div>
    <p v-if="!externalDocs.length">No external documents recorded.</p>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Document, Location } from '../types'

const props = defineProps<{
  documents: Document[]
  locations: Location[]
}>()

const externalDocs = computed(() =>
  props.documents.filter((d) => d.location_id),
)

function locationFor(doc: Document): Location | undefined {
  return props.locations.find((l) => l.id === doc.location_id)
}
</script>
