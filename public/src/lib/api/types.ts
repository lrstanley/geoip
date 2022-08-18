/**
 * Config are baseline configurations for the HTTP API client.
 *
 * @export
 * @interface Config
 * @typedef {Config}
 */
export interface Config {
  headers: Headers
  baseUrl: string
}

/**
 * GeoIPData is the geoip results of a given lookup from the API.
 *
 * @export
 * @interface GeoIPData
 * @typedef {GeoIPData}
 */
export interface GeoIPData {
  ip: string
  summary: string
  city: string
  subdivision: string
  country: string
  country_abbr: string
  continent: string
  continent_abbr: string
  latitude: number
  longitude: number
  timezone: string
  postal_code: string
  proxy: boolean
  host: string
}

/**
 * ClientState is the state of the client, including rate limit and geoip database
 * information.
 *
 * Marshalled from header data:
 *   X-Cache: MISS
 *   X-Maxmind-Build: 6-1659368358
 *   X-Maxmind-Type: GeoLite2-City
 *   X-Ratelimit-Limit: 1000000
 *   X-Ratelimit-Remaining: 999999
 *   X-Ratelimit-Reset: 1659751200
 *
 * @export
 * @interface ClientState
 * @typedef {ClientState}
 */
export interface ClientState {
  cached: boolean
  maxmind_build: string
  maxmind_type: string
  ratelimit_limit: number
  ratelimit_remaining: number
  ratelimit_reset: number
}

/**
 * APIResponse wrapper object for the API.
 *
 * @export
 * @interface APIResponse
 * @typedef {APIResponse}
 */
export interface APIResponse {
  query: string
  data: GeoIPData
  state: ClientState
  error: Error
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
export type FieldGroup = ByField<ValueOf<GeoIPData>, ValueOf<GeoIPData>>

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
