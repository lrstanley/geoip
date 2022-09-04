import type { App } from "vue"

const focus = {
  mounted: (el: HTMLElement) =>
    setTimeout(() => el.querySelector<HTMLElement>("input, textarea")?.focus(), 300),
}

export function registerDirectives(app: App) {
  app.directive("focus", focus)
}
