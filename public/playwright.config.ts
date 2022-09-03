import type { PlaywrightTestConfig } from "@playwright/test"
import { devices } from "@playwright/test"

const config: PlaywrightTestConfig = {
  testDir: "./tests",
  outputDir: "./tests/results/",
  webServer: {
    command: "cd .. && make node-preview",
    port: 8081,
    timeout: 30 * 1000,
    reuseExistingServer: !process.env.CI,
  },
  use: {
    actionTimeout: 0,
    baseURL: "http://localhost:8081",
    trace: "on-first-retry",
    bypassCSP: true,
    screenshot: "only-on-failure",
  },
  timeout: 30 * 1000,
  expect: {
    timeout: 7000,
  },
  fullyParallel: true,
  // repeatEach: 3,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 3 : 0,
  workers: process.env.CI ? 1 : undefined,
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
    {
      name: "webkit",
      use: {
        ...devices["Desktop Safari"],
      },
    },
    /* Test against mobile viewports. */
    // {
    //   name: 'Mobile Chrome',
    //   use: {
    //     ...devices['Pixel 5'],
    //   },
    // },
    // {
    //   name: 'Mobile Safari',
    //   use: {
    //     ...devices['iPhone 12'],
    //   },
    // },
  ],
}

export default config
