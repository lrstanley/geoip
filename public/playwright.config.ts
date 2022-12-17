import type { PlaywrightTestConfig } from "@playwright/test"
import { devices } from "@playwright/test"

const config: PlaywrightTestConfig = {
  testDir: "./tests",
  outputDir: "./tests/results/",
  webServer: {
    command: "cd .. && make node-preview",
    port: 8081,
    reuseExistingServer: !process.env.CI,
  },
  use: {
    baseURL: "http://localhost:8081",
    trace: "on-first-retry",
    bypassCSP: true,
    screenshot: "only-on-failure",
  },
  expect: {
    timeout: 7000,
  },
  // repeatEach: 3,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 3 : 0,
  // workers: process.env.CI ? 1 : undefined,
  reporter: process.env.CI ? "github" : "list",
  projects: [
    {
      name: "chromium",
      use: {
        ...devices["Desktop Chrome"],
      },
    },
    {
      name: "firefox",
      use: {
        ...devices["Desktop Firefox"],
      },
    },
  ],
}

export default config
