import { expect, test } from "./setup"

test("validate api docs", async ({ page, requests }) => {
  await page.goto("/lookup/docs#get-/lookup/-address-")
  await expect(page).toHaveTitle(/Documentation .*/)
  await expect(page.locator('h2:has-text("Lookup address")')).toBeVisible()

  const input = page.locator('api-request input[type="text"] >> nth=0')
  await input.focus()
  await input.fill(requests[0].body.query)
  await page.locator('button:has-text("TRY")').click()

  await expect(page.locator("text=Response Status: OK:200")).toBeVisible()
})
