import { expect, test } from "@playwright/test"

test("test", async ({ page }) => {
  await page.goto("/lookup/bulk")
  await page.locator("textarea").click()
  await page.locator("textarea").fill("8.8.8.8")
  await page.locator('button:has-text("search")').click()
  await page.locator("div:nth-child(8)").first().click()
  await page.locator("#app").click()
  await page.locator('button:has-text("clear")').click()
  await page.locator("text=reset").click()
})
