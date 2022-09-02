<template>
  <div ref="container" style="height: 200px" />
</template>

<script setup lang="ts">
import "leaflet/dist/leaflet.css"
import markerIcon from "leaflet/dist/images/marker-icon-2x.png"
import markerIconShadow from "leaflet/dist/images/marker-shadow.png"
import { Icon, Map, TileLayer, Marker, Circle, Control } from "leaflet"
import type { GeoResult } from "@/lib/api"

const props = defineProps<{
  value: GeoResult
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
  const tiles = new TileLayer(TILE_URL, { attribution: ATTRIBUTION, maxZoom: 18 })
  const marker = new Marker([props.value.latitude, props.value.longitude])
  const scale = new Control.Scale()

  let zoom = 5

  let circle: Circle | null = null
  if (props.value.accuracy_radius_km > 0) {
    circle = new Circle([props.value.latitude, props.value.longitude], {
      radius: props.value.accuracy_radius_km * 1000,
      weight: 1.5,
      color: "red",
      fillColor: "red",
      fillOpacity: 0.05,
    })

    if (props.value.accuracy_radius_km <= 25) {
      zoom = 9
    } else if (props.value.accuracy_radius_km <= 50) {
      zoom = 8
    } else if (props.value.accuracy_radius_km <= 200) {
      zoom = 6
    } else if (props.value.accuracy_radius_km <= 500) {
      zoom = 5
    } else if (props.value.accuracy_radius_km <= 1000) {
      zoom = 3
    } else {
      zoom = 2
    }
  }

  map.value = new Map(container.value, {
    preferCanvas: true,
    scrollWheelZoom: props.scrollWheel ?? true,
  }).setView([props.value.latitude, props.value.longitude], zoom)

  tiles.addTo(map.value)
  scale.addTo(map.value)

  if (circle) {
    circle.addTo(map.value)
  } else {
    marker.addTo(map.value)
  }

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
