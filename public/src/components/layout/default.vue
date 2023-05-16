<template>
  <div
    class="w-4xl max-h-full mx-auto mt-3 grid basis-auto shrink grow-0 gap-2 grid-cols-1 lg:mt-32 lg:grid-cols-[245px_minmax(0,_1fr)]"
  >
    <CoreNavigation />

    <div class="flex flex-col mx-3">
      <AnimateFade v-if="error">
        <n-alert title="An error occurred" type="error" class="m-2 md:m-6">
          {{ error }}
        </n-alert>
      </AnimateFade>

      <n-card v-else class="shadow drop-shadow-xl rounded-md" content-style="padding: 0;">
        <Suspense>
          <slot />
          <template #fallback>
            <n-spin class="flex flex-auto justify-center items-center p-20">
              <template #description>Loading...</template>
            </n-spin>
          </template>
        </Suspense>
      </n-card>

      <div class="py-4 px-2 mb-3 text-right text-sm lg:mb-20">
        Geo data from
        <a target="_blank" href="http://www.maxmind.com">Maxmind</a>
        &middot; GeoIP:
        <a target="_blank" href="https://github.com/lrstanley/geoip">FOSS</a>
        lookup service, made with
        <i-mdi-heart class="text-red-500" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const error = ref(null)

onErrorCaptured((e) => {
  error.value = e
  return true
})

const router = useRouter()

watch(router.currentRoute, () => {
  error.value = null
})
</script>
