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

    <v-dialog v-model="showForm" max-width="600">
      <v-card>
        <v-card-title>{{ editing ? 'Edit Policy' : 'Add Policy' }}</v-card-title>
        <v-card-text>
          <v-text-field v-model="form.provider" label="Insurance Provider" required />
          <v-text-field v-model="form.policy_number" label="Policy Number" />
          <v-text-field v-model="form.type" label="Type (e.g. Term Life, Whole Life)" />
          <v-text-field v-model.number="form.coverage_amount" label="Coverage Amount" type="number" prefix="$" />
          <v-text-field v-model="form.beneficiary" label="Beneficiary" />
          <v-text-field v-model="form.agent_name" label="Agent Name" />
          <v-text-field v-model="form.agent_phone" label="Agent Phone" />
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
import { useInsurancePolicyStore } from '../stores/insurancePolicies'
import type { InsurancePolicy } from '../types'

const store = useInsurancePolicyStore()
const showForm = ref(false)
const editing = ref<number | null>(null)
const form = reactive({
  provider: '',
  policy_number: '',
  type: '',
  coverage_amount: null as number | null,
  beneficiary: '',
  agent_name: '',
  agent_phone: '',
  notes: '',
})

function resetForm() {
  form.provider = ''
  form.policy_number = ''
  form.type = ''
  form.coverage_amount = null
  form.beneficiary = ''
  form.agent_name = ''
  form.agent_phone = ''
  form.notes = ''
  editing.value = null
}

function startEdit(policy: InsurancePolicy) {
  form.provider = policy.provider
  form.policy_number = policy.policy_number || ''
  form.type = policy.type || ''
  form.coverage_amount = policy.coverage_amount ?? null
  form.beneficiary = policy.beneficiary || ''
  form.agent_name = policy.agent_name || ''
  form.agent_phone = policy.agent_phone || ''
  form.notes = policy.notes || ''
  editing.value = policy.id
  showForm.value = true
}

function closeForm() {
  showForm.value = false
  resetForm()
}

async function save() {
  if (editing.value) {
    await store.updatePolicy(editing.value, { ...form })
  } else {
    await store.createPolicy({ ...form })
  }
  closeForm()
}

onMounted(() => {
  store.fetchPolicies()
})
</script>
