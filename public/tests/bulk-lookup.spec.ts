import { expect, test } from "./setup"

test("bulk lookup returns all results", async ({ page, requests }) => {
  await page.goto("/lookup/bulk")

  const input = page.locator("textarea")
  await expect(input).toBeVisible()

  // Fill in all sample data, and search.
  await input.fill(requests.map((r) => r.body.query).join("\n"))
  await expect(page.locator(`text=${requests.length} addresses`)).toBeVisible()

  await page.locator('button:has-text("search")').click()

  // Progress bar works.
  await expect(page.locator("#bulk-progress")).toHaveText(`${requests.length}/${requests.length}`)

  // Results are visible.
  for (const request of requests) {
    await expect(page.locator("#aggregate-country")).toContainText(request.body.country)
    await expect(page.locator("#aggregate-continent")).toContainText(request.body.continent)
    await expect(page.locator("#aggregate-asn")).toContainText(request.body.asn_org)
  }

  // Rate limit counter was updated.
  await expect(page.locator(`text=${100 - requests.length}% calls left`)).toBeVisible()

  // Confirm download points to an object URL, and is downloadable.
  await expect(page.locator("#bulk-results-download")).toHaveAttribute("href", /^blob:http.*/)
  const [download] = await Promise.all([
    page.waitForEvent("download"),
    page.locator("#bulk-results-download").click(),
  ])
  await expect(download.suggestedFilename()).toMatch(/.*\.json$/)

  // Confirm clearing results empties the results.
  await page.locator("#bulk-clear").click()
  await expect(page.locator("#aggregate-country")).not.toBeVisible()
  await expect(page.locator("#aggregate-continent")).not.toBeVisible()
  await expect(page.locator("#aggregate-asn")).not.toBeVisible()

  // Confirm hitting reset on the textarea clears the input.
  await page.locator("text=reset").click()
  await expect(page.locator("textarea")).toBeEmpty()
})
