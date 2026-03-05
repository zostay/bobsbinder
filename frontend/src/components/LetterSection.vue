<template>
  <v-card class="mb-4" :variant="section.visible ? 'elevated' : 'outlined'">
    <v-card-title class="d-flex align-center">
      <v-icon size="small" class="mr-2 cursor-grab">mdi-drag-horizontal</v-icon>

      <template v-if="editingTitle">
        <v-text-field
          v-model="titleText"
          density="compact"
          variant="outlined"
          hide-details
          class="flex-grow-1"
          @keyup.enter="saveTitle"
          @keyup.escape="cancelTitle"
        />
        <v-btn icon size="small" color="success" @click="saveTitle" class="ml-1">
          <v-icon size="small">mdi-check</v-icon>
        </v-btn>
        <v-btn icon size="small" @click="cancelTitle" class="ml-1">
          <v-icon size="small">mdi-close</v-icon>
        </v-btn>
      </template>
      <template v-else>
        <span class="flex-grow-1" :class="{ 'text-grey': !section.visible }">{{ section.title }}</span>
        <v-btn icon size="small" variant="text" @click="startTitleEdit">
          <v-icon size="small">mdi-pencil</v-icon>
        </v-btn>
      </template>

      <v-btn icon size="small" variant="text" @click="toggleVisibility">
        <v-icon size="small">{{ section.visible ? 'mdi-eye' : 'mdi-eye-off' }}</v-icon>
      </v-btn>
    </v-card-title>

    <v-card-text v-if="section.visible">
      <div v-if="visibleItems.length === 0 && suppressedItems.length === 0" class="text-grey text-body-2">
        No items in this section yet.
      </div>

      <LetterItem
        v-for="(item, idx) in visibleItems"
        :key="item.id"
        :item="item"
        :index="idx"
        @edit="(id, content) => $emit('editItem', id, content)"
        @delete="(id) => $emit('deleteItem', id)"
        @unsuppress="(id) => $emit('unsuppressItem', id)"
      />

      <v-divider v-if="suppressedItems.length > 0" class="my-2" />
      <div v-if="suppressedItems.length > 0" class="text-caption text-grey mb-1">Suppressed items:</div>
      <LetterItem
        v-for="(item, idx) in suppressedItems"
        :key="item.id"
        :item="item"
        :index="idx"
        @edit="(id, content) => $emit('editItem', id, content)"
        @delete="(id) => $emit('deleteItem', id)"
        @unsuppress="(id) => $emit('unsuppressItem', id)"
      />

      <v-btn
        size="small" variant="text" color="primary" prepend-icon="mdi-plus"
        class="mt-2" @click="showAddItem = true"
      >
        Add Item
      </v-btn>

      <div v-if="showAddItem" class="d-flex align-center mt-2">
        <v-text-field
          v-model="newItemContent"
          density="compact"
          variant="outlined"
          hide-details
          placeholder="Enter item text..."
          class="flex-grow-1"
          @keyup.enter="addItem"
          @keyup.escape="showAddItem = false"
        />
        <v-btn icon size="small" color="success" @click="addItem" class="ml-1">
          <v-icon size="small">mdi-check</v-icon>
        </v-btn>
        <v-btn icon size="small" @click="showAddItem = false" class="ml-1">
          <v-icon size="small">mdi-close</v-icon>
        </v-btn>
      </div>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { SurvivorLetterSection } from '../types'
import LetterItem from './LetterItem.vue'

const props = defineProps<{
  section: SurvivorLetterSection
}>()

const emit = defineEmits<{
  updateSection: [sectionId: number, updates: { title?: string; visible?: boolean }]
  addItem: [sectionId: number, content: string]
  editItem: [itemId: number, content: string]
  deleteItem: [itemId: number]
  unsuppressItem: [itemId: number]
}>()

const editingTitle = ref(false)
const titleText = ref('')
const showAddItem = ref(false)
const newItemContent = ref('')

const visibleItems = computed(() => props.section.items.filter((i) => !i.suppressed))
const suppressedItems = computed(() => props.section.items.filter((i) => i.suppressed))

function startTitleEdit() {
  titleText.value = props.section.title
  editingTitle.value = true
}

function saveTitle() {
  emit('updateSection', props.section.id, { title: titleText.value })
  editingTitle.value = false
}

function cancelTitle() {
  editingTitle.value = false
}

function toggleVisibility() {
  emit('updateSection', props.section.id, { visible: !props.section.visible })
}

function addItem() {
  if (newItemContent.value.trim()) {
    emit('addItem', props.section.id, newItemContent.value.trim())
    newItemContent.value = ''
    showAddItem.value = false
  }
}
</script>
