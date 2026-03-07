<template>
  <v-card v-if="!completed" class="mb-6">
    <v-card-title class="d-flex align-center">
      <span>Getting Started</span>
      <v-spacer />
      <v-btn variant="text" size="small" @click="dismiss">Skip for now</v-btn>
    </v-card-title>

    <v-stepper v-model="step" alt-labels>
      <v-stepper-header>
        <v-stepper-item title="Parties" :value="1" />
        <v-divider />
        <v-stepper-item title="Contacts" :value="2" />
        <v-divider />
        <v-stepper-item title="Locations" :value="3" />
        <v-divider />
        <v-stepper-item title="Documents" :value="4" />
        <v-divider />
        <v-stepper-item title="Obituary" :value="5" />
      </v-stepper-header>

      <v-stepper-window>
        <!-- Step 1: Parties -->
        <v-stepper-window-item :value="1">
          <div class="pa-4">
            <p class="text-body-1 mb-4">
              Parties represent the people your binder is organized for. Start with yourself,
              then add your spouse or dependents if applicable.
            </p>
            <v-list v-if="partyStore.parties.length" density="compact" class="mb-3">
              <v-list-item v-for="party in partyStore.parties" :key="party.id">
                <v-list-item-title>{{ party.name }}</v-list-item-title>
                <v-list-item-subtitle>{{ party.relationship }}</v-list-item-subtitle>
              </v-list-item>
            </v-list>
            <v-btn color="primary" variant="outlined" @click="showPartyDialog = true">
              Add Party
            </v-btn>
          </div>
        </v-stepper-window-item>

        <!-- Step 2: Primary Contacts -->
        <v-stepper-window-item :value="2">
          <div class="pa-4">
            <p class="text-body-1 mb-4">
              Add the people who should be contacted immediately. These are your key contacts
              like your attorney, pastor, financial advisor, or close family members.
            </p>
            <v-list v-if="primaryContacts.length" density="compact" class="mb-3">
              <v-list-item v-for="contact in primaryContacts" :key="contact.id">
                <v-list-item-title>{{ contact.name }}</v-list-item-title>
                <v-list-item-subtitle>{{ contact.role || contact.relationship }}</v-list-item-subtitle>
              </v-list-item>
            </v-list>
            <v-btn color="primary" variant="outlined" @click="showContactDialog = true">
              Add Primary Contact
            </v-btn>
          </div>
        </v-stepper-window-item>

        <!-- Step 3: Locations & Digital Access -->
        <v-stepper-window-item :value="3">
          <div class="pa-4">
            <p class="text-body-1 mb-4">
              Where can your loved ones find important things? Add physical locations
              (safe, filing cabinet) and digital access info (computers, password managers).
            </p>
            <v-row>
              <v-col cols="12" md="6">
                <h4 class="text-subtitle-1 mb-2">Locations</h4>
                <v-list v-if="locationStore.locations.length" density="compact" class="mb-3">
                  <v-list-item v-for="loc in locationStore.locations" :key="loc.id">
                    <v-list-item-title>{{ loc.name }}</v-list-item-title>
                    <v-list-item-subtitle>{{ loc.type }}</v-list-item-subtitle>
                  </v-list-item>
                </v-list>
                <v-btn color="primary" variant="outlined" @click="showLocationDialog = true">
                  Add Location
                </v-btn>
              </v-col>
              <v-col cols="12" md="6">
                <h4 class="text-subtitle-1 mb-2">Digital Access</h4>
                <v-list v-if="digitalStore.items.length" density="compact" class="mb-3">
                  <v-list-item v-for="item in digitalStore.items" :key="item.id">
                    <v-list-item-title>{{ item.name }}</v-list-item-title>
                    <v-list-item-subtitle>{{ item.type }}</v-list-item-subtitle>
                  </v-list-item>
                </v-list>
                <v-btn color="primary" variant="outlined" @click="showDigitalDialog = true">
                  Add Digital Access
                </v-btn>
              </v-col>
            </v-row>
          </div>
        </v-stepper-window-item>

        <!-- Step 4: Key Documents -->
        <v-stepper-window-item :value="4">
          <div class="pa-4">
            <p class="text-body-1 mb-4">
              These are the most important documents to have in your binder. Add them now
              or mark categories as not applicable.
            </p>
            <v-select
              v-if="partyStore.parties.length > 1"
              v-model="selectedPartyId"
              :items="partyStore.parties"
              item-title="name"
              item-value="id"
              label="Party"
              density="compact"
              class="mb-3"
              style="max-width: 300px"
            />
            <v-list density="compact" class="mb-3">
              <v-list-item v-for="cat in keyCategories" :key="cat.id">
                <template v-slot:prepend>
                  <v-icon :color="categoryStatusColor(cat)">
                    {{ categoryStatusIcon(cat) }}
                  </v-icon>
                </template>
                <v-list-item-title>{{ cat.name }}</v-list-item-title>
                <template v-slot:append>
                  <v-btn
                    size="small"
                    variant="text"
                    color="primary"
                    @click="openDocDialogForCategory(cat.id)"
                  >Add</v-btn>
                  <v-btn
                    size="small"
                    variant="text"
                    @click="markNA(cat.id)"
                  >N/A</v-btn>
                </template>
              </v-list-item>
            </v-list>
          </div>
        </v-stepper-window-item>

        <!-- Step 5: Obituary Information -->
        <v-stepper-window-item :value="5">
          <div class="pa-4">
            <p class="text-body-1 mb-4">
              Add biographical details for each party. This information will be used to help
              prepare obituary notices and memorial information.
            </p>
            <v-select
              v-if="partyStore.parties.length > 1"
              v-model="obituaryPartyId"
              :items="partyStore.parties"
              item-title="name"
              item-value="id"
              label="Party"
              density="compact"
              class="mb-3"
              style="max-width: 300px"
            />
            <v-list v-if="obituaryStore.items.length" density="compact" class="mb-3">
              <v-list-item v-for="item in obituaryStore.items" :key="item.id">
                <v-list-item-title>{{ item.name }}</v-list-item-title>
                <v-list-item-subtitle>{{ item.type }}{{ item.relationship ? ` - ${item.relationship}` : '' }}</v-list-item-subtitle>
              </v-list-item>
            </v-list>
            <v-btn color="primary" variant="outlined" @click="openObituaryDialog">
              Add Obituary Entry
            </v-btn>
          </div>
        </v-stepper-window-item>
      </v-stepper-window>

      <div class="d-flex justify-space-between pa-4">
        <v-btn v-if="step > 1" variant="text" @click="step--">Back</v-btn>
        <v-spacer />
        <v-btn v-if="step < 5" color="primary" @click="step++">Next</v-btn>
        <v-btn v-else color="primary" @click="finish">Finish</v-btn>
      </div>
    </v-stepper>

    <PartyFormDialog v-model="showPartyDialog" @saved="partyStore.fetchParties()" />
    <ContactFormDialog v-model="showContactDialog" :defaults="{ is_primary: true }" @saved="contactStore.fetchContacts()" />
    <LocationFormDialog v-model="showLocationDialog" @saved="locationStore.fetchLocations()" />
    <DigitalInfoFormDialog v-model="showDigitalDialog" @saved="digitalStore.fetchItems()" />
    <DocumentFormDialog
      v-model="showDocDialog"
      :initial-category-id="docCategoryId"
      :party-id="selectedPartyId"
      @saved="onDocSaved"
    />
    <DocumentFormDialog
      v-model="showObituaryDialog"
      :initial-type="'obituary_entry'"
      :party-id="obituaryPartyId"
      @saved="onObituarySaved"
    />
  </v-card>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { usePartyStore } from '../stores/parties'
import { useContactStore } from '../stores/contacts'
import { useLocationStore } from '../stores/locations'
import { useDigitalAccessStore } from '../stores/digitalAccess'
import { useChecklistStore } from '../stores/checklist'
import { useObituaryInfoStore } from '../stores/obituaryInfo'
import { useDocumentCategoryStore } from '../stores/documentCategories'
import PartyFormDialog from './PartyFormDialog.vue'
import ContactFormDialog from './ContactFormDialog.vue'
import LocationFormDialog from './LocationFormDialog.vue'
import DigitalInfoFormDialog from './DigitalInfoFormDialog.vue'
import DocumentFormDialog from './DocumentFormDialog.vue'

const STORAGE_KEY = 'bobsbinder_stepper_completed'

const partyStore = usePartyStore()
const contactStore = useContactStore()
const locationStore = useLocationStore()
const digitalStore = useDigitalAccessStore()
const checklistStore = useChecklistStore()
const obituaryStore = useObituaryInfoStore()
const categoryStore = useDocumentCategoryStore()

const step = ref(1)
const completed = ref(localStorage.getItem(STORAGE_KEY) === 'true')

const showPartyDialog = ref(false)
const showContactDialog = ref(false)
const showLocationDialog = ref(false)
const showDigitalDialog = ref(false)
const showDocDialog = ref(false)
const showObituaryDialog = ref(false)
const docCategoryId = ref<number | undefined>(undefined)

const selectedPartyId = ref<number | undefined>(undefined)
const obituaryPartyId = ref<number | undefined>(undefined)

const keyCategorySlugs = ['will', 'poa-medical', 'poa-financial', 'medical-directives', 'final-arrangements']

const keyCategories = computed(() =>
  categoryStore.categories.filter((c) => keyCategorySlugs.includes(c.slug)),
)

const primaryContacts = computed(() =>
  contactStore.contacts.filter((c) => c.is_primary),
)

const currentChecklist = computed(() =>
  checklistStore.checklists.find((c) => c.party_id === selectedPartyId.value),
)

watch(() => partyStore.selfParty, (self) => {
  if (self && !selectedPartyId.value) {
    selectedPartyId.value = self.id
    obituaryPartyId.value = self.id
  }
})

watch(obituaryPartyId, (id) => {
  if (id) obituaryStore.fetchItems(id)
})

function categoryStatusIcon(cat: { id: number }) {
  const item = currentChecklist.value?.items.find((i) => i.category_id === cat.id)
  if (item?.has_document) return 'mdi-check-circle'
  if (item?.status === 'not_applicable') return 'mdi-minus-circle'
  return 'mdi-circle-outline'
}

function categoryStatusColor(cat: { id: number }) {
  const item = currentChecklist.value?.items.find((i) => i.category_id === cat.id)
  if (item?.has_document) return 'success'
  if (item?.status === 'not_applicable') return 'grey'
  return 'grey-lighten-1'
}

function openDocDialogForCategory(categoryId: number) {
  docCategoryId.value = categoryId
  showDocDialog.value = true
}

async function markNA(categoryId: number) {
  if (selectedPartyId.value) {
    await checklistStore.updateStatus(selectedPartyId.value, categoryId, 'not_applicable')
  }
}

function openObituaryDialog() {
  showObituaryDialog.value = true
}

async function onDocSaved() {
  await checklistStore.fetchAll()
}

async function onObituarySaved() {
  if (obituaryPartyId.value) {
    await obituaryStore.fetchItems(obituaryPartyId.value)
  }
}

function dismiss() {
  completed.value = true
}

function finish() {
  completed.value = true
  localStorage.setItem(STORAGE_KEY, 'true')
}

onMounted(async () => {
  await Promise.all([
    partyStore.fetchParties(),
    contactStore.fetchContacts(),
    locationStore.fetchLocations(),
    digitalStore.fetchItems(),
    categoryStore.fetchCategories(),
    checklistStore.fetchAll(),
  ])
  if (partyStore.selfParty) {
    selectedPartyId.value = partyStore.selfParty.id
    obituaryPartyId.value = partyStore.selfParty.id
    obituaryStore.fetchItems(partyStore.selfParty.id)
  }
})
</script>
