<template>
  <v-app-bar color="primary" density="default">
    <v-app-bar-title>
      <router-link to="/" class="text-white text-decoration-none">
        Bob's Binder
      </router-link>
    </v-app-bar-title>

    <v-spacer />

    <template v-if="auth.isAuthenticated">
      <v-btn variant="text" to="/documents">Documents</v-btn>
      <v-btn variant="text" @click="auth.logout()">Logout</v-btn>
    </template>
    <template v-else>
      <v-btn variant="text" to="/login">Login</v-btn>
      <v-btn variant="outlined" to="/register">Register</v-btn>
    </template>
  </v-app-bar>
</template>

<script setup lang="ts">
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
</script>
