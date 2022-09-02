import "@unocss/reset/normalize.css"
import "uno.css"
import "@/css/main.css"
import { createPinia } from "pinia"
import { createApp } from "vue"
import router from "@/lib/core/router"
import App from "@/main.vue"
import { MotionPlugin } from "@vueuse/motion"

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(MotionPlugin)
app.mount("#app")
