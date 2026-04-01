import type { GeoResult } from "@/api/types.gen";
import { GeoFlag } from "@/components/geo/geo-flag";
import { Progress } from "@/components/ui/progress";
import { calculateGroupWidth, groupByField, type TruncateOptions } from "@/lib/utils";

type GeoAggregateProps = {
  id?: string;
  value: GeoResult[];
  field: keyof GeoResult;
  label: keyof GeoResult;
  truncate?: TruncateOptions;
  flag?: boolean;
  variant?: "default" | "secondary";
  className?: string;
};

export function GeoAggregate({
  id,
  value,
  field,
  label,
  truncate,
  flag,
  variant = "default",
  className,
}: GeoAggregateProps) {
  const stats = groupByField(value, field, label, truncate);
  const widths = calculateGroupWidth(stats, 20);

  return (
    <div id={id} className={className}>
      {stats.map((stat) => (
        <div key={String(stat.field)} className="list-group-item relative overflow-x-hidden contain-[layout]">
          <div className="my-1 flex flex-auto items-center gap-2">
            <div
              className="shrink-0 truncate text-right font-mono text-xs"
              style={{ width: `${widths.labelWidth}ch` }}
              title={
                truncate && truncate.label === stat.label
                  ? `${truncate.label} represents aggregated groups with less than ${truncate.count ?? `${truncate.percent ?? 0}%`} total results each`
                  : String(stat.label)
              }
            >
              {String(stat.label)}
            </div>
            {flag && <GeoFlag value={String(stat.field)} size={16} className="inline-flex" immediate />}
            <Progress
              value={stat.percent}
              className={
                variant === "secondary" ? "h-2 flex-1 **:data-[slot=progress-indicator]:bg-emerald-500" : "h-2 flex-1"
              }
            />
            <div className="hidden shrink-0 font-mono text-xs md:block" style={{ width: `${widths.countWidth}ch` }}>
              {stat.count}
            </div>
            <div
              className="hidden shrink-0 font-mono text-xs md:block"
              style={{ width: `${widths.percentWidth + 3}ch` }}
            >
              ({stat.percent}%)
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
