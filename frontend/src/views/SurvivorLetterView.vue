<template>
  <div>
    <div class="d-flex align-center mb-4">
      <h1 class="text-h4">Cover Letter</h1>
    </div>

    <v-progress-linear v-if="store.loading" indeterminate color="primary" />

    <v-row v-if="store.letter">
      <v-col cols="12" md="7" class="editor-column">
        <v-card class="mb-4">
          <v-card-text>
            <v-textarea v-model="boilerplate.greeting" label="Greeting" rows="2" />
          </v-card-text>
        </v-card>

        <v-card class="mb-4">
          <v-card-text>
            <v-textarea v-model="boilerplate.intro" label="Introduction" rows="4" />
          </v-card-text>
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

        <v-card class="mb-4">
          <v-card-text>
            <v-textarea v-model="boilerplate.closing" label="Closing" rows="3" />
          </v-card-text>
        </v-card>

        <v-card class="mb-4">
          <v-card-text>
            <v-textarea v-model="boilerplate.signature" label="Signature" rows="2" />
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" md="5">
        <DocumentPreviewPanel title="Cover Letter Preview">
          <CoverLetterPrintTemplate
            :letter="store.letter"
            :greeting="boilerplate.greeting"
            :intro="boilerplate.intro"
            :closing="boilerplate.closing"
            :signature="boilerplate.signature"
          />
        </DocumentPreviewPanel>
      </v-col>
    </v-row>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, watch } from 'vue'
import { useSurvivorLetterStore } from '../stores/survivorLetter'
import LetterSection from '../components/LetterSection.vue'
import DocumentPreviewPanel from '../components/DocumentPreviewPanel.vue'
import CoverLetterPrintTemplate from '../components/CoverLetterPrintTemplate.vue'

const store = useSurvivorLetterStore()

const boilerplate = reactive({
  greeting: '',
  intro: '',
  closing: '',
  signature: '',
})

let loaded = false
let debounceTimer: ReturnType<typeof setTimeout> | null = null

watch(
  () => store.letter,
  (letter) => {
    if (letter) {
      boilerplate.greeting = letter.greeting
      boilerplate.intro = letter.intro
      boilerplate.closing = letter.closing
      boilerplate.signature = letter.signature
      loaded = true
    }
  },
  { immediate: true },
)

watch(
  () => ({ ...boilerplate }),
  () => {
    if (!loaded) return

    if (debounceTimer) {
      clearTimeout(debounceTimer)
    }
    debounceTimer = setTimeout(() => {
      store.updateBoilerplate({ ...boilerplate })
    }, 1000)
  },
)

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
