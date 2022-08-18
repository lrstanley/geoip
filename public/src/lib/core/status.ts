import { createDiscreteApi, darkTheme } from "naive-ui"

import type { ComputedRef } from "vue"
import type { ConfigProviderProps } from "naive-ui"

const dark = usePreferredDark()

export const theme = computed(() => (dark.value ? darkTheme : null))

export const configProvider: ComputedRef<ConfigProviderProps> = computed(() => {
  return {
    theme: theme.value,
    abstract: true,
    preflightStyleDisabled: true,
  }
})

export const { loadingBar, notification } = createDiscreteApi(
  ["loadingBar", "notification"], // "message", "dialog"
  { configProviderProps: configProvider }
)
