import { createRouter, RouterProvider } from "@tanstack/react-router";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { client } from "@/api/client.gen";
import { applyRateLimitHeaders } from "@/lib/utils";
import { routeTree } from "./routeTree.gen";

import "./main.css";

client.setConfig({ baseUrl: "/api/v2" });
client.interceptors.response.use((response) => {
  applyRateLimitHeaders(response);
  return response;
});

const router = createRouter({
  routeTree,
  defaultPreload: "intent",
});

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

const rootEl = document.getElementById("root");
if (!rootEl) {
  throw new Error("missing #root element");
}
createRoot(rootEl).render(
  <StrictMode>
    <RouterProvider router={router} />
  </StrictMode>,
);
