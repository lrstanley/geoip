import { create } from "zustand";
import { persist } from "zustand/middleware";

import type { GeoResult } from "@/api/types.gen";

export type ClientState = {
  ratelimit_limit: number;
  ratelimit_remaining: number;
  ratelimit_reset: number;
};

type AppState = {
  clientState: ClientState;
  history: GeoResult[];
  setClientState: (c: ClientState) => void;
  addResult: (result: GeoResult) => void;
  clearHistory: () => void;
};

export const useAppStore = create<AppState>()(
  persist(
    (set, get) => ({
      clientState: {
        ratelimit_limit: 0,
        ratelimit_remaining: 0,
        ratelimit_reset: 0,
      },
      history: [],
      setClientState: (c) => set({ clientState: c }),
      addResult: (result) => {
        const { history } = get();
        const next = history.filter((h) => h.query !== result.query);
        next.push(result);
        set({ history: next });
      },
      clearHistory: () => set({ history: [] }),
    }),
    { name: "geoip-state-v5" },
  ),
);
