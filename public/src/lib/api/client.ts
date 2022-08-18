import pLimit from "p-limit"
import { DefaultConfig } from "@/lib/api/config"

import type { ClientState, APIResponse } from "@/lib/api/types"

/**
 * getClientState is a wrapper around an API response to parse the HTTP response
 * headers and return information about rate limiting and database information.
 *
 * @param {globalThis.Response} resp
 * @returns {ClientState}
 */
function getClientState(resp: globalThis.Response): ClientState {
  return {
    cached: resp.headers.get("X-Cache") == "HIT",
    maxmind_build: resp.headers.get("X-Maxmind-Build"),
    maxmind_type: resp.headers.get("X-Maxmind-Type"),
    ratelimit_limit: parseInt(resp.headers.get("X-Ratelimit-Limit") || "0"),
    ratelimit_remaining: parseInt(resp.headers.get("X-Ratelimit-Remaining") || "0"),
    ratelimit_reset: parseInt(resp.headers.get("X-Ratelimit-Reset") || "0"),
  }
}

/**
 * lookup queries the GeoIP API for a given IP address or hostname. If the lookup
 * is successful, the response is returned. If the lookup fails, an error is
 * returned via the error property of the APIResponse type. Use the save argument
 * to persist the lookup to localstorage state, which shows on the homepage.
 *
 * @export
 * @async
 * @param {string} addr
 * @param {boolean} save
 * @returns {Promise<APIResponse>}
 */
export async function lookup(addr: string, save: boolean): Promise<APIResponse> {
  const contoller = new AbortController()
  const options: RequestInit = {
    ...DefaultConfig,
    method: "GET",
    signal: contoller.signal,
  }

  const timeout = setTimeout(() => contoller.abort(), 15000)
  const resp = await fetch(`${DefaultConfig.baseUrl}/${addr}`, options)
  clearTimeout(timeout)

  let ok = resp.ok
  const data = await resp.json()

  if (resp.status != 200 || data.error) {
    ok = false
  }

  const result: APIResponse = {
    query: addr,
    state: getClientState(resp),
    data: ok ? data : null,
    error: ok ? null : new Error(data.error || resp.statusText || "Unknown error"),
  }

  const state = useState()
  state.clientState = result.state

  if (save && !result.error) {
    for (let i = 0; i < state.history.length; i++) {
      if (state.history[i].query == result.query) {
        state.history.splice(i, 1)
        break
      }
    }
    state.history.push(result)
  }

  return Promise.resolve(result)
}

/**
 * limit is the maximum number of concurrent requests to make at a time. This
 * number may not reflect the exact concurrency, as it's checked at an interval,
 * and there are also limits by the browser. The below number may be higher than
 * the browser's limit, but it should be low enough that the browser doesn't
 * exhaust resources queueing up all requests at once.
 */
const limit = pLimit(30)

/**
 * lookupMany allows looking up multiple addresses at once, in a bulk operation.
 * The concurrency is limited by the pLimit library (see definition above). Use
 * cb as a callback to get the results, which is called once per address result as
 * the results are received.
 *
 * @export
 * @param {string[]} addresses
 * @param {boolean} save
 * @param {(result: APIResponse) => void} cb
 * @returns {Promise<void[]>}
 */
export function lookupMany(
  addresses: string[],
  save: boolean,
  cb: (result: APIResponse) => void
): Promise<void[]> {
  const promises: Promise<void>[] = []

  for (const addr of addresses) {
    promises.push(
      limit(() => {
        return lookup(addr, save)
          .then(cb)
          .catch((e) => {
            cb({ query: addr, state: null, data: null, error: e })
          })
      })
    )
  }

  return Promise.all(promises)
}
