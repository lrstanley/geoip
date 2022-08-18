import { defineStore } from "pinia"
import { useStorage } from "@vueuse/core"

import type { APIResponse, ClientState } from "@/lib/api"

const version = 4

export const useState = defineStore("state", {
  state: () => {
    return useStorage(`state-${version}`, {
      clientState: {} as ClientState,
      history: [] as APIResponse[],
    })
  },
  getters: {
    hasSelf(state) {
      return state.history.some((item) => item.query == "self")
    },
  },
})
