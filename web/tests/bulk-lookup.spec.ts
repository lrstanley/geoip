import { expect, test } from "./setup";

test("bulk lookup returns all results", async ({ page, requests }) => {
  await page.goto("/lookup/bulk");

  const input = page.locator("textarea");
  await expect(input).toBeVisible();

  await input.fill(requests.map((r) => r.body.query).join("\n"));
  await expect(page.getByText(`${requests.length} addresses`)).toBeVisible();

  await page.getByRole("button", { name: /search/i }).click();

  await expect(page.locator("#bulk-progress")).toContainText(`${requests.length}/${requests.length}`);

  for (const request of requests) {
    await expect(page.locator("#aggregate-country")).toContainText(request.body.country);
    await expect(page.locator("#aggregate-continent")).toContainText(request.body.continent);
    await expect(page.locator("#aggregate-asn")).toContainText(request.body.asn_org);
  }

  await expect(page.getByText(`${100 - requests.length}% calls left`)).toBeVisible();

  await expect(page.locator("#bulk-results-download")).toHaveAttribute("href", /^blob:http.*/);
  const [download] = await Promise.all([page.waitForEvent("download"), page.locator("#bulk-results-download").click()]);
  await expect(download.suggestedFilename()).toMatch(/.*\.json$/);

  await page.locator("#bulk-clear").click();
  await expect(page.locator("#aggregate-country")).not.toBeVisible();
  await expect(page.locator("#aggregate-continent")).not.toBeVisible();
  await expect(page.locator("#aggregate-asn")).not.toBeVisible();

  await page.getByText("reset").click();
  await expect(input).toBeEmpty();
});
