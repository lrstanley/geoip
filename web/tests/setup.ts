import { test as base } from "@playwright/test";

export { expect } from "@playwright/test";

import type { GeoResult } from "../src/api/types.gen";
import type { APIRequest } from "./data";
import { requests as defaultRequests, openapi, toBulk } from "./data";

export const test = base.extend<{
  requests: APIRequest<GeoResult>[];
}>({
  requests: [defaultRequests, { option: true }],
  page: async ({ page, requests }, providePage) => {
    // Playwright uses the last matching route; register bulk after catch-all.
    await page.route("**/*", async (route, request) => {
      const url = new URL(request.url());
      const path = url.pathname;

      if (path === "/api/v2/openapi.yaml" && request.method() === "GET") {
        await route.fulfill({
          status: 200,
          body: openapi,
          contentType: "application/yaml",
        });
        return;
      }

      if (path === "/api/v2/lookup/self" && request.method() === "GET") {
        const r = requests[0];
        if (!r) {
          await route.continue();
          return;
        }
        await route.fulfill({
          status: 200,
          json: r.body,
          headers: r.headers,
        });
        return;
      }

      for (const req of requests) {
        if (path === `/api/v2/lookup/${req.body.ip}` && request.method() === "GET") {
          await route.fulfill({
            status: 200,
            json: req.body,
            headers: req.headers,
          });
          return;
        }
        if (req.body.host && path === `/api/v2/lookup/${req.body.host}` && request.method() === "GET") {
          await route.fulfill({
            status: 200,
            json: req.body,
            headers: req.headers,
          });
          return;
        }
      }

      await route.continue();
    });

    // String without leading * is resolved with test `baseURL` (see playwright.config baseURL)
    await page.route("/api/v2/bulk", async (route, request) => {
      if (request.method() !== "POST") {
        await route.continue();
        return;
      }
      const bulk = toBulk(requests);
      await route.fulfill({
        status: 200,
        json: bulk.body,
        headers: bulk.headers,
      });
    });

    await providePage(page);
  },
});
