import { expect, test } from "./setup";

test("validate api docs", async ({ page, requests }) => {
  await page.goto("/lookup/docs#get-/lookup/-address-");
  await expect(page).toHaveTitle(/Documentation.*/);

  await expect(page.getByRole("heading", { name: /lookup address/i })).toBeVisible({
    timeout: 15000,
  });

  const input = page.locator('api-request input[type="text"]').first();
  await input.focus();
  const first = requests[0];
  if (!first) {
    throw new Error("fixture missing request");
  }
  await input.fill(first.body.query);
  await page.locator('button:has-text("TRY")').click();

  await expect(page.locator("text=Response Status: OK:200")).toBeVisible({
    timeout: 15000,
  });
});
