import type { GeoResult, BulkGeoResult } from "@/lib/api"
import { readFileSync } from "fs"
import { resolve } from "path"

export { expect } from "@playwright/test"

export const openapi = readFileSync(
  resolve(__dirname, "../../internal/handlers/apihandler/openapi_v2.yaml"),
  {
    encoding: "utf8",
  }
)

export interface APIRequest<T> {
  body: T
  headers: { [key: string]: string }
}

export const requests: APIRequest<GeoResult>[] = [
  {
    body: {
      accuracy_radius_km: 1000,
      asn: "AS15169",
      asn_org: "GOOGLE",
      city: "",
      continent: "North America",
      continent_abbr: "NA",
      country: "United States",
      country_abbr: "US",
      host: "dns.google",
      ip: "8.8.8.8",
      ip_type: 4,
      latitude: 37.751,
      longitude: -97.822,
      network: "8.8.8.0/24",
      postal_code: "",
      query: "8.8.8.8",
      subdivision: "",
      summary: "United States, NA",
      timezone: "America/Chicago",
    },
    headers: {
      "Content-Type": "application/json",
      "X-Ratelimit-Limit": "100",
      "X-Ratelimit-Remaining": "99",
      "X-Ratelimit-Reset": "1659751200",
    },
  },
  {
    body: {
      accuracy_radius_km: 20,
      asn: "AS16276",
      asn_org: "OVH SAS",
      city: "Beauharnois",
      continent: "North America",
      continent_abbr: "NA",
      country: "Canada",
      country_abbr: "CA",
      host: "ovh.ca",
      ip: "66.70.178.39",
      ip_type: 4,
      latitude: 45.3161,
      longitude: -73.8736,
      network: "66.70.128.0/17",
      postal_code: "J6N",
      query: "ovh.ca",
      subdivision: "Quebec",
      summary: "Beauharnois, Quebec, CA",
      timezone: "America/Toronto",
    },
    headers: {
      "Content-Type": "application/json",
      "X-Ratelimit-Limit": "100",
      "X-Ratelimit-Remaining": "98",
      "X-Ratelimit-Reset": "1659751200",
    },
  },
  {
    body: {
      accuracy_radius_km: 20,
      asn: "AS14061",
      asn_org: "DIGITALOCEAN-ASN",
      city: "Singapore",
      continent: "Asia",
      continent_abbr: "AS",
      country: "Singapore",
      country_abbr: "SG",
      ip: "143.198.200.155",
      ip_type: 4,
      latitude: 1.3078,
      longitude: 103.6818,
      network: "143.198.192.0/19",
      postal_code: "62",
      query: "143.198.200.155",
      subdivision: "",
      summary: "Singapore, SG",
      timezone: "Asia/Singapore",
    },
    headers: {
      "Content-Type": "application/json",
      "X-Ratelimit-Limit": "100",
      "X-Ratelimit-Remaining": "97",
      "X-Ratelimit-Reset": "1659751200",
    },
  },
]

export function toBulk(requests: APIRequest<GeoResult>[]): APIRequest<BulkGeoResult> {
  return {
    headers: requests[requests.length - 1].headers,
    body: {
      errors: [],
      results: requests.map((r) => r.body),
    },
  }
}

export function toResponse<T>(request: APIRequest<T>, status = 200) {
  return {
    status: status,
    body: JSON.stringify(request.body, null, 4),
    headers: request.headers,
  }
}
