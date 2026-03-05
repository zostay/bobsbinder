<template>
  <v-card class="preview-panel">
    <v-card-title class="d-flex align-center">
      <span>{{ title }}</span>
      <v-spacer />
      <v-btn icon size="small" variant="text" @click="printDocument">
        <v-icon>mdi-printer</v-icon>
      </v-btn>
    </v-card-title>

    <v-divider />

    <div class="preview-content" ref="previewContentRef">
      <slot />
    </div>

    <iframe ref="printFrameRef" class="print-frame" />
  </v-card>
</template>

<script setup lang="ts">
import { ref } from 'vue'

defineProps<{
  title: string
}>()

const previewContentRef = ref<HTMLDivElement>()
const printFrameRef = ref<HTMLIFrameElement>()

function printDocument() {
  const content = previewContentRef.value
  const frame = printFrameRef.value
  if (!content || !frame) return

  const doc = frame.contentDocument || frame.contentWindow?.document
  if (!doc) return

  doc.open()
  doc.write(`<!DOCTYPE html>
<html>
<head>
<style>
  body {
    font-family: Georgia, 'Times New Roman', Times, serif;
    font-size: 12pt;
    line-height: 1.5;
    color: #000;
    background: #fff;
    margin: 1in;
  }
  h3 {
    font-size: 14pt;
    margin: 1em 0 0.5em;
  }
  p {
    margin: 0.5em 0;
  }
  .pre-wrap {
    white-space: pre-wrap;
  }
</style>
</head>
<body>${content.innerHTML}</body>
</html>`)
  doc.close()

  frame.contentWindow?.focus()
  frame.contentWindow?.print()
}
</script>

<style scoped>
.preview-panel {
  position: sticky;
  top: 80px;
}

.preview-content {
  max-height: calc(100vh - 96px - 56px);
  overflow-y: auto;
  padding: 16px;
}

.print-frame {
  position: absolute;
  width: 0;
  height: 0;
  border: 0;
  visibility: hidden;
}
</style>
