<template>
  <div ref="container" style="height: 200px" />
</template>

<script setup lang="ts">
import "leaflet/dist/leaflet.css"
import markerIcon from "leaflet/dist/images/marker-icon-2x.png"
import markerIconShadow from "leaflet/dist/images/marker-shadow.png"
import { Icon, Map, TileLayer, Marker } from "leaflet"
import type { GeoIPData } from "@/lib/api"

const props = defineProps<{
  value: GeoIPData
  zoom?: number
  scrollWheel?: boolean
}>()

// Fix leaflet not importing marker icons correctly.
if (import.meta.env.PROD) {
  Icon.Default.mergeOptions({
    iconRetinaUrl: markerIcon,
    iconUrl: markerIcon,
    shadowUrl: markerIconShadow,
  })
}

const TILE_URL = "https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
const ATTRIBUTION = '&copy; <a href="https://www.openstreetmap.org/copyright">OSM</a>'

const container = ref<HTMLElement>()
const map = ref<Map | null>()

function setupMap() {
  map.value = new Map(container.value, {
    preferCanvas: true,
    scrollWheelZoom: props.scrollWheel ?? true,
  }).setView([props.value.latitude, props.value.longitude], props.zoom ?? 5)

  new TileLayer(TILE_URL, { attribution: ATTRIBUTION, maxZoom: 18 }).addTo(map.value)
  new Marker([props.value.latitude, props.value.longitude]).addTo(map.value)

  setTimeout(() => {
    map.value.invalidateSize()
  }, 400)
}

// Lazy load the map.
const { stop } = useIntersectionObserver(container, ([{ isIntersecting }]) => {
  if (isIntersecting) {
    stop()
    setupMap()
  }
})
</script>

<style scoped>
@media (prefers-color-scheme: dark), (prefers-color-scheme: no-preference) {
  div {
    filter: invert(100%) hue-rotate(180deg) brightness(95%) contrast(90%);
  }
}
</style>
