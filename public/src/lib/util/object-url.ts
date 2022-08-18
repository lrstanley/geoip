import { computed } from "vue"
import { useObjectUrl } from "@vueuse/core"

import type { Ref } from "vue"

export function createJSONObjectURL(data: Ref<any>): Readonly<Ref<string>> {
  return useObjectUrl(
    computed(() => {
      if (!data.value) return null

      return new Blob([JSON.stringify(data.value, null, 4)], {
        type: "application/json",
      })
    })
  )
}
