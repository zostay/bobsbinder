<template>
  <div class="print-template">
    <h2>Obituary Information</h2>
    <template v-for="party in parties" :key="party.name">
      <h3>{{ party.name }} ({{ party.relationship }})</h3>
      <template v-for="(entries, groupType) in groupedEntries(party.entries)" :key="groupType">
        <h4 style="font-size: 12pt; margin: 0.8em 0 0.3em; text-transform: capitalize;">{{ groupType }}s</h4>
        <div v-for="entry in entries" :key="entry.id" class="list-item">
          &bull; <strong>{{ entry.name }}</strong>
          <span v-if="entry.relationship"> ({{ entry.relationship }})</span>
          <span v-if="entry.details"> &mdash; {{ entry.details }}</span>
          <span v-if="entry.event_date"> [{{ entry.event_date }}]</span>
        </div>
      </template>
      <p v-if="!party.entries.length" class="text-medium-emphasis">No entries yet.</p>
    </template>
  </div>
</template>

<script setup lang="ts">
import type { PartyObituaryInfo } from '../types'

defineProps<{
  parties: { name: string; relationship: string; entries: PartyObituaryInfo[] }[]
}>()

function groupedEntries(entries: PartyObituaryInfo[]) {
  const groups: Record<string, PartyObituaryInfo[]> = {}
  for (const e of entries) {
    if (!groups[e.type]) groups[e.type] = []
    groups[e.type].push(e)
  }
  return groups
}
</script>
