import { expect, test } from "@playwright/test"

test("validate home lookup IP", async ({ page }) => {
  await page.goto("/")

  const input = page.locator('input[type="text"]')

  await input.focus()
  await input.fill("8.8.8.8")
  await input.press("Enter")

  const result = page.locator("div.geo-result", { has: page.locator('button:has-text("8.8.8.8")') })
  await expect(result).toBeVisible()

  await result.locator('button:has-text("8.8.8.8")').click()
  await expect(page.locator("div.n-notification", { hasText: 'copied "8.8.8.8"' })).toBeVisible()

  await expect(result.locator('img[alt="Marker"]')).toBeVisible()
  await result.locator('[aria-label="Zoom in"]').click()
  await result.locator('[aria-label="Zoom out"]').click()

  // Validate the flag image is working properly.
  await expect(result.locator("button img")).toHaveAttribute("src", /.*\/flags\/us\.svg$/)

  // Validate that the rate limit info is working.
  await expect(page.locator("div.n-tag", { hasText: "calls left" })).toBeVisible()

  await page.locator('button:has-text("clear history")').click()
  await expect(result).toHaveCount(0) // was deleted.
})
