import { expect, test } from "./setup"

test("validate api docs", async ({ page, requests }) => {
  await page.goto("/lookup/docs")
  await expect(page).toHaveTitle(/Documentation .*/)
  await page.locator('h1:has-text("API Documentation")').isVisible()

  const anchor = page.locator("text=Filtering output >> a")
  await anchor.click()
  await expect(page).toHaveURL("/lookup/docs#filtering-output")
  await expect(anchor).toBeVisible()

  await page.locator(`td:has-text("${requests[0].headers["X-Maxmind-Type"]}")`).click()
})
