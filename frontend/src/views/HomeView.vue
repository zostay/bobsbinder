<template>
  <div>
    <h1 class="text-h4 mb-4">Dashboard</h1>

    <GettingStartedStepper />

    <HomeChecklist />

    <HomePrintableDocuments />
  </div>
</template>

<script setup lang="ts">
import GettingStartedStepper from '../components/GettingStartedStepper.vue'
import HomeChecklist from '../components/HomeChecklist.vue'
import HomePrintableDocuments from '../components/HomePrintableDocuments.vue'
</script>
