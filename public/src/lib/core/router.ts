import { setupLayouts } from "virtual:generated-layouts"
import { createRouter, createWebHistory } from "vue-router"
import generatedRoutes from "~pages"
import { api, saveResult } from "@/lib/api"
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

  const state = useState()
  if (!state.hasSelf) {
    // Kickoff the request for /api/self to lookup the user, in the background.
    api.lookup
      .getAddress({ address: "self" })
      .then((result) => {
        saveResult(result)
      })
      .catch(() => null)
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