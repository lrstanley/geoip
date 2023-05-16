<template>
  <CoreTooltip v-show="percent" :label="label">
    <n-tag type="info" v-bind="$attrs">
      <template #icon>
        <n-icon><i-mdi-timer-sync-outline /></n-icon>
      </template>
      {{ percent }}%
      <span :class="{ 'hidden md:inline-flex': props.allowSmall }">calls left</span>
    </n-tag>
  </CoreTooltip>
</template>

<script setup lang="ts">
const props = defineProps<{
  allowSmall?: boolean
}>()

const state = useState()

const remaining = computed(() => state.clientState.ratelimit_remaining ?? 0)
const limit = computed(() => state.clientState.ratelimit_limit ?? 0)
const label = computed(() => {
  return `${remaining.value.toLocaleString()} left of ${limit.value.toLocaleString()} limit`
})

const percent = computed(() => {
  return Math.floor((remaining.value / limit.value) * 100)
})
</script>
