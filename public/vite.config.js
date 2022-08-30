import { defineConfig } from "vite"

import path from "path"
import Vue from "@vitejs/plugin-vue"
import AutoImport from "unplugin-auto-import/vite"
import Components from "unplugin-vue-components/vite"
import { NaiveUiResolver, VueUseComponentsResolver } from "unplugin-vue-components/resolvers"
import Pages from "vite-plugin-pages"
import Icons from "unplugin-icons/vite"
import IconsResolver from "unplugin-icons/resolver"
import Unocss from "unocss/vite"
import Layouts from "vite-plugin-vue-layouts"
import { visualizer } from "rollup-plugin-visualizer"

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
    Layouts(),
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
    mode: "production",
  },
  preview: {
    port: 8081,
    mode: "production",
    open: false,
  },
  server: {
    base: "/",
    mode: "development",
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
