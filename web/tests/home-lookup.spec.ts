import { expect, test } from "./setup";

test("validate home lookup IP", async ({ page, requests }) => {
  await page.goto("/");
  await page.waitForLoadState("domcontentloaded");

  const input = page.getByPlaceholder(/Search IP address/i);

  for (const request of requests) {
    await input.focus();
    await input.fill(request.body.query);
    await input.press("Enter");

    const result = page.locator("div.geo-result", { hasText: request.body.query });
    await expect(result).toBeVisible();

    await result.locator(`button:has-text("${request.body.query}")`).first().click();
    await expect(
      page.locator("[data-sonner-toast]").filter({
        hasText: new RegExp(`copied.*${request.body.query.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")}`),
      }),
    ).toBeVisible();

    await result.locator(".leaflet-control-zoom-out").click();

    await expect(result.locator("button img")).toHaveAttribute(
      "src",
      new RegExp(`.*/flags/${request.body.country_abbr.toLowerCase()}.svg`),
    );
  }

  await expect(page.getByTestId("rate-limit-badge")).toContainText("calls left");

  await page.getByRole("button", { name: /clear history/i }).click();
  await expect(page.locator("div.geo-result")).toHaveCount(0);
});
