<template>
  <div>
    <div class="d-flex align-center mb-4">
      <h1 class="text-h4">Insurance Policies</h1>
      <v-spacer />
      <v-btn color="primary" prepend-icon="mdi-plus" @click="showForm = true">
        Add Policy
      </v-btn>
    </div>

    <v-progress-linear v-if="store.loading" indeterminate color="primary" />

    <v-card v-if="store.policies.length === 0 && !store.loading">
      <v-card-text class="text-center pa-8">
        <v-icon size="64" color="grey">mdi-shield-check-outline</v-icon>
        <p class="text-h6 mt-4">No insurance policies yet</p>
        <p class="text-body-2 text-grey">Add life insurance policies so your survivors know how to file claims.</p>
      </v-card-text>
    </v-card>

    <v-card v-for="policy in store.policies" :key="policy.id" class="mb-2">
      <v-card-title>{{ policy.provider }}</v-card-title>
      <v-card-subtitle v-if="policy.type">{{ policy.type }}</v-card-subtitle>
      <v-card-text>
        <div v-if="policy.policy_number">Policy #: {{ policy.policy_number }}</div>
        <div v-if="policy.coverage_amount">Coverage: ${{ policy.coverage_amount?.toLocaleString() }}</div>
        <div v-if="policy.beneficiary">Beneficiary: {{ policy.beneficiary }}</div>
        <div v-if="policy.agent_name">Agent: {{ policy.agent_name }} <span v-if="policy.agent_phone">- {{ policy.agent_phone }}</span></div>
        <div v-if="policy.notes" class="mt-2 text-grey">{{ policy.notes }}</div>
      </v-card-text>
      <v-card-actions>
        <v-btn size="small" variant="text" color="primary" @click="startEdit(policy)">Edit</v-btn>
        <v-btn size="small" variant="text" color="error" @click="store.deletePolicy(policy.id)">Delete</v-btn>
      </v-card-actions>
    </v-card>

    <DocumentFormDialog
      v-model="showForm"
      :edit-data="editingPolicy"
      initial-type="insurance_policy"
      @saved="onSaved"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useInsurancePolicyStore } from '../stores/insurancePolicies'
import DocumentFormDialog from '../components/DocumentFormDialog.vue'
import type { InsurancePolicy } from '../types'

const store = useInsurancePolicyStore()
const showForm = ref(false)
const editingPolicy = ref<InsurancePolicy | null>(null)

function startEdit(policy: InsurancePolicy) {
  editingPolicy.value = policy
  showForm.value = true
}

function onSaved() {
  editingPolicy.value = null
  store.fetchPolicies()
}

onMounted(() => {
  store.fetchPolicies()
})
</script>
