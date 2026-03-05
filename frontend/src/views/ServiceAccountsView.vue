<template>
  <div>
    <div class="d-flex align-center mb-4">
      <h1 class="text-h4">Service Accounts</h1>
      <v-spacer />
      <v-btn color="primary" prepend-icon="mdi-plus" @click="showForm = true">
        Add Account
      </v-btn>
    </div>

    <v-progress-linear v-if="store.loading" indeterminate color="primary" />

    <v-card v-if="store.accounts.length === 0 && !store.loading">
      <v-card-text class="text-center pa-8">
        <v-icon size="64" color="grey">mdi-briefcase-outline</v-icon>
        <p class="text-h6 mt-4">No service accounts yet</p>
        <p class="text-body-2 text-grey">Add tax preparers, financial tools, and backup services.</p>
      </v-card-text>
    </v-card>

    <v-card v-for="account in store.accounts" :key="account.id" class="mb-2">
      <v-card-title>{{ account.name }}</v-card-title>
      <v-card-subtitle>
        <v-chip size="small" color="info">{{ typeLabels[account.type] || account.type }}</v-chip>
      </v-card-subtitle>
      <v-card-text>
        <div v-if="account.provider">Provider: {{ account.provider }}</div>
        <div v-if="account.account_number">Account #: {{ account.account_number }}</div>
        <div v-if="account.contact_name">
          Contact: {{ account.contact_name }}
          <span v-if="account.contact_phone"> - {{ account.contact_phone }}</span>
          <span v-if="account.contact_email"> - {{ account.contact_email }}</span>
        </div>
        <div v-if="account.notes" class="mt-2 text-grey">{{ account.notes }}</div>
      </v-card-text>
      <v-card-actions>
        <v-btn size="small" variant="text" color="primary" @click="startEdit(account)">Edit</v-btn>
        <v-btn size="small" variant="text" color="error" @click="store.deleteAccount(account.id)">Delete</v-btn>
      </v-card-actions>
    </v-card>

    <v-dialog v-model="showForm" max-width="600">
      <v-card>
        <v-card-title>{{ editing ? 'Edit Account' : 'Add Account' }}</v-card-title>
        <v-card-text>
          <v-text-field v-model="form.name" label="Name" required />
          <v-select
            v-model="form.type"
            label="Type"
            :items="[
              { title: 'Financial Tool', value: 'financial_tool' },
              { title: 'Backup Service', value: 'backup_service' },
              { title: 'Tax Preparer', value: 'tax_preparer' },
            ]"
            required
          />
          <v-text-field v-model="form.provider" label="Provider" />
          <v-text-field v-model="form.account_number" label="Account Number" />
          <v-text-field v-model="form.contact_name" label="Contact Name" />
          <v-text-field v-model="form.contact_phone" label="Contact Phone" />
          <v-text-field v-model="form.contact_email" label="Contact Email" />
          <v-textarea v-model="form.notes" label="Notes" rows="2" />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn @click="closeForm">Cancel</v-btn>
          <v-btn color="primary" @click="save">Save</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { useServiceAccountStore } from '../stores/serviceAccounts'
import type { ServiceAccount } from '../types'

const store = useServiceAccountStore()
const showForm = ref(false)
const editing = ref<number | null>(null)

const typeLabels: Record<string, string> = {
  financial_tool: 'Financial Tool',
  backup_service: 'Backup Service',
  tax_preparer: 'Tax Preparer',
}

const form = reactive({
  name: '',
  type: 'financial_tool' as 'financial_tool' | 'backup_service' | 'tax_preparer',
  provider: '',
  account_number: '',
  contact_name: '',
  contact_phone: '',
  contact_email: '',
  notes: '',
})

function resetForm() {
  form.name = ''
  form.type = 'financial_tool'
  form.provider = ''
  form.account_number = ''
  form.contact_name = ''
  form.contact_phone = ''
  form.contact_email = ''
  form.notes = ''
  editing.value = null
}

function startEdit(account: ServiceAccount) {
  form.name = account.name
  form.type = account.type
  form.provider = account.provider || ''
  form.account_number = account.account_number || ''
  form.contact_name = account.contact_name || ''
  form.contact_phone = account.contact_phone || ''
  form.contact_email = account.contact_email || ''
  form.notes = account.notes || ''
  editing.value = account.id
  showForm.value = true
}

function closeForm() {
  showForm.value = false
  resetForm()
}

async function save() {
  if (editing.value) {
    await store.updateAccount(editing.value, { ...form })
  } else {
    await store.createAccount({ ...form })
  }
  closeForm()
}

onMounted(() => {
  store.fetchAccounts()
})
</script>
