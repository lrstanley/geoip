import { type ClassValue, clsx } from "clsx";
import ipRegex from "ip-regex";
import { useEffect, useMemo } from "react";
import { twMerge } from "tailwind-merge";
import type { GeoResult } from "@/api/types.gen";
import { useAppStore } from "@/lib/store";

const reHostname = /\b((?=[a-z0-9-]{1,63}\.)(xn--)?[a-z0-9]+(-[a-z0-9]+)*\.)+[a-z]{2,63}\b/gim;

export function matchAddresses(input: string): string[] {
  const ips = input.match(ipRegex());
  const hostnames = input.match(reHostname);
  return [...new Set([...(ips ?? []), ...(hostnames ?? [])])];
}

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export interface TruncateOptions {
  label: string;
  percent?: number;
  count?: number;
}

export interface ByField<Field, Label> {
  count: number;
  percent: number;
  field: Field;
  label: Label;
}

export type FieldGroup = ByField<GeoResult[keyof GeoResult], GeoResult[keyof GeoResult]>;

export interface GroupWidth {
  labelWidth: number;
  percentWidth: number;
  countWidth: number;
}

export function groupByField(
  results: GeoResult[],
  field: keyof GeoResult,
  label: keyof GeoResult,
  truncate?: TruncateOptions,
): FieldGroup[] {
  const calc: Record<string, FieldGroup> = {};

  for (const result of results) {
    const fieldVal = result[field];
    const labelVal = result[label];
    const key = String(fieldVal);

    if (!calc[key]) {
      calc[key] = {
        count: 1,
        percent: 0,
        field: fieldVal,
        label: labelVal,
      };
    } else {
      calc[key].count++;
    }
  }

  for (const key of Object.keys(calc)) {
    calc[key].percent = Math.round((calc[key].count / results.length) * 100);
  }

  const val = Object.values(calc).sort((a, b) => b.count - a.count);

  if (truncate) {
    const other: FieldGroup = {
      count: 0,
      percent: 0,
      field: truncate.label,
      label: truncate.label,
    };

    for (let i = val.length - 1; i >= 0; i--) {
      if (val[i].label === truncate.label) {
        continue;
      }

      if (truncate.percent && val[i].percent < truncate.percent) {
        other.count += val[i].count;
        val.splice(i, 1);
      } else if (truncate.count && val[i].count < truncate.count) {
        other.count += val[i].count;
        val.splice(i, 1);
      }
    }

    if (other.count > 0) {
      other.percent = Math.round((other.count / results.length) * 100);
      val.push(other);
    }
  }

  return val;
}

export function calculateGroupWidth(stats: FieldGroup[], maxLabelWidth = 30): GroupWidth {
  const label = stats.reduce((max, stat) => Math.max(max, String(stat.label).length), 0);
  const count = stats.reduce((max, stat) => Math.max(max, String(stat.count).length), 0);
  const percent = stats.reduce((max, stat) => Math.max(max, String(stat.percent).length), 0);

  return {
    labelWidth: Math.min(label, maxLabelWidth),
    countWidth: count,
    percentWidth: percent,
  };
}

export function useJsonObjectUrl(data: object | null): string | undefined {
  const url = useMemo(() => {
    if (!data) return undefined;
    const blob = new Blob([JSON.stringify(data, null, 4)], {
      type: "application/json",
    });
    return URL.createObjectURL(blob);
  }, [data]);

  useEffect(() => {
    if (!url) return;
    return () => {
      URL.revokeObjectURL(url);
    };
  }, [url]);

  return url;
}

export function applyRateLimitHeaders(res: Response): void {
  const limit = parseInt(res.headers.get("X-Ratelimit-Limit") ?? "0", 10);
  const remaining = parseInt(res.headers.get("X-Ratelimit-Remaining") ?? "0", 10);
  const reset = parseInt(res.headers.get("X-Ratelimit-Reset") ?? "0", 10);
  if (limit > 0 || remaining > 0 || reset > 0) {
    useAppStore.getState().setClientState({
      ratelimit_limit: limit,
      ratelimit_remaining: remaining,
      ratelimit_reset: reset,
    });
  }
}
