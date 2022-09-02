<template>
  <div v-bind="$attrs">
    <AnimateListGroup>
      <div
        v-for="stat in group.stats"
        :key="stat.field.toString()"
        class="flex flex-auto gap-2 items-center my-1"
      >
        <div
          class="shrink-0 text-right truncate"
          :style="{ width: group.labelWidth + 'ch' }"
          :title="
            props.truncate?.label == stat.label
              ? `${props.truncate.label} represents aggregated groups with less than ${
                  props.truncate.count || props.truncate.percent + '%'
                } total results each`
              : stat.label.toString()
          "
        >
          {{ stat.label }}
        </div>
        <GeoFlag
          v-if="props.flag"
          :value="stat.field.toString()"
          :size="16"
          class="inline-flex"
          immediate
        />
        <n-progress
          type="line"
          :percentage="stat.percent"
          :status="props.status || 'info'"
          :show-indicator="false"
        />
        <div class="shrink-0 hidden md:block" :style="{ width: group.countWidth + 'ch' }">
          {{ stat.count }}
        </div>
        <div class="shrink-0 hidden md:block" :style="{ width: group.percentWidth + 3 + 'ch' }">
          ({{ stat.percent }}%)
        </div>
      </div>
    </AnimateListGroup>
  </div>
</template>

<script setup lang="ts">
import { groupByField, calculateGroupWidth } from "@/lib/api"
import type { GeoResult, TruncateOptions } from "@/lib/api"
import type { Status as ProgressStatus } from "naive-ui/es/progress/src/interface"

const props = defineProps<{
  value: GeoResult[]
  field: keyof GeoResult
  label: keyof GeoResult
  truncate?: TruncateOptions
  flag?: boolean
  status?: ProgressStatus
}>()

const group = computed(() => {
  const stats = groupByField(props.value, props.field, props.label, props.truncate)

  return { stats: stats, ...calculateGroupWidth(stats, 20) }
})
</script>
