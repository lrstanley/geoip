import type { ReactNode } from "react";
import { Navigation } from "@/components/navigation";
import { Card } from "@/components/ui/card";

type DefaultLayoutProps = {
  children: ReactNode;
};

export function DefaultLayout({ children }: DefaultLayoutProps) {
  return (
    <div className="mx-auto mt-3 grid max-h-full w-full max-w-4xl shrink grow-0 basis-auto grid-cols-1 gap-2 lg:mt-32 lg:grid-cols-[245px_minmax(0,1fr)]">
      <Navigation />

      <div className="mx-3 flex flex-col">
        <Card className="rounded-md border bg-card py-0 shadow-xl drop-shadow-xl">{children}</Card>

        <div className="mb-3 px-2 py-4 text-right text-sm text-muted-foreground lg:mb-20">
          GeoLite data from{" "}
          <a className="text-primary underline" href="http://www.maxmind.com" target="_blank" rel="noreferrer">
            MaxMind
          </a>
          {" · "}
          GeoIP:{" "}
          <a
            className="text-primary underline"
            href="https://github.com/lrstanley/geoip"
            target="_blank"
            rel="noreferrer"
          >
            FOSS
          </a>{" "}
          lookup service, made with{" "}
          <span className="text-red-500" aria-hidden>
            ♥
          </span>
        </div>
      </div>
    </div>
  );
}
