import { useState } from "@/lib/core/state"

import type { ApiError, GeoResult, Error } from "@/lib/api/openapi"
import type { ClientState } from "@/lib/api/types"

export function saveResult(result: GeoResult) {
  const state = useState()

  for (let i = 0; i < state.history.length; i++) {
    if (state.history[i].query == result.query) {
      state.history.splice(i, 1)
      break
    }
  }
  state.history.push(result)
}

/**
 * getClientState is a wrapper around an API response to parse the HTTP response
 * headers and return information about rate limiting and database information.
 *
 * @param {globalThis.Response} resp
 * @returns {ClientState}
 */
export function getClientState(resp: globalThis.Response): ClientState {
  return {
    ratelimit_limit: parseInt(resp.headers.get("X-Ratelimit-Limit") || "0"),
    ratelimit_remaining: parseInt(resp.headers.get("X-Ratelimit-Remaining") || "0"),
    ratelimit_reset: parseInt(resp.headers.get("X-Ratelimit-Reset") || "0"),
  }
}

/**
 * getError is a wrapper around an API response to convert it into a GeoIP Error
 * object.
 *
 * @export
 * @param {ApiError} e
 * @returns {Error}
 */
export function getError(e: ApiError): Error {
  return e.body
}
