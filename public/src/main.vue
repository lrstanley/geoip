<template>
  <n-config-provider v-bind="configProvider">
    <router-view />
  </n-config-provider>
</template>

<script setup lang="ts">
import { configProvider } from "@/lib/core/status"
</script>
