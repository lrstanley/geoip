import { readFileSync } from "fs"
import { resolve } from "path"
import { test as base } from "@playwright/test"
import requests from "./data/requests.json"

const openapi = readFileSync(resolve(__dirname, "../../internal/handlers/apihandler/openapi_v2.yaml"), {
  encoding: "utf8",
})

export { expect } from "@playwright/test"

export interface Request {
  body: Response
  headers: { [key: string]: string }
}

export interface Response {
  accuracy_radius_km: number
  asn: string
  asn_org: string
  city: string
  continent: string
  continent_abbr: string
  country: string
  country_abbr: string
  host?: string
  ip: string
  ip_type: number
  latitude: number
  longitude: number
  network: string
  postal_code: string
  query: string
  subdivision: string
  summary: string
  timezone: string
}

export type TestOptions = {
  requests: Request[]
}

function requestResponse(request: any) {
  return {
    status: 200,
    body: JSON.stringify(request.body, null, 4),
    headers: request.headers,
  }
}

export const test = base.extend<TestOptions>({
  requests,
  page: async ({ page }, use) => {
    // Register both IP and domain/host lookup methods.
    for (const request of requests) {
      await page.route(
        (url) => url.pathname == `/api/v2/lookup/${request.body.ip}`,
        (route) => route.fulfill(requestResponse(request))
      )
      await page.route(
        (url) => url.pathname == `/api/v2/lookup/${request.body.host}`,
        (route) => route.fulfill(requestResponse(request))
      )
    }

    await page.route("/api/v2/openapi.yaml", (route) =>
      route.fulfill({
        status: 200,
        body: openapi,
        contentType: "application/yaml",
      })
    )

    // /api/v2/lookup/self and others
    await page.route(
      (url) => url.pathname == "/api/v2/lookup/self",
      (route) => route.fulfill(requestResponse(requests[0]))
    )

    await page.route(
      (url) => url.pathname == "/api/v2/bulk",
      (route) =>
        route.fulfill({
          status: 200,
          headers: requests[requests.length - 1].headers,
          contentType: "application/json",
          body: JSON.stringify({ results: requests.map((r) => r.body), errors: [] }, null, 4),
        })
    )

    await use(page)
  },
})
