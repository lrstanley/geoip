import { createFileRoute } from "@tanstack/react-router";
import { Database, ExternalLink, HardDriveDownload, Search, Trash2 } from "lucide-react";
import { useEffect, useMemo, useState } from "react";
import { getManyAddresses } from "@/api/sdk.gen";
import type { BulkError, GeoResult } from "@/api/types.gen";
import { DefaultLayout } from "@/components/default-layout";
import { GeoAggregate } from "@/components/geo/geo-aggregate";
import { GeoMultiError } from "@/components/geo/geo-multi-error";
import { RateLimitCounter } from "@/components/rate-limit-counter";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Progress } from "@/components/ui/progress";
import { Separator } from "@/components/ui/separator";
import { Textarea } from "@/components/ui/textarea";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { BULK_SAMPLE_LINES } from "@/lib/bulk-sample-data";
import { cn, matchAddresses, useJsonObjectUrl } from "@/lib/utils";

export const Route = createFileRoute("/lookup/bulk")({
  component: BulkPage,
});

function BulkPage() {
  const [input, setInput] = useState("");
  const [results, setResults] = useState<GeoResult[]>([]);
  const [errors, setErrors] = useState<BulkError[]>([]);
  const [loading, setLoading] = useState(false);
  const [completed, setCompleted] = useState(0);
  const [total, setTotal] = useState(0);
  const [searchCount, setSearchCount] = useState(0);

  const percent = total > 0 ? Math.round((completed / total) * 100) : 0;

  const showProgress = loading || total > 0;

  const jsonPayload = useMemo(() => (results.length > 0 ? { data: results } : null), [results]);
  const jsonUrl = useJsonObjectUrl(jsonPayload);

  useEffect(() => {
    const t = setTimeout(() => {
      setSearchCount(matchAddresses(input).length);
    }, 500);
    return () => clearTimeout(t);
  }, [input]);

  useEffect(() => {
    document.title = "Bulk Lookup · GeoIP";
  }, []);

  async function search() {
    const addresses = matchAddresses(input);
    if (!input.trim() || addresses.length < 1) return;

    setLoading(true);
    setTotal(addresses.length);
    setCompleted(0);
    setErrors([]);
    const chunks: string[][] = [];
    const copy = [...addresses];
    while (copy.length) {
      chunks.push(copy.splice(0, 25));
    }

    for (const chunk of chunks) {
      try {
        const { data } = await getManyAddresses({
          body: {
            addresses: chunk,
            options: { disable_host_lookup: true },
          },
          throwOnError: true,
        });
        setResults((r) => [...r, ...data.results]);
        setErrors((e) => [...e, ...data.errors]);
      } catch (err) {
        setErrors((e) => [...e, { error: String(err), query: "-" }]);
      } finally {
        setCompleted((c) => c + chunk.length);
      }
    }

    setLoading(false);
  }

  function clearResults() {
    setCompleted(0);
    setTotal(0);
    setResults([]);
    setErrors([]);
  }

  return (
    <DefaultLayout>
      <div className="transition-all duration-500 ease-in-out">
        <div className="relative p-4">
          <div className="absolute top-6 right-6 z-10 flex gap-1">
            {searchCount > 0 && (
              <Badge variant="secondary" className="font-mono">
                {searchCount} addresses
              </Badge>
            )}
            {input.length > 0 && (
              <Badge
                variant="outline"
                className="cursor-pointer font-mono"
                onClick={() => {
                  setInput("");
                  setSearchCount(0);
                }}
              >
                reset
              </Badge>
            )}
          </div>
          <Textarea
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder="Bulk search IPs (8.8.8.8) or hosts (google.com)"
            rows={10}
            disabled={loading}
            className="scrollbar-compact max-h-[min(70vh,36rem)] min-h-[240px] resize-y overflow-y-auto font-mono text-sm"
          />
        </div>

        <div className="flex flex-row flex-wrap gap-2 px-4 pb-4">
          <RateLimitCounter allowSmall />

          <span className="ml-auto" />

          {input.length < 1 && (
            <Tooltip>
              <TooltipTrigger asChild>
                <Button
                  type="button"
                  variant="outline"
                  size="icon"
                  className="hidden md:flex"
                  onClick={() => setInput(BULK_SAMPLE_LINES.join("\n"))}
                >
                  <Database className="size-4" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>Use sample data</TooltipContent>
            </Tooltip>
          )}

          {(results.length > 0 || errors.length > 0) && !loading && (
            <Button id="bulk-clear" type="button" variant="secondary" size="sm" className="h-8" onClick={clearResults}>
              <Trash2 className="size-3.5" />
              clear
            </Button>
          )}

          {results.length > 0 && !loading && jsonUrl && (
            <>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button id="bulk-results-download" variant="outline" size="icon" className="hidden md:flex" asChild>
                    <a href={jsonUrl} download="geoip-results.json">
                      <HardDriveDownload className="size-4" />
                    </a>
                  </Button>
                </TooltipTrigger>
                <TooltipContent>Download results as JSON</TooltipContent>
              </Tooltip>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button id="bulk-results-open" variant="outline" size="icon" className="hidden md:flex" asChild>
                    <a href={jsonUrl} target="_blank" rel="noreferrer">
                      <ExternalLink className="size-4" />
                    </a>
                  </Button>
                </TooltipTrigger>
                <TooltipContent>Open JSON in new tab</TooltipContent>
              </Tooltip>
            </>
          )}

          <Button type="button" size="sm" className="h-8" disabled={loading} onClick={() => void search()}>
            <Search className="size-3.5" />
            search
          </Button>
        </div>

        {showProgress && (
          <>
            <Separator />
            <div id="bulk-progress" className="p-4">
              <div className="flex min-w-0 items-center gap-3">
                <Progress
                  value={percent}
                  className={cn("min-w-0 flex-1", loading ? "" : "**:data-[slot=progress-indicator]:bg-emerald-500")}
                />
                <p className="shrink-0 font-mono text-xs tabular-nums text-muted-foreground">
                  {completed}/{total}
                </p>
              </div>
            </div>
          </>
        )}

        {errors.length > 0 && !loading && (
          <>
            <Separator />
            <GeoMultiError value={errors} className="p-4" />
          </>
        )}

        {results.length > 0 && (
          <>
            <Separator />
            <GeoAggregate
              id="aggregate-country"
              className="p-4"
              value={results}
              truncate={{ percent: 5, label: "Other" }}
              field="country_abbr"
              label="country"
              flag
            />
            <Separator />
            <GeoAggregate
              id="aggregate-continent"
              className="p-4"
              value={results}
              field="continent_abbr"
              label="continent"
              truncate={{ percent: 5, label: "Other" }}
              variant="secondary"
            />
            <Separator />
            <GeoAggregate
              id="aggregate-asn"
              className="p-4"
              value={results}
              field="asn_org"
              label="asn_org"
              truncate={{ percent: 5, label: "Other" }}
              variant="secondary"
            />
          </>
        )}
      </div>
    </DefaultLayout>
  );
}
