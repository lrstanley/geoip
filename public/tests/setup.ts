import { test as base } from "@playwright/test"
export { expect } from "@playwright/test"
import { requests, openapi } from "./data"
import type { APIRequest } from "./data"
import { toBulk, toResponse } from "./data"
import type { GeoResult } from "@/lib/api"

export const test = base.extend<{
  requests: APIRequest<GeoResult>[]
}>({
  requests: [requests, { option: true }],
  page: async ({ page }, use) => {
    // Register both IP and domain/host lookup methods.
    for (const request of requests) {
      await page.route(
        (url) => url.pathname == `/api/v2/lookup/${request.body.ip}`,
        (route) => route.fulfill(toResponse(request))
      )
      await page.route(
        (url) => url.pathname == `/api/v2/lookup/${request.body.host}`,
        (route) => route.fulfill(toResponse(request))
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
      (route) => route.fulfill(toResponse(requests[0]))
    )

    await page.route(
      (url) => url.pathname == "/api/v2/bulk",
      (route) => route.fulfill(toResponse(toBulk(requests)))
    )

    await use(page)
  },
})
