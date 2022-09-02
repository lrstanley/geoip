<template>
  <span ref="container">
    <AnimateFade>
      <n-avatar
        v-if="loaded"
        round
        size="small"
        :src="flag"
        :fallback-src="FLAG_URI + '/xx.svg'"
        v-bind="$attrs"
      />
    </AnimateFade>
  </span>
</template>

<script setup lang="ts">
const props = defineProps<{
  value: string
  immediate?: boolean
}>()

const FLAG_URI = "https://hatscripts.github.io/circle-flags/flags"
const flag = computed(() => {
  let code = props.value.toLowerCase()
  if (!code || code == "other") {
    code = "xx"
  }
  return `${FLAG_URI}/${code}.svg`
})

const loaded = ref(false)
const container = ref<HTMLElement>(null)

// Lazy load the flag if requested.
const { stop } = useIntersectionObserver(container, ([{ isIntersecting }]) => {
  if (isIntersecting) {
    stop()
    loaded.value = true
  }
})

onMounted(() => {
  if (props.immediate) {
    stop()
    loaded.value = true
  }
})
</script>
