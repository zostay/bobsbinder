<template>
  <v-card class="mb-6">
    <v-card-title>Binder Checklist</v-card-title>

    <v-tabs v-if="checklistStore.checklists.length > 1" v-model="activeTab" density="compact">
      <v-tab v-for="party in checklistStore.checklists" :key="party.party_id" :value="party.party_id">
        {{ party.party_name }}
      </v-tab>
    </v-tabs>

    <v-card-text>
      <v-list v-if="activeChecklist" density="compact">
        <v-list-item v-for="item in activeChecklist.items" :key="item.category_id">
          <template v-slot:prepend>
            <v-icon
              :color="statusColor(item)"
              @click="cycleStatus(item)"
              style="cursor: pointer"
            >
              {{ statusIcon(item) }}
            </v-icon>
          </template>
          <v-list-item-title>{{ item.category_name }}</v-list-item-title>
          <template v-slot:append>
            <v-btn
              size="small"
              variant="text"
              color="primary"
              :to="'/documents'"
            >Add Document</v-btn>
          </template>
        </v-list-item>
      </v-list>
      <p v-else class="text-body-2 text-medium-emphasis">
        No parties found. Add a party to get started.
      </p>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useChecklistStore } from '../stores/checklist'
import type { ChecklistItem } from '../types'

const checklistStore = useChecklistStore()
const activeTab = ref<number | undefined>(undefined)

const activeChecklist = computed(() =>
  checklistStore.checklists.find((c) => c.party_id === activeTab.value),
)

watch(() => checklistStore.checklists, (lists) => {
  if (lists.length && !activeTab.value) {
    activeTab.value = lists[0].party_id
  }
}, { immediate: true })

function statusIcon(item: ChecklistItem) {
  if (item.has_document) return 'mdi-check-circle'
  if (item.status === 'complete') return 'mdi-check-circle'
  if (item.status === 'not_applicable') return 'mdi-minus-circle'
  return 'mdi-circle-outline'
}

function statusColor(item: ChecklistItem) {
  if (item.has_document) return 'success'
  if (item.status === 'complete') return 'success'
  if (item.status === 'not_applicable') return 'grey'
  return 'grey-lighten-1'
}

async function cycleStatus(item: ChecklistItem) {
  if (item.has_document) return // complete via document, don't cycle
  const nextStatus = item.status === 'pending' ? 'not_applicable' : 'pending'
  if (activeTab.value) {
    await checklistStore.updateStatus(activeTab.value, item.category_id, nextStatus)
  }
}

onMounted(() => {
  checklistStore.fetchAll()
})
</script>
