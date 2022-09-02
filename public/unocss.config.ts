import { defineConfig, presetTypography, presetUno, transformerDirectives } from "unocss"

export default defineConfig({
  presets: [
    presetUno({
      dark: "media",
    }),
    presetTypography(),
  ],
  transformers: [transformerDirectives()],
})
