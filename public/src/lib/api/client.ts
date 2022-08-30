import { getClientState } from "@/lib/api/helpers"
import { HTTPClient } from "@/lib/api/openapi"
import { useState } from "@/lib/core/state"

// Monkey patch the fetch function to add the client state to global state.
const { fetch: originalFetch } = window
window.fetch = async (...args) => {
  const [resource, config] = args

  const resp = await originalFetch(resource, config)

  const state = useState()
  state.clientState = getClientState(resp)

  return resp
}

// Define the API client.
export const api = new HTTPClient({
  HEADERS: {
    "Content-Type": "application/json",
  },
  BASE: "/api/v2",
})
