<route>
meta:
  title: Lookup
</route>

<template>
  <div class="transition-all ease-in-out duration-500">
    <div class="p-4" :class="{ 'my-10': history.length < 1 }">
      <n-input
        v-model:value="address"
        v-focus
        type="text"
        size="large"
        placeholder="Search IP address (e.g 1.2.3.4) or host (e.g google.com)"
        :loading="loading || undefined"
        @blur="resultError = null"
        @keyup.enter="search"
      >
        <template v-if="!loading" #suffix>
          <n-icon @click="search"><i-mdi-search /></n-icon>
        </template>
      </n-input>

      <AnimateFade v-if="resultError">
        <n-alert type="error" class="mt-2 alert-small" :show-icon="false">
          error: {{ resultError }}
        </n-alert>
      </AnimateFade>
    </div>

    <n-divider v-show="history.length > 0" class="m-0!" />

    <div v-show="history.length > 0" class="p-4">
      <div class="flex flex-row mb-4">
        <RateLimitCounter size="small" />
        <n-button size="tiny" type="primary" class="ml-auto" @click="clearHistory()">
          <n-icon><i-mdi-broom /></n-icon>
          clear history
        </n-button>
      </div>

      <AnimateListGroup>
        <GeoObject v-for="item of history" :key="item.query" :value="item" />
      </AnimateListGroup>
    </div>
  </div>
</template>

<script setup async lang="ts">
import { api, saveResult, getError } from "@/lib/api"
import vFocus from "@/lib/directives/focus"

const router = useRouter()
const route = useRoute()

const state = useState()
const address = ref<string>("")
const loading = ref<boolean>(false)
const resultError = ref<string | null>()

async function search() {
  if (!address.value) {
    return
  }

  loading.value = true
  resultError.value = null

  // Update the route query param.
  router.push({ query: { q: address.value } })

  try {
    const result = await api.lookup.getAddress({ address: address.value })

    saveResult(result)
    address.value = ""
  } catch (resp) {
    resultError.value = getError(resp).error
  } finally {
    loading.value = false
  }
}

const history = computed(() => {
  return state.history.slice().reverse()
})

function clearHistory() {
  state.history = []
}

onMounted(() => {
  if (route.query.q && route.query.q.length > 0) {
    address.value = typeof route.query.q === "string" ? route.query.q : route.query.q[0]
    search()
  }
})
</script>
