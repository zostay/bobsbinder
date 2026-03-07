<template>
  <v-card class="mb-6">
    <v-card-title>Printable Documents</v-card-title>
    <v-card-text>
      <v-expansion-panels variant="accordion">
        <!-- 1. Cover Letter -->
        <v-expansion-panel title="Cover Letter" @group:selected="loadCoverLetter">
          <v-expansion-panel-text>
            <template v-if="letterStore.letter">
              <DocumentPreviewPanel title="Cover Letter">
                <CoverLetterPrintTemplate
                  :letter="letterStore.letter"
                  :greeting="letterStore.letter.greeting"
                  :intro="letterStore.letter.intro"
                  :closing="letterStore.letter.closing"
                  :signature="letterStore.letter.signature"
                />
              </DocumentPreviewPanel>
            </template>
            <v-progress-linear v-else indeterminate />
          </v-expansion-panel-text>
        </v-expansion-panel>

        <!-- 2. Confidential Insert -->
        <v-expansion-panel title="Confidential Insert" @group:selected="loadConfidential">
          <v-expansion-panel-text>
            <template v-if="confidentialSections">
              <DocumentPreviewPanel title="Confidential Insert">
                <ConfidentialPrintTemplate :sections="confidentialSections" />
              </DocumentPreviewPanel>
            </template>
            <v-progress-linear v-else-if="confidentialLoading" indeterminate />
            <p v-else class="text-medium-emphasis">No confidential data.</p>
          </v-expansion-panel-text>
        </v-expansion-panel>

        <!-- 3. Obituary Information -->
        <v-expansion-panel title="Obituary Information" @group:selected="loadObituary">
          <v-expansion-panel-text>
            <template v-if="obituaryParties.length">
              <DocumentPreviewPanel title="Obituary Information">
                <ObituaryPrintTemplate :parties="obituaryParties" />
              </DocumentPreviewPanel>
            </template>
            <v-progress-linear v-else-if="obituaryLoading" indeterminate />
            <p v-else class="text-medium-emphasis">No obituary data.</p>
          </v-expansion-panel-text>
        </v-expansion-panel>

        <!-- 4. Individual Documents -->
        <v-expansion-panel title="Individual Documents" @group:selected="loadDocuments">
          <v-expansion-panel-text>
            <template v-if="documentsLoaded">
              <DocumentPreviewPanel title="Individual Documents">
                <DocumentsPrintTemplate :documents="documentStore.documents" :categories="categoryStore.categories" />
              </DocumentPreviewPanel>
            </template>
            <v-progress-linear v-else indeterminate />
          </v-expansion-panel-text>
        </v-expansion-panel>

        <!-- 5. External Documents & Locations -->
        <v-expansion-panel title="External Documents &amp; Locations" @group:selected="loadDocuments">
          <v-expansion-panel-text>
            <template v-if="documentsLoaded">
              <DocumentPreviewPanel title="External Documents & Locations">
                <ExternalDocsPrintTemplate :documents="documentStore.documents" :locations="locationStore.locations" />
              </DocumentPreviewPanel>
            </template>
            <v-progress-linear v-else indeterminate />
          </v-expansion-panel-text>
        </v-expansion-panel>

        <!-- 6. All Locations -->
        <v-expansion-panel title="All Locations" @group:selected="loadLocations">
          <v-expansion-panel-text>
            <DocumentPreviewPanel title="All Locations">
              <LocationsPrintTemplate :locations="locationStore.locations" />
            </DocumentPreviewPanel>
          </v-expansion-panel-text>
        </v-expansion-panel>

        <!-- 7. Contact Page -->
        <v-expansion-panel title="Contact Page" @group:selected="loadContacts">
          <v-expansion-panel-text>
            <DocumentPreviewPanel title="Contacts">
              <ContactsPrintTemplate :contacts="contactStore.contacts" />
            </DocumentPreviewPanel>
          </v-expansion-panel-text>
        </v-expansion-panel>
      </v-expansion-panels>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useSurvivorLetterStore } from '../stores/survivorLetter'
import { useDocumentStore } from '../stores/documents'
import { useDocumentCategoryStore } from '../stores/documentCategories'
import { useLocationStore } from '../stores/locations'
import { useContactStore } from '../stores/contacts'
import { usePartyStore } from '../stores/parties'
import { useObituaryInfoStore } from '../stores/obituaryInfo'
import api from '../services/api'
import type { ConfidentialSection, PartyObituaryInfo } from '../types'
import DocumentPreviewPanel from './DocumentPreviewPanel.vue'
import CoverLetterPrintTemplate from './CoverLetterPrintTemplate.vue'
import ConfidentialPrintTemplate from './ConfidentialPrintTemplate.vue'
import ObituaryPrintTemplate from './ObituaryPrintTemplate.vue'
import DocumentsPrintTemplate from './DocumentsPrintTemplate.vue'
import ExternalDocsPrintTemplate from './ExternalDocsPrintTemplate.vue'
import LocationsPrintTemplate from './LocationsPrintTemplate.vue'
import ContactsPrintTemplate from './ContactsPrintTemplate.vue'

const letterStore = useSurvivorLetterStore()
const documentStore = useDocumentStore()
const categoryStore = useDocumentCategoryStore()
const locationStore = useLocationStore()
const contactStore = useContactStore()
const partyStore = usePartyStore()
const obituaryStore = useObituaryInfoStore()

const confidentialSections = ref<ConfidentialSection[] | null>(null)
const confidentialLoading = ref(false)
const obituaryParties = ref<{ name: string; relationship: string; entries: PartyObituaryInfo[] }[]>([])
const obituaryLoading = ref(false)
const documentsLoaded = ref(false)

let coverLetterLoaded = false
let confidentialLoaded = false
let obituaryLoaded = false
let docsLoaded = false

async function loadCoverLetter() {
  if (coverLetterLoaded) return
  coverLetterLoaded = true
  await letterStore.fetchLetter()
}

async function loadConfidential() {
  if (confidentialLoaded) return
  confidentialLoaded = true
  confidentialLoading.value = true
  try {
    const { data } = await api.get<ConfidentialSection[]>('/confidential')
    confidentialSections.value = data
  } finally {
    confidentialLoading.value = false
  }
}

async function loadObituary() {
  if (obituaryLoaded) return
  obituaryLoaded = true
  obituaryLoading.value = true
  try {
    await partyStore.fetchParties()
    const results: { name: string; relationship: string; entries: PartyObituaryInfo[] }[] = []
    for (const party of partyStore.parties) {
      const { data } = await api.get<PartyObituaryInfo[]>(`/parties/${party.id}/obituary-info`)
      results.push({ name: party.name, relationship: party.relationship, entries: data })
    }
    obituaryParties.value = results
  } finally {
    obituaryLoading.value = false
  }
}

async function loadDocuments() {
  if (docsLoaded) return
  docsLoaded = true
  await Promise.all([
    documentStore.fetchDocuments(),
    categoryStore.fetchCategories(),
    locationStore.fetchLocations(),
  ])
  documentsLoaded.value = true
}

async function loadLocations() {
  await locationStore.fetchLocations()
}

async function loadContacts() {
  await contactStore.fetchContacts()
}
</script>
