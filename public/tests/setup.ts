import { test as base } from "@playwright/test"
import requests from "./data/requests.json"

export { expect } from "@playwright/test"

export interface Request {
  body: Response
  headers: { [key: string]: string }
}

export interface Response {
  city: string
  continent: string
  continent_abbr: string
  country: string
  country_abbr: string
  host: string
  ip: string
  latitude: number
  longitude: number
  postal_code: string
  proxy: boolean
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
      page.route(`/api/${request.body.ip}`, (route) => route.fulfill(requestResponse(request)))
      page.route(`/api/${request.body.host}`, (route) => route.fulfill(requestResponse(request)))
    }

    // Also register /api/self, and make a catch-all for other items.
    page.route("/api/self", (route) => route.fulfill(requestResponse(requests[0])))

    await use(page)
  },
})
