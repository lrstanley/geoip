<route>
meta:
  title: Bulk Lookup
</route>

<template>
  <div class="transition-all ease-in-out duration-500">
    <div class="p-4">
      <AnimateFade>
        <div class="absolute top-6 right-6 z-1000 flex gap-1">
          <n-tag v-if="searchCount > 0"> {{ searchCount }} addresses </n-tag>
          <n-tag
            v-if="input.length > 0"
            class="cursor-pointer"
            @click=";(input = ''), (searchCount = 0)"
          >
            reset
          </n-tag>
        </div>
      </AnimateFade>
      <n-input
        v-model:value="input"
        v-focus
        type="textarea"
        size="large"
        rows="10"
        placeholder="Bulk search IPs (8.8.8.8) or hosts (google.com)"
        :disabled="loading"
      />
    </div>
    <div class="flex flex-row px-4 pb-4 gap-2">
      <RateLimitCounter allow-small />

      <span class="ml-auto" />

      <CoreTooltip v-if="input.length < 1" label="Use sample data">
        <n-button size="small" type="info" class="hidden md:flex" dashed @click="useSampleData()">
          <n-icon><i-mdi-database-plus-outline /></n-icon>
        </n-button>
      </CoreTooltip>
      <n-button
        v-show="(results.length > 0 || errors.length > 0) && !loading"
        size="small"
        type="warning"
        @click="clearResults()"
      >
        <n-icon><i-mdi-broom /></n-icon>
        clear
      </n-button>
      <CoreTooltip v-if="results.length > 0 && !loading" label="Download results as JSON">
        <n-button
          size="small"
          type="info"
          tag="a"
          :href="jsonUrl"
          download="geoip-results.json"
          class="hidden md:flex"
        >
          <n-icon><i-mdi-floppy /></n-icon>
        </n-button>
      </CoreTooltip>
      <CoreTooltip v-if="results.length > 0 && !loading" label="Open JSON in new tab">
        <n-button
          size="small"
          type="info"
          tag="a"
          :href="jsonUrl"
          target="_blank"
          class="hidden md:flex"
        >
          <n-icon><i-mdi-open-in-new /></n-icon>
        </n-button>
      </CoreTooltip>
      <n-button
        size="small"
        type="primary"
        :disabled="loading"
        :loading="loading"
        @click="!loading && search()"
      >
        <n-icon><i-mdi-magnify /></n-icon>
        search
      </n-button>
    </div>

    <n-divider v-show="loading || percent > 0" class="m-0!" />
    <AnimateFade v-show="loading || percent > 0">
      <div class="p-4">
        <n-progress
          type="line"
          :status="loading ? 'info' : 'success'"
          :percentage="percent"
          :processing="loading"
        >
          {{ completed }}/{{ total }}
        </n-progress>
      </div>
    </AnimateFade>

    <n-divider v-show="errors.length > 0 && !loading" class="m-0!" />
    <GeoMultiError v-show="errors.length > 0 && !loading" :value="errors" class="p-4" />

    <n-divider v-show="results.length > 0" class="m-0!" />
    <GeoAggregate
      v-show="results.length > 0"
      class="p-4"
      :value="results"
      :truncate="{ percent: 5, label: 'Other' }"
      field="country_abbr"
      label="country"
      flag
    />

    <n-divider v-show="results.length > 0" class="m-0!" />
    <GeoAggregate
      v-show="results.length > 0"
      class="p-4"
      :value="results"
      field="continent_abbr"
      label="continent"
      :truncate="{ percent: 5, label: 'Other' }"
      status="success"
    />
  </div>
</template>

<script setup async lang="ts">
import sampleData from "@/lib/data/sample-bulk.txt?raw"
import { lookupMany } from "@/lib/api"
import { matchAddresses } from "@/lib/util/match"
import { createJSONObjectURL } from "@/lib/util/object-url"
import vFocus from "@/lib/directives/focus"
import type { APIResponse } from "@/lib/api"

const input = ref("")
const results = ref<APIResponse[]>([]) // Aggregate of all results (not just the current input).
const errors = ref<APIResponse[]>([]) // Aggregate of all failed results.
const loading = ref(false)
const completed = ref(0) // Completed requests.
const total = ref(0) // Total addresses to lookup.
const percent = computed(() => Math.round((completed.value / total.value) * 100))

const jsonUrl = createJSONObjectURL(
  computed(() => {
    return { data: results.value.map(({ data }) => data) }
  })
)

async function search() {
  const addresses = matchAddresses(input.value)

  if (!input.value || addresses.length < 1) return

  loading.value = true
  total.value = addresses.length
  completed.value = 0
  errors.value = []

  await lookupMany(addresses, false, (result) => {
    nextTick(() => {
      completed.value++

      if (result.error) {
        errors.value.push(result)
      } else {
        results.value.push(result)
      }
    })
  })

  loading.value = false
}

function clearResults() {
  completed.value = 0
  total.value = 0
  results.value = []
  errors.value = []
}

function useSampleData() {
  input.value = sampleData
}

const searchCount = ref(0)

watchDebounced(input, () => (searchCount.value = matchAddresses(input.value).length), {
  debounce: 500,
  maxWait: 1000,
})
</script>
