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

    <DigitalInfoFormDialog
      v-model="showForm"
      :edit-data="editingAccount"
      initial-type="financial_tool"
      @saved="onSaved"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useServiceAccountStore } from '../stores/serviceAccounts'
import DigitalInfoFormDialog from '../components/DigitalInfoFormDialog.vue'
import type { ServiceAccount } from '../types'

const store = useServiceAccountStore()
const showForm = ref(false)
const editingAccount = ref<ServiceAccount | null>(null)

const typeLabels: Record<string, string> = {
  financial_tool: 'Financial Tool',
  backup_service: 'Backup Service',
  tax_preparer: 'Tax Preparer',
}

function startEdit(account: ServiceAccount) {
  editingAccount.value = account
  showForm.value = true
}

function onSaved() {
  editingAccount.value = null
  store.fetchAccounts()
}

onMounted(() => {
  store.fetchAccounts()
})
</script>
