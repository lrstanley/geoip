<template>
  <n-card class="geo-result" content-style="padding: 0;">
    <div class="flex flex-auto flex-row gap-1 mb-2 px-2 pt-2">
      <CoreTooltip :clipboard="geo.query" label="search query" class="truncate">
        <n-button size="tiny" dashed>
          <template #icon>
            <n-icon><i-mdi-magnify /></n-icon>
          </template>

          {{ geo.query }}
        </n-button>
      </CoreTooltip>

      <span class="ml-auto" />

      <CoreTooltip
        v-if="geo.asn_org"
        :label="'Autonomous System Organization: ' + geo.network"
        class="hidden! md:flex!"
      >
        <n-button
          size="tiny"
          tag="a"
          tertiary
          :bordered="false"
          type="primary"
          :href="'https://bgp.he.net/' + geo.asn + '#_whois'"
          target="_blank"
        >
          <template #icon>
            <n-icon><i-mdi-server-network /></n-icon>
          </template>

          {{ geo.asn_org }}
        </n-button>
      </CoreTooltip>

      <CoreTooltip v-if="geo.host" :clipboard="geo.host" class="truncate" label="reverse dns">
        <n-button size="tiny" tertiary type="primary" :bordered="false">
          <template #icon>
            <n-icon><i-mdi-web-check /></n-icon>
          </template>
          {{ geo.host }}
        </n-button>
      </CoreTooltip>
    </div>

    <GeoMap :value="geo" />

    <div class="flex flex-auto flex-row gap-1 mt-2 px-2 pb-1">
      <CoreTooltip
        v-if="geo.timezone"
        :clipboard="geo.timezone"
        label="timezone"
        class="hidden! md:flex!"
      >
        <n-button size="tiny" tertiary :bordered="false">
          <template #icon>
            <n-icon><i-mdi-clock /></n-icon>
          </template>

          {{ geo.timezone }}
        </n-button>
      </CoreTooltip>

      <span class="ml-auto" />

      <CoreTooltip
        v-if="geo.postal_code"
        :clipboard="geo.postal_code"
        label="postal code"
        class="truncate"
      >
        <n-button size="tiny" tertiary type="info" :bordered="false">
          <template #icon>
            <n-icon><i-mdi-sign-direction /></n-icon>
          </template>
          {{ geo.postal_code }}
        </n-button>
      </CoreTooltip>

      <CoreTooltip v-if="geo.summary" :clipboard="geo.summary" label="location" class="truncate">
        <n-button size="tiny" tertiary type="info" :bordered="false">
          <template #icon>
            <n-icon><GeoFlag :value="geo.country_abbr" :size="16" class="mt-1px" /></n-icon>
          </template>
          {{ geo.summary }}
        </n-button>
      </CoreTooltip>

      <CoreTooltip label="open in Google Maps">
        <n-button
          size="tiny"
          tag="a"
          tertiary
          :bordered="false"
          :href="googleMaps(geo.latitude, geo.longitude)"
          target="_blank"
        >
          <template #icon>
            <n-icon><i-mdi-google-maps class="text-red-500" /></n-icon>
          </template>
        </n-button>
      </CoreTooltip>
    </div>
  </n-card>
</template>

<script setup lang="ts">
import type { GeoResult } from "@/lib/api"

const props = defineProps<{
  value: GeoResult
}>()

const geo = computed(() => props.value)

function googleMaps(lat: number, lng: number) {
  return `https://google.com/maps/place/${lat},${lng}/@${lat},${lng},5z/`
}
</script>

<style scoped>
.geo-result {
  @apply shadow-md rounded-md;
}

.geo-result:not(:last-child) {
  @apply mb-4;
}
</style>
