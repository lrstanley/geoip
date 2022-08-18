<template>
  <br v-if="props.nopadding === false" />
  <div :id="props.id" class="flex flex-auto flex-row">
    <n-icon v-if="props.icon" size="50" class="mr-4">
      <slot name="icon" />
    </n-icon>
    <div class="flex flex-auto flex-col">
      <component
        :is="props.el ?? 'h2'"
        class="m-0 text-gradient bg-gradient-to-br from-pink-500 via-red-500 to-yellow-500"
      >
        {{ props.title }}
        <a v-if="props.id" :href="'#' + id" class="text-emerald-800" @click="(e) => scrollIntoView(e)">
          <i-mdi-link />
        </a>
      </component>
      <span>{{ props.subtitle }}</span>
    </div>
  </div>
  <n-divider class="mt-4!" />
</template>

<script setup lang="ts">
const props = defineProps<{
  title: string
  subtitle: string
  icon?: boolean
  id?: string
  el?: string
  nopadding?: boolean
}>()

function scrollIntoView(e: MouseEvent) {
  e.preventDefault()
  history.pushState({}, null, location.pathname + "#" + encodeURIComponent(props.id))

  const el = document.getElementById(props.id)

  if (el) {
    el.scrollIntoView({ behavior: "smooth" })
  }
}
</script>
