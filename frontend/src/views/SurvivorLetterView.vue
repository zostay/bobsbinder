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
          @add-structured="handleAddStructured"
          @edit-structured="handleEditStructured"
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

        <v-btn v-if="!showConfidentialPreview" color="warning" variant="tonal" prepend-icon="mdi-lock"
          class="mt-3" block @click="fetchConfidential">
          Print Confidential Supplement
        </v-btn>
        <v-btn v-else color="grey" variant="tonal" prepend-icon="mdi-eye-off"
          class="mt-3" block @click="hideConfidential">
          Hide Confidential Supplement
        </v-btn>

        <DocumentPreviewPanel v-if="showConfidentialPreview && confidentialSections.length > 0"
          title="Confidential Supplement" class="mt-3">
          <ConfidentialPrintTemplate :sections="confidentialSections" />
        </DocumentPreviewPanel>
      </v-col>
    </v-row>

    <ContactFormDialog v-model="showContactDialog" :edit-data="editingContact" @saved="handleStructuredSaved" />
    <ReferenceDocFormDialog v-model="showDocumentDialog" :edit-data="editingDocument" @saved="handleStructuredSaved" />
    <LocationFormDialog v-model="showLocationDialog" :edit-data="editingLocation" @saved="handleStructuredSaved" />
    <DigitalInfoFormDialog v-model="showDigitalInfoDialog" :edit-data="editingDigitalInfo" @saved="handleStructuredSaved" />
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useSurvivorLetterStore } from '../stores/survivorLetter'
import LetterSection from '../components/LetterSection.vue'
import DocumentPreviewPanel from '../components/DocumentPreviewPanel.vue'
import CoverLetterPrintTemplate from '../components/CoverLetterPrintTemplate.vue'
import ConfidentialPrintTemplate from '../components/ConfidentialPrintTemplate.vue'
import ContactFormDialog from '../components/ContactFormDialog.vue'
import ReferenceDocFormDialog from '../components/ReferenceDocFormDialog.vue'
import LocationFormDialog from '../components/LocationFormDialog.vue'
import DigitalInfoFormDialog from '../components/DigitalInfoFormDialog.vue'
import api from '../services/api'
import type { Contact, Document, Location, DigitalAccess, ServiceAccount, ConfidentialSection } from '../types'

const router = useRouter()
const store = useSurvivorLetterStore()

const showContactDialog = ref(false)
const showDocumentDialog = ref(false)
const showLocationDialog = ref(false)
const showDigitalInfoDialog = ref(false)

const editingContact = ref<Contact | null>(null)
const editingDocument = ref<Document | null>(null)
const editingLocation = ref<Location | null>(null)
const editingDigitalInfo = ref<DigitalAccess | ServiceAccount | null>(null)

const confidentialSections = ref<ConfidentialSection[]>([])
const showConfidentialPreview = ref(false)

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

function handleAddStructured(sectionKey: string) {
  switch (sectionKey) {
    case 'contacts':
      showContactDialog.value = true
      break
    case 'documents':
      showDocumentDialog.value = true
      break
    case 'locations':
      showLocationDialog.value = true
      break
    case 'digital_info':
      showDigitalInfoDialog.value = true
      break
  }
}

const sourceTypeEndpoints: Record<string, string> = {
  contact: '/contacts',
  document: '/documents',
  location: '/locations',
  digital_access: '/digital-access',
  insurance_policy: '/insurance-policies',
  service_account: '/service-accounts',
}

async function handleEditStructured(sourceType: string, sourceId: number) {
  const endpoint = sourceTypeEndpoints[sourceType]
  if (!endpoint) return

  const { data } = await api.get(`${endpoint}/${sourceId}`)

  switch (sourceType) {
    case 'contact':
      editingContact.value = data as Contact
      showContactDialog.value = true
      break
    case 'document': {
      const doc = data as Document
      if (doc.doc_type === 'typed') {
        router.push({ name: 'document-edit', params: { id: doc.id } })
      } else {
        editingDocument.value = doc
        showDocumentDialog.value = true
      }
      break
    }
    case 'insurance_policy':
      editingDocument.value = data as Document
      showDocumentDialog.value = true
      break
    case 'location':
      editingLocation.value = data as Location
      showLocationDialog.value = true
      break
    case 'digital_access':
      editingDigitalInfo.value = data as DigitalAccess
      showDigitalInfoDialog.value = true
      break
    case 'service_account':
      editingDigitalInfo.value = data as ServiceAccount
      showDigitalInfoDialog.value = true
      break
  }
}

async function fetchConfidential() {
  const { data } = await api.get('/confidential')
  confidentialSections.value = data as ConfidentialSection[]
  showConfidentialPreview.value = true
}

function hideConfidential() {
  showConfidentialPreview.value = false
  confidentialSections.value = []
}

async function handleStructuredSaved() {
  editingContact.value = null
  editingDocument.value = null
  editingLocation.value = null
  editingDigitalInfo.value = null
  await store.fetchLetter()
}

onMounted(() => {
  store.fetchLetter()
})
</script>
