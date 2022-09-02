import ipRegex from "ip-regex"

const reHostname = /\b((?=[a-z0-9-]{1,63}\.)(xn--)?[a-z0-9]+(-[a-z0-9]+)*\.)+[a-z]{2,63}\b/gm

export function matchAddresses(input: string): string[] {
  const ips = input.match(ipRegex())
  const hostnames = input.match(reHostname)

  return [...new Set([...(ips ?? []), ...(hostnames ?? [])])]
}
