<template>
  <AppNavbar />
  <v-main>
    <v-container>
      <slot />
    </v-container>
  </v-main>
</template>

<script setup lang="ts">
import AppNavbar from '../components/AppNavbar.vue'
</script>
