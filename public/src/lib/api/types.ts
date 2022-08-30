import type { GeoResult } from "@/lib/api/openapi"

/**
 * ClientState is the state of the client, including rate limit.
 *
 * Marshalled from header data:
 *   X-Ratelimit-Limit: 1000000
 *   X-Ratelimit-Remaining: 999999
 *   X-Ratelimit-Reset: 1659751200
 *
 * @export
 * @interface ClientState
 * @typedef {ClientState}
 */
export interface ClientState {
  ratelimit_limit: number
  ratelimit_remaining: number
  ratelimit_reset: number
}

/**
 * ValueOf is a type guard for values of a given interface/type representing K:V
 * objects.
 *
 * @export
 * @typedef {ValueOf}
 * @template T
 */
export type ValueOf<T> = T[keyof T]

/**
 * ByField is used for grouping and sorting an input by a given field and label.
 *
 * @export
 * @interface ByField
 * @typedef {ByField}
 * @template Field
 * @template Label
 */
export interface ByField<Field, Label> {
  count: number
  percent: number
  field: Field
  label: Label
}

/**
 * ByField is used for grouping and sorting GeoIPData by a given field and label.
 *
 * @export
 * @typedef {FieldGroup}
 */
export type FieldGroup = ByField<ValueOf<GeoResult>, ValueOf<GeoResult>>

/**
 * TruncateOptions is used by grouping functions to configure how (if at all)
 * truncating low-volume results into a "other" type group.
 *
 * @export
 * @interface TruncateOptions
 * @typedef {TruncateOptions}
 */
export interface TruncateOptions {
  label: string
  percent?: number
  count?: number
}

export interface GroupWidth {
  labelWidth: number
  percentWidth: number
  countWidth: number
}
