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

export const Route = createFileRoute("/lookup/docs")({
  component: DocsPage,
});

function DocsPage() {
  const origin = typeof window !== "undefined" ? window.location.origin : "";

  useEffect(() => {
    document.title = "Documentation · GeoIP";
  }, []);

  const doc = useMemo(() => {
    const props: Record<string, string> = {
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
      "allow-api-list-style-selection": "false",
      "load-fonts": "false",
      "regular-font": "Consolas, monaco, monospace",
      "mono-font": "Consolas, monaco, monospace",
      "font-size": "large",
      "server-url": `${origin}/api/v2`,
      "default-api-server": `${origin}/api/v2`,
      "nav-bg-color": THEME.navBg,
      "nav-hover-bg-color": THEME.navHover,
      "bg-color": THEME.bg,
      "text-color": THEME.text,
      "primary-color": THEME.primary,
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
