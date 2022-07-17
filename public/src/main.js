import Vue from "vue"
import App from "@/App.vue"
import VueLocalStorage from "vue-localstorage"
import VueRouter from "vue-router"
import VueResource from "vue-resource"
import VueProgressBar from "vue-progressbar"

// For routes.
import Lookup from "@/components/Lookup.vue"
import BulkLookup from "@/components/BulkLookup.vue"
import Docs from "@/components/Docs.vue"
import NotFound from "@/components/NotFound.vue"

const routes = [
    {
        name: "lookup",
        path: "/",
        component: Lookup,
        meta: { title: "Lookup address" },
    },
    {
        name: "bulkLookup",
        path: "/lookup/bulk",
        component: BulkLookup,
        meta: { title: "Bulk lookup" },
    },
    {
        name: "apidocs",
        path: "/lookup/docs",
        component: Docs,
        meta: { title: "API docs" },
    },
    {
        name: "catchall",
        path: "*",
        redirect: "/404",
    },
    {
        name: "404",
        path: "/404",
        component: NotFound,
        meta: { title: "Page not found" },
    },
]

Vue.use(VueLocalStorage, { name: "ls" })
Vue.use(VueResource)
Vue.use(VueRouter)
Vue.use(VueProgressBar, {
    color: "#0074D9",
    failedColor: "#FF4136",
    thickness: "3px",
    transition: {
        speed: "0.2s",
        opacity: "0.6s",
        termination: 300,
    },
    autoRevert: true,
    location: "top",
    inverse: false,
})

const router = new VueRouter({ routes, mode: "history" })
router.beforeEach((to, from, next) => {
    if (to.meta.title !== undefined) {
        document.title = `${to.meta.title} · GeoIP`
    } else {
        document.title = "GeoIP"
    }

    next()
})

new Vue({ router, el: "#vue", render: (h) => h(App) })
