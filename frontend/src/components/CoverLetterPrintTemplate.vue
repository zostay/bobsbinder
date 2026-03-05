<template>
  <div class="print-template">
    <p v-if="greeting">{{ greeting }}</p>

    <p v-if="intro" class="pre-wrap">{{ intro }}</p>

    <template v-for="section in previewSections" :key="section.id">
      <h3>{{ section.title }}</h3>
      <template v-for="item in activeItems(section)" :key="item.id">
        <p v-if="item.item_type === 'paragraph'">{{ item.content }}</p>
        <p v-else-if="item.item_type === 'bulleted'" class="list-item">&bull; {{ item.content }}</p>
        <p v-else class="list-item">{{ numberedIndex(section, item) }}. {{ item.content }}</p>
      </template>
    </template>

    <p v-if="closing" class="pre-wrap">{{ closing }}</p>

    <p v-if="signature">{{ signature }}</p>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { FullSurvivorLetter, SurvivorLetterSection, SurvivorLetterItem } from '../types'

const props = defineProps<{
  letter: FullSurvivorLetter
  greeting: string
  intro: string
  closing: string
  signature: string
}>()

const previewSections = computed(() =>
  props.letter.sections.filter(
    (s) => s.visible && s.items.some((i) => !i.suppressed),
  ),
)

function activeItems(section: SurvivorLetterSection): SurvivorLetterItem[] {
  return section.items
    .filter((i) => !i.suppressed)
    .sort((a, b) => a.sort_order - b.sort_order)
}

function numberedIndex(section: SurvivorLetterSection, item: SurvivorLetterItem): number {
  const items = activeItems(section)
  let counter = 0
  for (const i of items) {
    if (i.item_type === 'numbered') counter++
    if (i.id === item.id) return counter
  }
  return counter
}
</script>

<style>
.print-template {
  font-family: Georgia, 'Times New Roman', Times, serif;
  font-size: 12pt;
  line-height: 1.5;
  color: #000;
  background: #fff;
}

.print-template h3 {
  font-size: 14pt;
  margin: 1em 0 0.5em;
}

.print-template p {
  margin: 0.5em 0;
}

.print-template .pre-wrap {
  white-space: pre-wrap;
}

.print-template .list-item {
  padding-left: 1.5em;
  text-indent: -1.5em;
}
</style>
