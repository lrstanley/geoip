import { expect, test } from "./setup"

test("validate home lookup IP", async ({ page, requests }) => {
  await page.goto("/")

  const input = page.locator('input[type="text"]')

  for (const request of requests) {
    await input.focus()
    await input.fill(request.body.query)
    await input.press("Enter")

    const result = page.locator("div.geo-result", { hasText: request.body.query })
    await expect(result).toBeVisible()

    // Validate we can click the query and copy to clipboard.
    await result.locator(`button:has-text("${request.body.query}") >> nth=0`).click()
    await expect(
      page.locator("div.n-notification", { hasText: `copied "${request.body.query}"` })
    ).toBeVisible()

    await result.locator('[aria-label="Zoom out"]').click()

    // Validate the flag image is working properly.
    await expect(result.locator("button img")).toHaveAttribute(
      "src",
      new RegExp(`.*/flags/${request.body.country_abbr.toLowerCase()}.svg`)
    )
  }

  // Validate that the rate limit info is working.
  await expect(page.locator("div.n-tag", { hasText: "calls left" })).toBeVisible()

  await page.locator('button:has-text("clear history")').click()
  await expect(page.locator("div.geo-result")).toHaveCount(0) // all results were removed.
})
