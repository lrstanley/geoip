<route>
meta:
  title: Documentation
  layout: bare
  disableAnchor: true
</route>

<template>
  <LayoutBare>
    <rapi-doc
      render-style="focused"
      spec-url="/api/v2/openapi.yaml"
      show-header="false"
      show-info="false"
      show-method-in-nav-bar="as-colored-text"
      nav-item-spacing="relaxed"
      allow-authentication="false"
      allow-server-selection="false"
      allow-spec-url-load="false"
      allow-spec-file-load="false"
      allow-spec-file-download="false"
      allow-api-list-style-selection="false"
      load-fonts="false"
      regular-font="Consolas, monaco, monospace"
      mono-font="Consolas, monaco, monospace"
      font-size="large"
      :server-url="origin + '/api/v2'"
      :default-api-server="origin + '/api/v2'"
      :nav-bg-color="theme.nav.bg"
      :nav-hover-bg-color="theme.nav.hover"
      :bg-color="theme.bg"
      :text-color="theme.text"
      :primary-color="theme.primary"
      class="h-full w-full"
    >
      <div slot="nav-logo">
        <CoreNavigation slim hide-source class="pb-2" />
      </div>
    </rapi-doc>
  </LayoutBare>
</template>

<script setup lang="ts">
import "rapidoc"

const dark = usePreferredDark()

const origin = location.origin

const themeDefaults = {
  dark: {
    nav: {
      bg: "#0F0F0F",
      hover: "#1B1B1B",
    },
    bg: "#18181C",
    text: "#BBBBBB",
    primary: "#10B981",
  },
  light: {
    nav: {
      bg: "#FAFAFA",
      hover: "#E5E7EB",
    },
    bg: "#FFFFFF",
    text: "#444444",
    primary: "#36AD6A",
  },
}

const theme = computed(() => {
  return dark.value ? themeDefaults.dark : themeDefaults.light
})
</script>

<style scoped>
rapi-doc::part(section-navbar-scroll) {
  padding: 0px 15px;
}
</style>
