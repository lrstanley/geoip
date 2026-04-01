import path from "node:path";
import { heyApiPlugin } from "@hey-api/vite-plugin";
import tailwindcss from "@tailwindcss/vite";
import { tanstackRouter } from "@tanstack/router-plugin/vite";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";
import biomePlugin from "vite-plugin-biome";

export default defineConfig({
  plugins: [
    heyApiPlugin({
      config: {
        input: path.resolve(__dirname, "../internal/handlers/apihandler/openapi_v2.yaml"),
        output: {
          path: path.resolve(__dirname, "src/api"),
          entryFile: false,
        },
        plugins: ["@hey-api/client-ofetch"],
      },
    }),
    tanstackRouter({
      target: "react",
      autoCodeSplitting: true,
    }),
    react(),
    tailwindcss(),
    biomePlugin({
      applyFixes: true,
      hotUpdateMode: "changed",
      mode: "lint",
    }),
  ],
  resolve: {
    tsconfigPaths: true,
  },
  publicDir: path.resolve(__dirname, "./src/assets"),
  base: "/",
  build: {
    sourcemap: "hidden",
    emptyOutDir: true,
  },
  preview: {
    port: 8081,
    open: false,
    strictPort: true,
    proxy: {
      "^/(api|security\\.txt|robots\\.txt)(/.*|$)": {
        target: "http://localhost:8080",
        changeOrigin: true,
        xfwd: true,
      },
    },
  },
  server: {
    port: 8081,
    open: false,
    strictPort: true,
    proxy: {
      "^/(api|security\\.txt|robots\\.txt)(/.*|$)": {
        target: "http://localhost:8080",
        changeOrigin: true,
        xfwd: true,
      },
    },
  },
  optimizeDeps: {
    include: ["react", "react-dom", "ofetch"],
  },
});
