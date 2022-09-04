<route>
meta:
  title: Not found
</route>

<template>
  <LayoutDefault>
    <div class="p-10">
      <n-result status="404" title="Page not found" description="You know life is always ridiculous.">
        <template #footer>
          <n-button-group>
            <n-button @click="$router.back()">
              <n-icon class="mr-1"><i-mdi-undo-variant /></n-icon>
              Go back
            </n-button>
            <n-button @click="$router.push('/')">
              <n-icon class="mr-1"><i-mdi-magnify /></n-icon>
              Lookup
            </n-button>
          </n-button-group>
        </template>
      </n-result>
    </div>
  </LayoutDefault>
</template>
