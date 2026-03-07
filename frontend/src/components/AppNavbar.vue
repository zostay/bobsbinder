<template>
  <v-app-bar color="primary" density="default">
    <v-app-bar-title>
      <router-link to="/" class="text-white text-decoration-none">
        Bob's Binder
      </router-link>
    </v-app-bar-title>

    <v-spacer />

    <template v-if="auth.isAuthenticated">
      <v-btn variant="text" to="/">Dashboard</v-btn>
      <v-btn variant="text" to="/cover-letter">Cover Letter</v-btn>
      <v-btn variant="text" to="/documents">Documents</v-btn>
      <v-menu>
        <template v-slot:activator="{ props }">
          <v-btn variant="text" v-bind="props" append-icon="mdi-chevron-down">Data</v-btn>
        </template>
        <v-list>
          <v-list-item to="/contacts" title="Contacts" />
          <v-list-item to="/locations" title="Locations" />
          <v-list-item to="/digital-access" title="Digital Access" />
          <v-list-item to="/insurance-policies" title="Insurance Policies" />
          <v-list-item to="/service-accounts" title="Service Accounts" />
        </v-list>
      </v-menu>
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
