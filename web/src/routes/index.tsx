import { createFileRoute } from "@tanstack/react-router";
import { BrushCleaning, Search } from "lucide-react";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { getAddress } from "@/api/sdk.gen";
import { DefaultLayout } from "@/components/default-layout";
import { GeoObject } from "@/components/geo/geo-object";
import { RateLimitCounter } from "@/components/rate-limit-counter";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { useAppStore } from "@/lib/store";
import { cn } from "@/lib/utils";

export const Route = createFileRoute("/")({
  validateSearch: (search: Record<string, unknown>) => ({
    q: typeof search.q === "string" ? search.q : undefined,
  }),
  component: HomePage,
});

function errorMessage(e: unknown): string {
  if (typeof e === "object" && e !== null && "error" in e) {
    const err = e as { error: unknown };
    if (typeof err.error === "string") return err.error;
  }
  if (e instanceof Error) return e.message;
  return String(e);
}

function HomePage() {
  const { q } = Route.useSearch();
  const navigate = Route.useNavigate();
  const history = useAppStore((s) => s.history);
  const addResult = useAppStore((s) => s.addResult);
  const clearHistory = useAppStore((s) => s.clearHistory);

  const [address, setAddress] = useState("");
  const [loading, setLoading] = useState(false);
  const [resultError, setResultError] = useState<string | null>(null);

  const reversed = useMemo(() => [...history].reverse(), [history]);

  const searchInputRef = useRef<HTMLInputElement>(null);

  const focusSearchInput = useCallback(() => {
    requestAnimationFrame(() => {
      searchInputRef.current?.focus();
    });
  }, []);

  const runLookup = useCallback(
    async (query: string) => {
      if (!query.trim()) return;
      setLoading(true);
      setResultError(null);
      void navigate({ search: { q: query } });
      try {
        const { data } = await getAddress({
          path: { address: query },
          throwOnError: true,
        });
        addResult(data);
        setAddress("");
      } catch (e) {
        setResultError(errorMessage(e));
      } finally {
        setLoading(false);
        focusSearchInput();
      }
    },
    [addResult, focusSearchInput, navigate],
  );

  const booted = useRef(false);
  useEffect(() => {
    if (booted.current) return;
    booted.current = true;
    if (q) {
      void runLookup(q);
    }
  }, [q, runLookup]);

  // Only run auto self-lookup once per mount. After clearHistory(), history is empty again;
  // without this guard the effect would re-fetch self and repopulate history (e2e flake on Firefox).
  const selfLookupRan = useRef(false);
  useEffect(() => {
    if (q || history.length > 0) return;
    if (selfLookupRan.current) return;
    selfLookupRan.current = true;
    void (async () => {
      try {
        const { data } = await getAddress({
          path: { address: "self" },
          throwOnError: true,
        });
        addResult(data);
      } catch {
        /* ignore */
      }
    })();
  }, [q, history.length, addResult]);

  useEffect(() => {
    document.title = "Lookup · GeoIP";
  }, []);

  return (
    <DefaultLayout>
      <div>
        <div className={`p-4 ${reversed.length < 1 ? "my-10" : ""}`}>
          <div className="relative">
            <Input
              ref={searchInputRef}
              type="text"
              value={address}
              autoFocus
              onChange={(e) => setAddress(e.target.value)}
              onBlur={() => setResultError(null)}
              onKeyDown={(e) => {
                if (e.key === "Enter") void runLookup(address);
              }}
              placeholder="Search IP address (e.g 1.2.3.4) or host (e.g google.com)"
              aria-label="Search IP address"
              className="h-11 pr-10 font-mono text-base"
              disabled={loading}
            />
            {!loading && (
              <Button
                type="button"
                variant="ghost"
                size="icon"
                className="absolute top-1/2 right-1 -translate-y-1/2"
                onClick={() => void runLookup(address)}
                aria-label="search"
              >
                <Search className="size-5" />
              </Button>
            )}
            {loading && (
              <span className="absolute top-1/2 right-3 -translate-y-1/2 text-xs text-muted-foreground">…</span>
            )}
          </div>

          {resultError && (
            <Alert variant="destructive" className="mt-2 py-2">
              <AlertDescription className="font-mono text-sm">error: {resultError}</AlertDescription>
            </Alert>
          )}
        </div>

        {reversed.length > 0 && <Separator />}

        {reversed.length > 0 && (
          <div className="px-4 pt-4 pb-2">
            <div className="mb-4 flex w-full min-w-0 flex-wrap items-center gap-2">
              <div className="min-w-0 shrink">
                <RateLimitCounter />
              </div>
              <Button
                type="button"
                variant="default"
                size="sm"
                className="ml-auto h-7 shrink-0 gap-1 text-xs"
                onClick={() => {
                  void navigate({ search: { q: undefined } });
                  clearHistory();
                }}
              >
                <BrushCleaning className="size-3.5" />
                clear history
              </Button>
            </div>

            <div className="divide-y divide-border/40">
              {reversed.map((item, index) => (
                <div key={item.query} className="list-group-item relative overflow-x-hidden contain-[layout]">
                  <div className={cn(index > 0 && "pt-px", index === reversed.length - 1 ? "pb-0" : "pb-px")}>
                    <GeoObject value={item} isLastInList={index === reversed.length - 1} />
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </DefaultLayout>
  );
}
