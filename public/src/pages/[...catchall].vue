<template>
  <div class="flex flex-col justify-center flex-auto">
    <n-result
      :status="errorCode"
      :title="'Error code: ' + errorTitle"
      :description="props.error ? props.error.toString() : 'You know life is always ridiculous.'"
    >
      <template #footer>
        <n-button-group>
          <n-button @click="$router.back()">
            <n-icon class="mr-1"><i-mdi-undo-variant /></n-icon>
            Go back
          </n-button>
          <n-button @click="$router.push('/')">
            <n-icon class="mr-1"><i-mdi-home /></n-icon>
            Home
          </n-button>
        </n-button-group>
      </template>
    </n-result>
  </div>
</template>

<script setup>
const props = defineProps({
  catchall: {
    type: [String, Array],
    default: "",
  },
  error: {
    type: Error,
    default: null,
  },
})

const source = computed(() => (typeof props.catchall === "string" ? [props.catchall] : props.catchall))
const errorCode = ref("0")
const errorTitle = ref("")
const supported = ["info", "success", "warning", "error", "404", "403", "500", "418"]

onMounted(() => {
  for (const item of source.value) {
    if (supported.includes(item)) {
      errorCode.value = item
    } else if (item.match(/^[45][0-9]+$/)) {
      errorCode.value = "error"
    } else if (item == "CombinedError") {
      errorCode.value = "error"
      errorTitle.value = "query error"
    } else {
      errorCode.value = "404"
      errorTitle.value = "not found"
      break
    }
  }

  if (errorTitle.value == "") {
    errorTitle.value = errorCode.value
  }
})
</script>
