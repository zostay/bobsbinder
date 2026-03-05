<template>
  <v-row justify="center" class="mt-8">
    <v-col cols="12" sm="8" md="4">
      <v-card>
        <v-card-title class="text-h5 text-center pa-4">Login</v-card-title>
        <v-card-text>
          <v-form @submit.prevent="handleLogin">
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
              Login
            </v-btn>
          </v-form>
        </v-card-text>
        <v-card-actions class="justify-center">
          <span>Don't have an account?</span>
          <v-btn variant="text" color="primary" to="/register">Register</v-btn>
        </v-card-actions>
      </v-card>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

async function handleLogin() {
  loading.value = true
  error.value = ''
  try {
    await auth.login(email.value, password.value)
  } catch (e: any) {
    error.value = e.response?.data?.error || 'Login failed'
  } finally {
    loading.value = false
  }
}
</script>
