<template>
  <router-view v-slot="{ Component, route }">
    <AnimateFade>
      <component :is="Component" :key="route" />
    </AnimateFade>
  </router-view>
</template>

<script setup lang="ts"></script>

<style scoped></style>
