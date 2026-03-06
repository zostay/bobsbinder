<template>
  <v-dialog :model-value="modelValue" max-width="600" @update:model-value="$emit('update:modelValue', $event)">
    <v-card>
      <v-card-title>{{ dialogTitle }}</v-card-title>
      <v-card-text>
        <v-alert type="info" variant="tonal" density="compact" class="mb-3">
          Fields marked with <v-icon size="small">mdi-lock</v-icon> are confidential
          and will not be sent to the printing service. Only place passwords, PINs,
          access codes, or account numbers in locked fields.
        </v-alert>
        <v-select
          v-if="!editData"
          v-model="documentType"
          label="Type"
          :items="typeOptions"
          class="mb-2"
        />

        <template v-if="documentType === 'document'">
          <v-text-field v-model="docForm.title" label="Title" required />
          <v-select
            v-model="docForm.category_id"
            label="Category"
            :items="categoryItems"
            item-title="name"
            item-value="id"
          />
          <v-textarea v-model="docForm.content" label="Content" rows="3" />
          <v-select
            v-model="docForm.status"
            label="Status"
            :items="['draft', 'complete']"
          />
          <v-textarea v-model="docForm.secure_notes" label="Confidential Notes"
            rows="2" prepend-inner-icon="mdi-lock"
            hint="Will NOT appear on the printed cover letter."
            persistent-hint />
        </template>

        <template v-else-if="documentType === 'insurance_policy'">
          <v-text-field v-model="policyForm.provider" label="Insurance Provider" required />
          <v-text-field v-model="policyForm.policy_number" label="Policy Number" prepend-inner-icon="mdi-lock" />
          <v-text-field v-model="policyForm.type" label="Type (e.g. Term Life, Whole Life)" />
          <v-text-field v-model.number="policyForm.coverage_amount" label="Coverage Amount" type="number" prefix="$" />
          <v-text-field v-model="policyForm.beneficiary" label="Beneficiary" />
          <v-text-field v-model="policyForm.agent_name" label="Agent Name" />
          <v-text-field v-model="policyForm.agent_phone" label="Agent Phone" />
          <v-textarea v-model="policyForm.notes" label="Notes" rows="2" />
          <v-textarea v-model="policyForm.secure_notes" label="Confidential Notes"
            rows="2" prepend-inner-icon="mdi-lock"
            hint="Will NOT appear on the printed cover letter."
            persistent-hint />
        </template>

        <template v-else-if="documentType === 'obituary_entry'">
          <v-select
            v-model="obituaryForm.type"
            label="Entry Type"
            :items="[
              { title: 'Survivor', value: 'survivor' },
              { title: 'Predeceased', value: 'predeceased' },
              { title: 'Event', value: 'event' },
            ]"
          />
          <v-text-field v-model="obituaryForm.name" label="Name" required />
          <v-text-field v-model="obituaryForm.relationship" label="Relationship" />
          <v-textarea v-model="obituaryForm.details" label="Details" rows="3" />
          <v-text-field v-model="obituaryForm.event_date" label="Event Date" type="date" />
        </template>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn @click="close">Cancel</v-btn>
        <v-btn color="primary" @click="save">Save</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { useDocumentStore } from '../stores/documents'
import { useInsurancePolicyStore } from '../stores/insurancePolicies'
import { useObituaryInfoStore } from '../stores/obituaryInfo'
import { useDocumentCategoryStore } from '../stores/documentCategories'
import { usePartyStore } from '../stores/parties'
import type { Document, InsurancePolicy } from '../types'

const props = defineProps<{
  modelValue: boolean
  editData?: Document | InsurancePolicy | null
  initialType?: 'document' | 'insurance_policy' | 'obituary_entry'
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  saved: []
}>()

const documentStore = useDocumentStore()
const policyStore = useInsurancePolicyStore()
const obituaryStore = useObituaryInfoStore()
const categoryStore = useDocumentCategoryStore()
const partyStore = usePartyStore()

const documentType = ref<'document' | 'insurance_policy' | 'obituary_entry'>('document')

const typeOptions = [
  { title: 'Document', value: 'document' },
  { title: 'Insurance Policy', value: 'insurance_policy' },
  { title: 'Obituary Entry', value: 'obituary_entry' },
]

const categoryItems = computed(() => categoryStore.categories)

const dialogTitle = computed(() => {
  if (props.editData) return 'Edit'
  const labels: Record<string, string> = {
    document: 'Add Document',
    insurance_policy: 'Add Insurance Policy',
    obituary_entry: 'Add Obituary Entry',
  }
  return labels[documentType.value] || 'Add'
})

const docForm = reactive({
  title: '',
  category_id: null as number | null,
  content: '',
  status: 'draft' as 'draft' | 'complete',
  secure_notes: '',
})

const policyForm = reactive({
  provider: '',
  policy_number: '',
  type: '',
  coverage_amount: null as number | null,
  beneficiary: '',
  agent_name: '',
  agent_phone: '',
  notes: '',
  secure_notes: '',
})

const obituaryForm = reactive({
  type: 'survivor' as 'survivor' | 'predeceased' | 'event',
  name: '',
  relationship: '',
  details: '',
  event_date: '',
})

function resetForms() {
  docForm.title = ''
  docForm.category_id = null
  docForm.content = ''
  docForm.status = 'draft'
  docForm.secure_notes = ''

  policyForm.provider = ''
  policyForm.policy_number = ''
  policyForm.type = ''
  policyForm.coverage_amount = null
  policyForm.beneficiary = ''
  policyForm.agent_name = ''
  policyForm.agent_phone = ''
  policyForm.notes = ''
  policyForm.secure_notes = ''

  obituaryForm.type = 'survivor'
  obituaryForm.name = ''
  obituaryForm.relationship = ''
  obituaryForm.details = ''
  obituaryForm.event_date = ''
}

watch(() => props.modelValue, (open) => {
  if (open) {
    documentType.value = props.initialType || 'document'
    if (props.editData) {
      // Populate form based on type
      if ('title' in props.editData) {
        documentType.value = 'document'
        const d = props.editData as Document
        docForm.title = d.title
        docForm.category_id = d.category_id
        docForm.content = d.content || ''
        docForm.status = d.status
        docForm.secure_notes = d.secure_notes || ''
      } else if ('provider' in props.editData) {
        documentType.value = 'insurance_policy'
        const p = props.editData as InsurancePolicy
        policyForm.provider = p.provider
        policyForm.policy_number = p.policy_number || ''
        policyForm.type = p.type || ''
        policyForm.coverage_amount = p.coverage_amount ?? null
        policyForm.beneficiary = p.beneficiary || ''
        policyForm.agent_name = p.agent_name || ''
        policyForm.agent_phone = p.agent_phone || ''
        policyForm.notes = p.notes || ''
        policyForm.secure_notes = p.secure_notes || ''
      }
    } else {
      resetForms()
    }
  }
})

function close() {
  emit('update:modelValue', false)
  resetForms()
}

async function save() {
  if (documentType.value === 'document') {
    if (props.editData && 'title' in props.editData) {
      await documentStore.updateDocument(props.editData.id, { ...docForm })
    } else {
      await documentStore.createDocument({ ...docForm })
    }
  } else if (documentType.value === 'insurance_policy') {
    if (props.editData && 'provider' in props.editData) {
      await policyStore.updatePolicy(props.editData.id, { ...policyForm })
    } else {
      await policyStore.createPolicy({ ...policyForm })
    }
  } else if (documentType.value === 'obituary_entry') {
    const selfParty = partyStore.selfParty
    if (selfParty) {
      await obituaryStore.createItem(selfParty.id, {
        ...obituaryForm,
        event_date: obituaryForm.event_date || null,
      })
    }
  }
  close()
  emit('saved')
}

onMounted(() => {
  categoryStore.fetchCategories()
  partyStore.fetchParties()
})
</script>
