import { defineConfig, devices } from "@playwright/test";

export default defineConfig({
  testDir: "./tests",
  outputDir: "./tests/results/",
  webServer: {
    command: "bun run build && bun run preview",
    port: 8081,
    reuseExistingServer: !process.env.CI,
  },
  use: {
    baseURL: "http://localhost:8081",
    trace: "on-first-retry",
    bypassCSP: true,
    screenshot: "only-on-failure",
    serviceWorkers: "block",
  },
  expect: {
    timeout: 7000,
  },
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 3 : 0,
  reporter: process.env.CI ? "github" : "list",
  projects: [
    {
      name: "chromium",
      use: {
        ...devices["Desktop Chrome"],
        // navigator.clipboard in headless (firefox does not support these permissions)
        permissions: ["clipboard-read", "clipboard-write"],
      },
    },
    {
      name: "firefox",
      use: {
        ...devices["Desktop Firefox"],
      },
    },
  ],
});
