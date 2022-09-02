<template>
  <aside aria-label="Navigation">
    <div class="overflow-y-auto rounded-md" :class="{ 'py-4 px-3': !props.slim }" v-bind="$attrs">
      <span class="flex items-center pl-2">
        <n-icon>
          <i-mdi-map-search-outline class="text-3xl" />
        </n-icon>

        <span class="self-center ml-2 text-3xl font-semibold whitespace-nowrap text-emerald-500">
          GeoIP
        </span>
      </span>

      <n-divider class="my-2!" />

      <ul class="space-y-2 list-none p-0 m-0">
        <li v-for="route in menuOptions" :key="route.label">
          <component
            :is="route.attrs.href ? 'a' : 'router-link'"
            v-bind="route.attrs"
            active-class="route-active"
            class="route"
          >
            <n-icon>
              <component :is="route.icon" />
            </n-icon>

            <span class="ml-3 text-gradient bg-gradient-to-r from-emerald-500 to-lightblue-500">
              {{ route.label }}
            </span>
          </component>
        </li>
      </ul>

      <template v-if="!props.hideSource">
        <n-divider class="my-2!" />

        <ul class="space-y-2 list-none p-0 m-0">
          <li>
            <a href="https://github.com/lrstanley/geoip" class="route" target="_blank">
              <n-icon>
                <i-mdi-github />
              </n-icon>

              <span class="ml-3 text-gradient bg-gradient-to-r from-emerald-500 to-lightblue-500">
                Github Project
              </span>
            </a>
          </li>
        </ul>
      </template>
    </div>
  </aside>
</template>
<script setup lang="ts">
const props = defineProps<{
  slim?: boolean
  hideSource?: boolean
}>()

interface MenuOption {
  attrs: { [key: string]: any }
  label: string
  icon: string
}

const menuOptions: MenuOption[] = [
  {
    attrs: { to: { name: "index" } },
    label: "Lookup Address",
    icon: IconMdiMagnify,
  },
  {
    attrs: { to: { name: "lookup-bulk" } },
    label: "Bulk Lookup",
    icon: IconMdiDatabaseSearchOutline,
  },
  {
    attrs: { to: { name: "lookup-docs" } },
    label: "API Documentation",
    icon: IconMdiApplicationBracketsOutline,
  },
]
</script>

<style scoped>
.route {
  @apply flex items-center py-2 px-3 rounded;
  @apply transition-all ease-in-out duration-250;
  @apply text-base font-normal decoration-none;
  @apply text-dark-900 dark:text-white;
  @apply hover:bg-gray-200 dark:hover:bg-dark-700;
}
.route-active {
  @apply bg-gray-200 dark:bg-dark-700 dark:text-emerald-500;
}
</style>
