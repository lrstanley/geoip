import path from "path"
import { visualizer } from "rollup-plugin-visualizer"
import Unocss from "unocss/vite"
import AutoImport from "unplugin-auto-import/vite"
import IconsResolver from "unplugin-icons/resolver"
import Icons from "unplugin-icons/vite"
import { NaiveUiResolver, VueUseComponentsResolver } from "unplugin-vue-components/resolvers"
import Components from "unplugin-vue-components/vite"
import { defineConfig } from "vite"
import Pages from "vite-plugin-pages"
import Vue from "@vitejs/plugin-vue"

export default defineConfig({
  resolve: {
    alias: {
      "@/": `${path.resolve("src")}/`,
    },
  },
  publicDir: `${path.resolve("src")}/assets`,
  plugins: [
    visualizer({
      filename: "./dist/stats.html",
    }),
    Pages({
      extensions: ["vue", "md"],
      dirs: "src/pages",
      routeBlockLang: "yaml",
    }),
    Vue({
      include: [/\.vue$/, /\.md$/],
      template: {
        compilerOptions: {
          isCustomElement: (tag) => ["rapi-doc"].includes(tag),
        },
      },
    }),
    Components({
      extensions: ["vue", "md"],
      dirs: ["src/components"],
      include: [/\.vue$/, /\.vue\?vue/, /\.md$/],
      directoryAsNamespace: true,
      resolvers: [
        VueUseComponentsResolver(),
        NaiveUiResolver(),
        IconsResolver({ componentPrefix: "i", enabledCollections: ["mdi"] }),
      ],
    }),
    AutoImport({
      dts: true,
      imports: [
        "vue",
        "vue/macros",
        "vue-router",
        "@vueuse/core",
        {
          "@/lib/core/state": ["useState"],
        },
      ],
      resolvers: [IconsResolver({ componentPrefix: "icon", enabledCollections: ["mdi"] })],
      eslintrc: {
        enabled: true,
      },
    }),
    Unocss(),
    Icons({
      autoInstall: true,
      defaultClass: "icon",
    }),
  ],
  base: "/",
  build: {
    sourcemap: "hidden",
    emptyOutDir: true,
  },
  preview: {
    port: 8081,
    open: false,
  },
  server: {
    base: "/",
    port: 8081,
    open: false,
    strictPort: true,
    proxy: {
      "^/(api|security\\.txt|robots\\.txt)(/.*|$)": {
        target: "http://localhost:8080",
        xfwd: true,
      },
    },
  },
  optimizeDeps: {
    include: ["vue", "vue-router", "@vueuse/core"],
  },
})
