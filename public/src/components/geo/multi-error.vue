<template>
  <AnimateFade>
    <div v-bind="$attrs">
      <n-alert v-if="!showErrors" type="error" class="alert-small" :show-icon="false">
        <div class="flex flex-auto">
          <div class="pl-2">
            {{ props.value.length }} {{ props.value.length > 1 ? "errors" : "error" }} occurred
          </div>
          <div class="ml-auto">
            <n-button size="tiny" @click="showErrors = true">show</n-button>
          </div>
        </div>
      </n-alert>
      <AnimateListGroup v-else>
        <n-alert
          v-for="result in props.value"
          :key="result.query"
          :show-icon="false"
          type="error"
          class="alert-small mb-2 last:mb-0"
        >
          <n-tag size="small">Q: {{ result.query }}</n-tag>
          error: {{ result.error }}
        </n-alert>
      </AnimateListGroup>
    </div>
  </AnimateFade>
</template>

<script setup lang="ts">
import type { BulkError } from "@/lib/api"
const props = defineProps<{
  value: BulkError[]
}>()

watch(props, () => {
  showErrors.value = false
})

const showErrors = ref(false) // Whether the errors are currently visible.
</script>
