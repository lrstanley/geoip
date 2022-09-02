<template>
  <div class="inline-flex" v-bind="$attrs" @click="clipboard && copyToClipboard()">
    <n-tooltip>
      <template #trigger>
        <slot />
      </template>
      <span>{{ props.label }}</span>
    </n-tooltip>
  </div>
</template>

<script setup lang="ts">
import { notification } from "@/lib/core/status"

const props = defineProps<{
  label: string
  clipboard?: string
}>()

const { copy } = useClipboard()

function copyToClipboard() {
  copy(props.clipboard)
  notification.success({
    duration: 2000,
    content: `copied "${props.clipboard}"`,
  })
}
</script>
