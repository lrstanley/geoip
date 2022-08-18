import type { Config } from "@/lib/api/types"

export const DefaultConfig: Config = {
  headers: new Headers({
    Accept: "application/json",
  }),
  baseUrl: "/api",
}
