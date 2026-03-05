<template>
  <div>
    <div class="d-flex align-center mb-4">
      <h1 class="text-h4">Survivor Letter</h1>
    </div>

    <v-progress-linear v-if="store.loading" indeterminate color="primary" />

    <template v-if="store.letter">
      <v-card class="mb-4">
        <v-card-title>Letter Boilerplate</v-card-title>
        <v-card-text>
          <v-textarea v-model="boilerplate.greeting" label="Greeting" rows="2" />
          <v-textarea v-model="boilerplate.intro" label="Introduction" rows="4" />
          <v-textarea v-model="boilerplate.closing" label="Closing" rows="3" />
          <v-textarea v-model="boilerplate.signature" label="Signature" rows="2" />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn color="primary" @click="saveBoilerplate" :loading="savingBoilerplate">
            Save Boilerplate
          </v-btn>
        </v-card-actions>
      </v-card>

      <h2 class="text-h5 mb-3">Sections</h2>

      <LetterSection
        v-for="section in store.letter.sections"
        :key="section.id"
        :section="section"
        @update-section="handleUpdateSection"
        @add-item="handleAddItem"
        @edit-item="handleEditItem"
        @delete-item="handleDeleteItem"
        @unsuppress-item="handleUnsuppressItem"
      />
    </template>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, watch } from 'vue'
import { useSurvivorLetterStore } from '../stores/survivorLetter'
import LetterSection from '../components/LetterSection.vue'

const store = useSurvivorLetterStore()
const savingBoilerplate = ref(false)

const boilerplate = reactive({
  greeting: '',
  intro: '',
  closing: '',
  signature: '',
})

watch(
  () => store.letter,
  (letter) => {
    if (letter) {
      boilerplate.greeting = letter.greeting
      boilerplate.intro = letter.intro
      boilerplate.closing = letter.closing
      boilerplate.signature = letter.signature
    }
  },
  { immediate: true },
)

async function saveBoilerplate() {
  savingBoilerplate.value = true
  try {
    await store.updateBoilerplate({ ...boilerplate })
  } finally {
    savingBoilerplate.value = false
  }
}

async function handleUpdateSection(sectionId: number, updates: { title?: string; visible?: boolean }) {
  await store.updateSection(sectionId, updates)
}

async function handleAddItem(sectionId: number, content: string) {
  await store.addItem(sectionId, content)
}

async function handleEditItem(itemId: number, content: string) {
  await store.editItem(itemId, content)
}

async function handleDeleteItem(itemId: number) {
  await store.deleteItem(itemId)
}

async function handleUnsuppressItem(itemId: number) {
  await store.unsuppressItem(itemId)
}

onMounted(() => {
  store.fetchLetter()
})
</script>
