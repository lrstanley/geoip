import { createFileRoute } from "@tanstack/react-router";
import { createElement, useEffect, useMemo } from "react";
import { Navigation } from "@/components/navigation";

import "rapidoc";

const THEME = {
  navBg: "#0F0F0F",
  navHover: "#1B1B1B",
  bg: "#18181C",
  text: "#BBBBBB",
  primary: "#10B981",
};

const FONT = "Consolas, monaco, monospace";

export const Route = createFileRoute("/lookup/docs")({
  component: DocsPage,
});

function DocsPage() {
  const origin = typeof window !== "undefined" ? window.location.origin : "";

  useEffect(() => {
    document.title = "Documentation · GeoIP";
  }, []);

  const doc = useMemo(() => {
    const props: Record<string, unknown> = {
      "render-style": "focused",
      "spec-url": "/api/v2/openapi.yaml",
      "show-header": "false",
      "show-info": "false",
      "show-method-in-nav-bar": "as-colored-text",
      "nav-item-spacing": "relaxed",
      "allow-authentication": "false",
      "allow-server-selection": "false",
      "allow-spec-url-load": "false",
      "allow-spec-file-load": "false",
      "allow-spec-file-download": "false",
      "color-scheme": "dark",
      scale: "lg",
      "server-url": `${origin}/api/v2`,
      "default-api-server": `${origin}/api/v2`,
      "primary-color": THEME.primary,
      // RapiDoc v10+ themes through CSS custom properties on the host element,
      // rather than individual color/font attributes.
      style: {
        "--font-sans": FONT,
        "--font-mono": FONT,
        "--background": THEME.bg,
        "--foreground": THEME.text,
        "--sidebar": THEME.navBg,
        "--sidebar-accent": THEME.navHover,
        // Don't clip the slotted navigation, which is taller than a plain logo.
        "--sidebar-logo-max-height": "none",
      },
      className: "h-full w-full min-h-0 flex-1",
    };
    return createElement(
      "rapi-doc",
      props,
      createElement(
        "div",
        { slot: "nav-logo", className: "pb-2" },
        createElement(Navigation, { slim: true, hideSource: true }),
      ),
    );
  }, [origin]);

  return <div className="flex h-dvh flex-col overflow-hidden bg-background">{doc}</div>;
}
