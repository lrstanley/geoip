import { defineConfig, presetTypography, presetUno, transformerDirectives } from "unocss"

export default defineConfig({
  presets: [
    presetUno({
      dark: "media",
    }),
    presetTypography(),
  ],
  transformers: [transformerDirectives()],
  safelist:
    "prose dark:prose-invert max-w-full mx-auto text-left text-dark-900 dark:text-white/82 p-4".split(
      " "
    ),
})
