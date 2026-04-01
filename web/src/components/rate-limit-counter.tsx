import { Timer } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { useAppStore } from "@/lib/store";

export function RateLimitCounter({ allowSmall }: { allowSmall?: boolean }) {
  const clientState = useAppStore((s) => s.clientState);
  const remaining = clientState.ratelimit_remaining ?? 0;
  const limit = clientState.ratelimit_limit ?? 0;
  const percent = limit > 0 ? Math.floor((remaining / limit) * 100) : 0;

  const label = `${remaining.toLocaleString()} left of ${limit.toLocaleString()} limit`;

  if (!percent && !limit) {
    return null;
  }

  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <Badge
          variant="secondary"
          className="max-w-full min-w-0 justify-start font-mono text-xs"
          data-testid="rate-limit-badge"
        >
          <Timer className="size-3.5 shrink-0" aria-hidden />
          <span className="min-w-0 truncate">
            {percent}%<span className={allowSmall ? "max-md:hidden" : "hidden md:inline"}> calls left</span>
          </span>
        </Badge>
      </TooltipTrigger>
      <TooltipContent>{label}</TooltipContent>
    </Tooltip>
  );
}
