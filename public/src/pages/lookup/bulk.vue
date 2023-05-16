<route>
meta:
  title: Bulk Lookup
</route>

<template>
  <LayoutDefault>
    <div class="transition-all ease-in-out duration-500">
      <div class="p-4">
        <AnimateFade>
          <div class="absolute top-6 right-6 z-1000 flex gap-1">
            <n-tag v-if="searchCount > 0">{{ searchCount }} addresses</n-tag>
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
          id="bulk-clear"
          size="small"
          type="warning"
          @click="clearResults()"
        >
          <n-icon><i-mdi-broom /></n-icon>
          clear
        </n-button>
        <CoreTooltip v-if="results.length > 0 && !loading" label="Download results as JSON">
          <n-button
            id="bulk-results-download"
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
            id="bulk-results-open"
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
        <div id="bulk-progress" class="p-4">
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
        id="aggregate-country"
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
        id="aggregate-continent"
        class="p-4"
        :value="results"
        field="continent_abbr"
        label="continent"
        :truncate="{ percent: 5, label: 'Other' }"
        status="success"
      />

      <n-divider v-show="results.length > 0" class="m-0!" />
      <GeoAggregate
        v-show="results.length > 0"
        id="aggregate-asn"
        class="p-4"
        :value="results"
        field="asn_org"
        label="asn_org"
        :truncate="{ percent: 5, label: 'Other' }"
        status="success"
      />
    </div>
  </LayoutDefault>
</template>

<script setup async lang="ts">
import { api } from "@/lib/api"
import { matchAddresses } from "@/lib/util/match"
import { createJSONObjectURL } from "@/lib/util/object-url"
import type { GeoResult, BulkError } from "@/lib/api"

const input = ref("")
const results = ref<GeoResult[]>([]) // Aggregate of all results (not just the current input).
const errors = ref<BulkError[]>([]) // Aggregate of all failed results.
const loading = ref(false)
const completed = ref(0) // Completed requests.
const total = ref(0) // Total addresses to lookup.
const percent = computed(() => Math.round((completed.value / total.value) * 100))

const jsonUrl = createJSONObjectURL(
  computed(() => {
    return { data: results.value }
  })
)

async function search() {
  const addresses = matchAddresses(input.value)

  if (!input.value || addresses.length < 1) return

  loading.value = true
  total.value = addresses.length
  completed.value = 0
  errors.value = []

  // Chunk the addresses into groups of 25, the maximum allowed by the API.
  const chunks = []
  while (addresses.length > 0) {
    chunks.push(addresses.splice(0, 25))
  }

  for (const chunk of chunks) {
    try {
      const resp = await api.lookup.getManyAddresses({
        requestBody: {
          addresses: chunk,
          options: {
            disable_host_lookup: true,
          },
        },
      })

      results.value.push(...resp.results)
      errors.value.push(...resp.errors)
    } catch (error) {
      errors.value.push({
        error: error,
        query: "-",
      })
    } finally {
      completed.value += chunk.length
    }
  }

  loading.value = false
}

function clearResults() {
  completed.value = 0
  total.value = 0
  results.value = []
  errors.value = []
}

const searchCount = ref(0)

watchDebounced(input, () => (searchCount.value = matchAddresses(input.value).length), {
  debounce: 500,
  maxWait: 1000,
})

function useSampleData() {
  input.value = sampleData.join("\n")
}

const sampleData = [
  "google.com",
  "facebook.com",
  "ebay.co.uk",
  "github.com",
  "amazon.ca",
  "101.255.65.138",
  "101.32.11.149",
  "1.0.233.26",
  "103.102.42.42",
  "103.105.130.83",
  "103.109.74.14",
  "103.129.221.188",
  "142.93.84.194",
  "143.110.251.175",
  "143.198.123.124",
  "143.198.154.97",
  "143.198.171.44",
  "143.198.200.155",
  "143.198.211.87",
  "143.198.231.219",
  "143.244.130.229",
  "143.244.145.146",
  "143.244.154.61",
  "143.244.161.152",
  "143.244.162.174",
  "14.3.9.46",
  "144.34.250.161",
  "144.48.227.75",
  "144.91.113.135",
  "144.91.117.81",
  "144.91.68.182",
  "146.185.137.240",
  "146.190.227.169",
  "146.190.60.149",
  "14.63.212.60",
  "147.182.135.41",
  "147.182.167.232",
  "104.191.173.169",
  "104.206.128.22",
  "104.208.118.5",
  "104.211.77.31",
  "104.218.164.12",
  "104.225.159.240",
  "104.236.151.120",
  "104.236.182.223",
  "104.244.108.67",
  "104.248.117.154",
  "119.65.149.106",
  "120.89.46.71",
  "121.130.111.133",
  "121.150.205.245",
  "121.155.168.49",
  "121.159.171.57",
  "121.169.150.161",
  "121.179.208.82",
  "121.184.138.195",
  "121.185.123.67",
  "121.187.251.210",
  "121.200.55.93",
  "121.227.89.109",
  "121.30.226.73",
  "121.46.24.111",
  "121.46.30.135",
  "12.171.207.202",
  "122.116.47.83",
  "122.117.248.97",
  "122.117.25.149",
  "122.117.28.77",
  "122.117.7.175",
  "122.117.88.125",
  "122.146.196.217",
  "122.168.125.16",
  "122.170.13.184",
  "122.170.9.211",
  "122.176.19.65",
  "122.180.243.249",
  "122.254.114.2",
  "1.226.228.82",
  "123.125.194.157",
  "123.126.106.88",
  "123.143.203.67",
  "123.194.184.204",
  "123.205.58.175",
  "123.24.206.219",
  "123.41.0.35",
  "1.234.58.214",
  "1.234.58.225",
  "1.235.128.206",
  "1.235.192.218",
  "124.106.69.18",
  "124.30.44.214",
  "125.129.82.220",
  "125.132.41.164",
  "125.160.99.98",
  "125.212.203.113",
  "107.175.222.27",
  "107.179.222.3",
  "107.182.129.240",
  "109.234.39.250",
  "109.62.200.55",
]
</script>
