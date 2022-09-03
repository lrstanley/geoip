import { defineStore } from "pinia"
import { useStorage } from "@vueuse/core"

import type { GeoResult, ClientState } from "@/lib/api"

const version = 4

export const useState = defineStore("state", {
  state: () => {
    return useStorage(`state-${version}`, {
      clientState: {} as ClientState,
      history: [] as GeoResult[],
    })
  },
  getters: {
    hasSelf(state) {
      return state.history.some((item) => item.query == "self")
    },
  },
})
