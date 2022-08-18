<template>
  <n-card class="geo-result" content-style="padding: 0;">
    <div class="flex flex-auto flex-row gap-1 mb-2 px-2 pt-2">
      <CoreTooltip :clipboard="props.value.query" label="search query" class="truncate">
        <n-button size="tiny" dashed>
          <template #icon>
            <n-icon><i-mdi-magnify /></n-icon>
          </template>

          {{ props.value.query }}
        </n-button>
      </CoreTooltip>

      <CoreTooltip v-if="data.host" :clipboard="data.host" class="ml-auto truncate" label="reverse dns">
        <n-button size="tiny" tertiary type="primary" :bordered="false">
          <template #icon>
            <n-icon><i-mdi-web-check /></n-icon>
          </template>
          {{ data.host }}
        </n-button>
      </CoreTooltip>
    </div>

    <GeoMap :value="data" />

    <div class="flex flex-auto flex-row gap-1 mt-2 px-2 pb-1">
      <CoreTooltip
        v-if="data.timezone"
        :clipboard="data.timezone"
        label="timezone"
        class="hidden! md:flex!"
      >
        <n-button size="tiny" tertiary :bordered="false">
          <template #icon>
            <n-icon><i-mdi-clock /></n-icon>
          </template>

          {{ data.timezone }}
        </n-button>
      </CoreTooltip>

      <span class="ml-auto" />

      <CoreTooltip
        v-if="data.postal_code"
        :clipboard="data.postal_code"
        label="postal code"
        class="truncate"
      >
        <n-button size="tiny" tertiary type="info" :bordered="false">
          <template #icon>
            <n-icon><i-mdi-sign-direction /></n-icon>
          </template>
          {{ data.postal_code }}
        </n-button>
      </CoreTooltip>

      <CoreTooltip v-if="data.summary" :clipboard="data.summary" label="location" class="truncate">
        <n-button size="tiny" tertiary type="info" :bordered="false">
          <template #icon>
            <n-icon><GeoFlag :value="props.value.data.country_abbr" :size="16" class="mt-1px" /></n-icon>
          </template>
          {{ data.summary }}
        </n-button>
      </CoreTooltip>
    </div>
  </n-card>
</template>

<script setup lang="ts">
import type { APIResponse } from "@/lib/api"

const props = defineProps<{
  value: APIResponse
}>()

const data = computed(() => props.value.data)
</script>

<style scoped>
.geo-result {
  @apply shadow-md rounded-md;
}

.geo-result:not(:last-child) {
  @apply mb-4;
}
</style>
