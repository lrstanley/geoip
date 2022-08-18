<route>
meta:
  title: Documentation
</route>
<script setup>
  import { useState } from "@/lib/core/state"
  const state = useState()
</script>

<DocsHeader title="API Documentation" subtitle="Query the GeoIP API through REST and JSON" el="h1" icon nopadding>
  <template #icon><i-mdi-cogs /></template>
</DocsHeader>

The HTTP API takes GET requests in the following scheme:

```bash
http[s]://geoip.pw/api{/endpoint}
```

<DocsHeader
  title="Lookup an IP address or hostname"
  subtitle="How do I obtain the geographical info for a hostname or IP?"
  el="h2"
  id="lookup"
/>

The GeoIP API supports both IP (IPv4 and IPv6) lookups, as well as hostname & domain
lookups. To obtain JSON information about a given IP address or hostname, simply use
the following format (replacing `{query}` with the IP or hostname):

```bash
https://geoip.pw/api/{ip}
```

<br />
For example:

```json
$ curl "https://geoip.pw/api/8.8.8.8"
{"ip":"8.8.8.8","summary":"United States, NA","city":"","subdivision":"","country":"United States","country_abbr":"US","continent":"North America","continent_abbr":"NA","latitude":37.751,"longitude":-97.822,"timezone":"America/Chicago","postal_code":"","proxy":false,"host":"dns.google"}
```

<br />
We can take that one step further, and prettify the JSON:

```json
$ curl "https://geoip.pw/api/8.8.8.8?pretty=1"
{
  "ip": "8.8.8.8",
  "summary": "United States, NA",
  "city": "",
  "subdivision": "",
  "country": "United States",
  "country_abbr": "US",
  "continent": "North America",
  "continent_abbr": "NA",
  "latitude": 37.751,
  "longitude": -97.822,
  "timezone": "America/Chicago",
  "postal_code": "",
  "proxy": false,
  "host": "dns.google"
}
```

<br />
GeoIP also supports providing a `lang` query parameter (or `Accept-Language` header)
to specify the language of the response:

```json
$ curl "https://geoip.pw/api/8.8.8.8?pretty=1&lang=zh-CN"
{
  "ip": "8.8.8.8",
  "summary": "美国, NA",
  "city": "",
  "subdivision": "",
  "country": "美国",
  "country_abbr": "US",
  "continent": "北美洲",
  "continent_abbr": "NA",
  "latitude": 37.751,
  "longitude": -97.822,
  "timezone": "America/Chicago",
  "postal_code": "",
  "proxy": false,
  "host": "dns.google"
}
```

<DocsHeader
  title="Filtering output"
  subtitle="What if I only want a specific field to reduce query time?"
  el="h2"
  id="filtering-output"
/>

The query API supports filtering and obtaining only specific fields. The supported
fields are what you currently see in the JSON output above. The full lookup endpoint
does things like reverse DNS lookups, so using this, will allow you to speed up
potential queries if not all the information is being used.

<br />
The syntax is as follows:

```bash
http[s]://geoip.pw/api/{query}/{comma-list-of-fields}
```

Warning: the `{comma-list-of-fields}` must be a set of valid field parameters, or
the API will return an error. Errors returned by this endpoint will be plaintext
(e.g. `error: too many filters supplied`), which should hopefully cater to those
using this endpoint for shell scripts.

The response to this query will be a pipe (`|`) separated array of values. Pipe
is used primarily as it's very unlikely returned fields themselves will contain
pipe symbols (compared to commas, for example), which allows things like shell
scripts to easily parse the resulting output.

<br />
For example:

```bash
$ curl "https://geoip.pw/api/8.8.8.8/summary,country,host"
United States, NA|United States|google-public-dns-a.google.com
```

<DocsHeader
  title="Additional headers"
  subtitle="How much of a rate limit do I have? How old is the GeoIP data?"
  el="h2"
  id="addl-headers"
/>

With each API response, the API will tag a few headers that will help you better
determine what your current rate limit is. If you want to obtain these headers
without submitting a query for a given IP/hostname, simply use the below endpoint:

```bash
http[s]://geoip.pw/api/ping
```

#### Access-Control headers

GeoIP also has support to restrict connections based on the `Access-Control-Allow-Origin`
header, as long as end browsers/clients support it. By default, this isn't enabled,
however self-hosted versions may utilize this.

<br />
The following are the currently supported headers:

| **Header**              | **Always Present** | **Description**                                                     |
| ----------------------- | ------------------ |-------------------------------------------------------------------- |
| `X-Maxmind-Build`       | <DocsCheck />      | The [Maxmind][maxmind] database version                             |
| `X-Maxmind-Type`        | <DocsCheck />      | The [Maxmind][maxmind] db type -- see [this page][db] for more info |
| `X-Ratelimit-Limit`     | -                  | Number of queries allowed within `reset` period                     |
| `X-Ratelimit-Remaining` | -                  | Number of queries remaining within `limit`                          |
| `X-Ratelimit-Reset`     | -                  | Seconds until `remaining` is reset                                  |

<br />
The current values of these headers:

| **Header**              | **Value**                                   |
| ----------------------- | ------------------------------------------- |
| `X-Maxmind-Build`       | {{ state.clientState.maxmind_build }}       |
| `X-Maxmind-Type`        | {{ state.clientState.maxmind_type }}        |
| `X-Ratelimit-Limit`     | {{ state.clientState.ratelimit_limit }}     |
| `X-Ratelimit-Remaining` | {{ state.clientState.ratelimit_remaining }} |
| `X-Ratelimit-Reset`     | {{ state.clientState.ratelimit_reset }}     |

<n-alert type="warning" :show-icon="false">
Warning: not all self-hosted versions of GeoIP may have rate-limiting enabled, and
as such, the header may not always be present. (see the "supported headers" table,
showing which headers always will be present.)
</n-alert>

<br />
An example response may look like this:

```bash
$ curl -i "https://geoip.pw/api/ping"
HTTP/1.1 200 OK
Content-Type: application/json
X-Maxmind-Build: 6-1501721259
X-Maxmind-Type: GeoLite2-City
X-Ratelimit-Limit: 2000
X-Ratelimit-Remaining: 1993
X-Ratelimit-Reset: 3593

{"pong":true}
```

<n-alert type="info" :show-icon="false">
Note: the <code>/api/ping</code> endpoint is the only one which currently supports HTTP HEAD
requests (which still contain the necessary headers.)
</n-alert>

[maxmind]: http://www.maxmind.com/
[db]: https://dev.maxmind.com/geoip/geoip2/geolite2/
