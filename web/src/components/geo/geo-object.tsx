import { Clock, Globe, MapPin, Network, Search, Signpost } from "lucide-react";
import type { ReactNode } from "react";
import { toast } from "sonner";
import type { GeoResult } from "@/api/types.gen";
import { GeoFlag } from "@/components/geo/geo-flag";
import { GeoMap } from "@/components/geo/geo-map";
import { Button } from "@/components/ui/button";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";

function googleMapsUrl(lat: number, lng: number) {
  return `https://google.com/maps/place/${lat},${lng}/@${lat},${lng},5z/`;
}

async function copyText(text: string) {
  try {
    await navigator.clipboard.writeText(text);
    toast.success(`copied "${text}"`);
  } catch {
    toast.error("could not copy");
  }
}

const PILL_BASE =
  "min-w-0 shrink gap-1 overflow-hidden border font-mono text-[0.6875rem] leading-tight shadow-sm transition-colors [&_svg]:size-3";
const PRIMARY_COL = "flex shrink-0 flex-col justify-center sm:py-1 sm:pr-0";
const SECONDARY_CLUSTER = "flex min-w-0 flex-1 flex-wrap items-center gap-1.5 [column-gap:0.5rem] [row-gap:0.375rem]";
const SECONDARY_CLUSTER_ASN_HOST =
  "flex min-w-0 w-full flex-wrap items-center justify-end gap-1 [column-gap:0.375rem] [row-gap:0.25rem] sm:flex-1";
const ASN_HOST_COMPACT = "flex-[0_1_auto] max-w-[min(100%,10rem)] text-[0.625rem] [&_svg]:size-2.5";
const META_SHRINK = "flex-[0_1_auto] max-w-[min(100%,10rem)]";
const SUMMARY_SHRINK = "flex-[0_1_auto] max-w-[min(100%,12rem)]";

type InfoPillProps = {
  tooltip: string;
  icon: ReactNode;
  children: ReactNode;
  colorClass: string;
  className?: string;
  copy?: string;
  href?: string;
};

function InfoPill({ tooltip, icon, children, colorClass, className, copy, href }: InfoPillProps) {
  const content = (
    <>
      {icon}
      <span className="min-w-0 truncate">{children}</span>
    </>
  );

  const btnClass = cn(PILL_BASE, colorClass, className);

  return (
    <Tooltip>
      <TooltipTrigger asChild>
        {href ? (
          <Button variant="outline" size="xs" className={btnClass} asChild>
            <a href={href} target="_blank" rel="noreferrer" className="flex min-w-0 items-center gap-1">
              {content}
            </a>
          </Button>
        ) : (
          <Button variant="outline" size="xs" className={btnClass} onClick={copy ? () => copyText(copy) : undefined}>
            {content}
          </Button>
        )}
      </TooltipTrigger>
      <TooltipContent>{tooltip}</TooltipContent>
    </Tooltip>
  );
}

type GeoObjectProps = {
  value: GeoResult;
  isLastInList?: boolean;
};

export function GeoObject({ value: geo, isLastInList = false }: GeoObjectProps) {
  return (
    <div className="geo-result" data-testid="geo-result">
      <div className="flex flex-col gap-2 pt-2 pb-1.5 sm:flex-row sm:items-stretch sm:gap-x-4 sm:pb-0">
        <div className={PRIMARY_COL}>
          <InfoPill
            tooltip="search query"
            icon={<Search className="shrink-0" aria-hidden />}
            copy={geo.query}
            colorClass="max-w-full border-emerald-500/45 bg-emerald-500/10 text-emerald-300 hover:bg-emerald-500/20 hover:text-emerald-200 sm:max-w-[min(100%,13rem)]"
          >
            {geo.query}
          </InfoPill>
        </div>

        <div className={SECONDARY_CLUSTER_ASN_HOST}>
          {geo.asn_org && (
            <InfoPill
              tooltip={`Autonomous System Organization: ${geo.network}`}
              icon={<Network className="shrink-0" aria-hidden />}
              href={`https://bgp.he.net/${geo.asn}#_whois`}
              colorClass="hidden max-w-full border-sky-500/45 bg-sky-500/10 text-sky-200 hover:bg-sky-500/20 hover:text-sky-100 md:inline-flex"
              className={ASN_HOST_COMPACT}
            >
              {geo.asn_org}
            </InfoPill>
          )}

          {geo.host && (
            <InfoPill
              tooltip="reverse dns"
              icon={<Globe className="shrink-0" aria-hidden />}
              copy={geo.host}
              colorClass="max-w-full border-violet-500/45 bg-violet-500/10 text-violet-200 hover:bg-violet-500/20 hover:text-violet-100"
              className={ASN_HOST_COMPACT}
            >
              {geo.host}
            </InfoPill>
          )}
        </div>
      </div>

      <GeoMap value={geo} />

      <div
        className={cn(
          geo.timezone
            ? "flex flex-col gap-2 pt-1.5 sm:pt-0 md:flex-row md:items-stretch md:gap-x-4"
            : "flex w-full min-w-0 flex-wrap items-center gap-1.5 pt-1.5 sm:pt-0",
          isLastInList ? "pb-0" : "pb-2",
        )}
      >
        {geo.timezone && (
          <div className={cn(PRIMARY_COL, "hidden md:flex")}>
            <InfoPill
              tooltip="timezone"
              icon={<Clock className="shrink-0" aria-hidden />}
              copy={geo.timezone}
              colorClass="max-w-full border-zinc-500/55 bg-zinc-800/90 text-zinc-200 hover:bg-zinc-700/90 hover:text-zinc-50 md:max-w-[min(100%,16rem)]"
            >
              {geo.timezone}
            </InfoPill>
          </div>
        )}

        <div className={cn(SECONDARY_CLUSTER, geo.timezone && "md:flex-1 md:justify-end")}>
          {geo.postal_code && (
            <InfoPill
              tooltip="postal code"
              icon={<Signpost className="shrink-0" aria-hidden />}
              copy={geo.postal_code}
              colorClass="border-amber-500/45 bg-amber-500/10 text-amber-200 hover:bg-amber-500/20 hover:text-amber-50"
              className={META_SHRINK}
            >
              {geo.postal_code}
            </InfoPill>
          )}

          {geo.summary && (
            <InfoPill
              tooltip="location"
              icon={<GeoFlag value={geo.country_abbr} size={14} className="shrink-0" />}
              copy={geo.summary}
              colorClass="border-teal-500/45 bg-teal-500/10 text-teal-200 hover:bg-teal-500/20 hover:text-teal-50"
              className={SUMMARY_SHRINK}
            >
              {geo.summary}
            </InfoPill>
          )}

          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                variant="outline"
                size="icon-xs"
                className="shrink-0 border-red-500/45 bg-red-500/10 text-red-400 hover:bg-red-500/20 hover:text-red-300 [&_svg]:size-3.5"
                asChild
              >
                <a href={googleMapsUrl(geo.latitude, geo.longitude)} target="_blank" rel="noreferrer">
                  <MapPin aria-hidden />
                </a>
              </Button>
            </TooltipTrigger>
            <TooltipContent>open in Google Maps</TooltipContent>
          </Tooltip>
        </div>
      </div>
    </div>
  );
}
