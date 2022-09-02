import type { GeoResult } from "@/lib/api/openapi"
import type { FieldGroup, TruncateOptions, GroupWidth } from "@/lib/api/types"

/**
 * groupByField groups GeoResult by a given field and label on the GeoResult object,
 * allowing truncating low-grouped values into a "other" group. Returned results are
 * sorted by the count per group, descending, and the "other" group (if configured)
 * is last.
 *
 * @export
 * @param {GeoResult[]} results
 * @param {keyof GeoResult} field
 * @param {keyof GeoResult} label
 * @param {?TruncateOptions} [truncate]
 * @returns {FieldGroup[]}
 */
export function groupByField(
  results: GeoResult[],
  field: keyof GeoResult,
  label: keyof GeoResult,
  truncate?: TruncateOptions
): FieldGroup[] {
  const calc: { [key: string]: FieldGroup } = {}

  for (const result of results) {
    const fieldVal = result[field]
    const labelVal = result[label]
    const key = fieldVal.toString()

    if (!calc[key]) {
      calc[key] = {
        count: 1,
        percent: 0,
        field: fieldVal,
        label: labelVal,
      }
    } else {
      calc[key].count++
    }
  }

  // Calculate percentages.
  for (const key in calc) {
    calc[key].percent = Math.round((calc[key].count / results.length) * 100)
  }

  // Sort by count.
  const val: FieldGroup[] = Object.values(calc).sort((a, b) => b.count - a.count)

  // Truncate all low-percentage results into "Other".
  if (truncate) {
    const other: FieldGroup = {
      count: 0,
      percent: 0,
      field: truncate.label,
      label: truncate.label,
    }

    // Iterate backwards so we can remove items from the array easily.
    for (let i = val.length - 1; i >= 0; i--) {
      if (val[i].label == truncate.label) {
        continue
      }

      if (truncate.percent && val[i].percent < truncate.percent) {
        other.count += val[i].count
        val.splice(i, 1)
      } else if (truncate.count && val[i].count < truncate.count) {
        other.count += val[i].count
        val.splice(i, 1)
      }
    }

    if (other.count > 0) {
      other.percent = Math.round((other.count / results.length) * 100)
      val.push(other)
    }
  }

  return val
}

/**
 * calculateGroupWidth calculates the width of the count and label fields of an
 * array of FieldGroup objects. Percentage does not include the percent sign.
 *
 * @export
 * @param {FieldGroup[]} stats
 * @returns {GroupWidth}
 */
export function calculateGroupWidth(stats: FieldGroup[], maxLabelWidth = 30): GroupWidth {
  const label = stats.reduce((max, stat) => Math.max(max, stat.label.toString().length), 0)
  const count = stats.reduce((max, stat) => Math.max(max, stat.count.toString().length), 0)
  const percent = stats.reduce((max, stat) => Math.max(max, stat.percent.toString().length), 0)

  return {
    labelWidth: Math.min(label, maxLabelWidth),
    countWidth: count,
    percentWidth: percent,
  }
}
