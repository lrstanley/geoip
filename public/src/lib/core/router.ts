import { setupLayouts } from "virtual:generated-layouts"
import { createRouter, createWebHistory } from "vue-router"
import generatedRoutes from "~pages"
import { loadingBar } from "@/lib/core/status"

const routes = setupLayouts(generatedRoutes)

const router = createRouter({
  history: createWebHistory("/"),
  routes,
})

router.beforeEach(async (to, from, next) => {
  if (from.name != to.name || JSON.stringify(from.params) != JSON.stringify(to.params)) {
    loadingBar.start()
  }

  next()
})

router.afterEach((to) => {
  document.title = `${to.meta.title} · GeoIP`

  nextTick(() => {
    loadingBar.finish()

    if (location.hash && !to.meta.disableAnchor) {
      const el = document.getElementById(location.hash.slice(1))
      if (el) {
        el.scrollIntoView()
      }
    }
  })
})

export default router
