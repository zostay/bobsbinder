<template>
  <v-row justify="center" class="mt-8">
    <v-col cols="12" sm="8" md="4">
      <v-card>
        <v-card-title class="text-h5 text-center pa-4">Create Account</v-card-title>
        <v-card-text>
          <v-form @submit.prevent="handleRegister">
            <v-text-field
              v-model="name"
              label="Full Name"
              required
            />
            <v-text-field
              v-model="email"
              label="Email"
              type="email"
              required
              :error-messages="error ? [error] : []"
            />
            <v-text-field
              v-model="password"
              label="Password"
              type="password"
              required
            />
            <v-btn
              type="submit"
              color="primary"
              block
              size="large"
              :loading="loading"
              class="mt-4"
            >
              Register
            </v-btn>
          </v-form>
        </v-card-text>
        <v-card-actions class="justify-center">
          <span>Already have an account?</span>
          <v-btn variant="text" color="primary" to="/login">Login</v-btn>
        </v-card-actions>
      </v-card>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const name = ref('')
const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

async function handleRegister() {
  loading.value = true
  error.value = ''
  try {
    await auth.register(email.value, password.value, name.value)
  } catch (e: any) {
    error.value = e.response?.data?.error || 'Registration failed'
  } finally {
    loading.value = false
  }
}
</script>
