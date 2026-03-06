<template>
  <v-dialog :model-value="modelValue" max-width="600" @update:model-value="$emit('update:modelValue', $event)">
    <v-card>
      <v-card-title>{{ dialogTitle }}</v-card-title>
      <v-card-text>
        <v-select
          v-model="infoType"
          label="Type"
          :items="typeOptions"
          :disabled="!!editData"
          class="mb-2"
        />

        <template v-if="isDigitalAccess">
          <v-text-field v-model="accessForm.name" label="Name" required />
          <v-text-field v-model="accessForm.username" label="Username" />
          <v-textarea v-model="accessForm.instructions" label="Access Instructions" rows="3" />
        </template>

        <template v-else>
          <v-text-field v-model="serviceForm.name" label="Name" required />
          <v-text-field v-model="serviceForm.provider" label="Provider" />
          <v-text-field v-model="serviceForm.account_number" label="Account Number" />
          <v-text-field v-model="serviceForm.contact_name" label="Contact Name" />
          <v-text-field v-model="serviceForm.contact_phone" label="Contact Phone" />
          <v-text-field v-model="serviceForm.contact_email" label="Contact Email" />
          <v-textarea v-model="serviceForm.notes" label="Notes" rows="2" />
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
import { ref, reactive, computed, watch } from 'vue'
import { useDigitalAccessStore } from '../stores/digitalAccess'
import { useServiceAccountStore } from '../stores/serviceAccounts'
import type { DigitalAccess, ServiceAccount } from '../types'

type InfoType = 'computer' | 'phone' | 'password_manager' | 'financial_tool' | 'backup_service' | 'tax_preparer'

const props = defineProps<{
  modelValue: boolean
  editData?: DigitalAccess | ServiceAccount | null
  initialType?: InfoType
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  saved: []
}>()

const accessStore = useDigitalAccessStore()
const serviceStore = useServiceAccountStore()

const infoType = ref<InfoType>('computer')

const digitalAccessTypes = ['computer', 'phone', 'password_manager']

const typeOptions = [
  { title: 'Computer', value: 'computer' },
  { title: 'Phone', value: 'phone' },
  { title: 'Password Manager', value: 'password_manager' },
  { title: 'Financial Tool', value: 'financial_tool' },
  { title: 'Backup Service', value: 'backup_service' },
  { title: 'Tax Preparation Software', value: 'tax_preparer' },
]

const isDigitalAccess = computed(() => digitalAccessTypes.includes(infoType.value))

const dialogTitle = computed(() => {
  if (props.editData) return 'Edit'
  return 'Add Digital Info'
})

const accessForm = reactive({
  name: '',
  username: '',
  instructions: '',
})

const serviceForm = reactive({
  name: '',
  provider: '',
  account_number: '',
  contact_name: '',
  contact_phone: '',
  contact_email: '',
  notes: '',
})

function resetForms() {
  accessForm.name = ''
  accessForm.username = ''
  accessForm.instructions = ''

  serviceForm.name = ''
  serviceForm.provider = ''
  serviceForm.account_number = ''
  serviceForm.contact_name = ''
  serviceForm.contact_phone = ''
  serviceForm.contact_email = ''
  serviceForm.notes = ''
}

watch(() => props.modelValue, (open) => {
  if (open) {
    infoType.value = props.initialType || 'computer'
    if (props.editData) {
      if ('username' in props.editData) {
        const da = props.editData as DigitalAccess
        infoType.value = da.type
        accessForm.name = da.name
        accessForm.username = da.username || ''
        accessForm.instructions = da.instructions || ''
      } else {
        const sa = props.editData as ServiceAccount
        infoType.value = sa.type as InfoType
        serviceForm.name = sa.name
        serviceForm.provider = sa.provider || ''
        serviceForm.account_number = sa.account_number || ''
        serviceForm.contact_name = sa.contact_name || ''
        serviceForm.contact_phone = sa.contact_phone || ''
        serviceForm.contact_email = sa.contact_email || ''
        serviceForm.notes = sa.notes || ''
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
  if (isDigitalAccess.value) {
    if (props.editData && 'username' in props.editData) {
      await accessStore.updateItem(props.editData.id, {
        ...accessForm,
        type: infoType.value as DigitalAccess['type'],
      })
    } else {
      await accessStore.createItem({
        ...accessForm,
        type: infoType.value as DigitalAccess['type'],
      })
    }
  } else {
    if (props.editData && !('username' in props.editData)) {
      await serviceStore.updateAccount(props.editData.id, {
        ...serviceForm,
        type: infoType.value as ServiceAccount['type'],
      })
    } else {
      await serviceStore.createAccount({
        ...serviceForm,
        type: infoType.value as ServiceAccount['type'],
      })
    }
  }
  close()
  emit('saved')
}
</script>
