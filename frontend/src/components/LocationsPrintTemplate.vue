<template>
  <div class="print-template">
    <h2>Important Locations</h2>
    <div v-for="loc in locations" :key="loc.id" class="mb-3">
      <p>
        <strong>{{ loc.name }}</strong>
        <span class="text-medium-emphasis"> ({{ loc.type }})</span>
      </p>
      <p v-if="loc.description" class="ml-4">{{ loc.description }}</p>
      <p v-if="loc.address" class="ml-4">Address: {{ loc.address }}</p>
    </div>
    <p v-if="!locations.length">No locations recorded.</p>
  </div>
</template>

<script setup lang="ts">
import type { Location } from '../types'

defineProps<{
  locations: Location[]
}>()
</script>
