<template>
  <div class="confidential-print">
    <h2 class="text-center mb-1">CONFIDENTIAL &mdash; For Authorized Eyes Only</h2>
    <p class="text-center text-caption mb-4">
      This document contains sensitive information. Keep it in a secure location.
    </p>

    <template v-for="section in sections" :key="section.title">
      <div v-if="section.items.length > 0" class="mb-4">
        <h3 class="mb-2">{{ section.title }}</h3>
        <div v-for="item in section.items" :key="item.label" class="mb-2 ml-2">
          <strong>{{ item.label }}</strong>
          <div v-for="field in item.fields" :key="field.name" class="ml-4">
            <span class="text-medium-emphasis">{{ field.name }}:</span>
            <span class="pre-wrap"> {{ field.value }}</span>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import type { ConfidentialSection } from '../types'

defineProps<{
  sections: ConfidentialSection[]
}>()
</script>

<style scoped>
.confidential-print {
  font-family: Georgia, 'Times New Roman', Times, serif;
  font-size: 11pt;
  line-height: 1.4;
}

.pre-wrap {
  white-space: pre-wrap;
}
</style>
