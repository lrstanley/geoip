## :gear: Usage

#### Application Options
| Environment vars | Flags | Type | Description |
| --- | --- | --- | --- |
| - | `-v, --version` | bool | prints version information and exits |
| - | `--version-json` | bool | prints version information in JSON format and exits |
| `DEBUG` | `-D, --debug` | bool | enables debug mode |

#### HTTP Server options
| Environment vars | Flags | Type | Description |
| --- | --- | --- | --- |
| `HTTP_BIND_ADDR` | `--http.bind-addr` | string | ip:port pair to bind to [**required**] [**default: :8080**] |
| `HTTP_TRUSTED_PROXIES` | `--http.trusted-proxies` | []string | CIDR ranges that we trust the X-Forwarded-For header from (addl opts: local, *, cloudflare, and/or custom header to use) |
| `HTTP_MAX_CONCURRENT` | `--http.max-concurrent` | int | limit total max concurrent requests across all connections (0 for no limit) |
| `HTTP_LIMIT` | `--http.limit` | int | number of requests/ip/hour [**default: 2000**] |
| `HTTP_HSTS` | `--http.hsts` | bool | enable HTTP Strict Transport Security |
| `HTTP_CORS` | `--http.cors` | []string | CORS allowed origins [**default: ***] |
| `HTTP_METRICS` | `--http.metrics` | bool | enable prometheus metrics on /metrics to internal IPs |

#### DB Options
| Environment vars | Flags | Type | Description |
| --- | --- | --- | --- |
| `DB_GEOIP_PATH` | `--db.geoip-path` | string | path to read/store GeoIP Maxmind DB [**default: geoip.db**] |
| `DB_GEOIP_UPDATE_URL` | `--db.geoip-update-url` | string | GeoIP database file download location (must be gzipped) [**default: https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=%s&suffix=tar.gz**] |
| `DB_ASN_PATH` | `--db.asn-path` | string | path to read/store ASN Maxmind DB [**default: asn.db**] |
| `DB_ASN_UPDATE_URL` | `--db.asn-update-url` | string | ASN database file download location (must be gzipped) [**default: https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-ASN&license_key=%s&suffix=tar.gz**] |
| `DB_LICENSE_KEY` | `--db.license-key` | string | maxmind license key (must register for a maxmind account) [**required**] |
| `DB_UPDATE_INTERVAL` | `--db.update-interval` | time.Duration | interval of time between database update checks [**default: 12h**] |
| `DB_CACHE_SIZE` | `--db.size` | int | total number of lookups to keep in ARC cache (50% most recent, 50% most requested) [**default: 1000**] |
| `DB_CACHE_EXPIRE` | `--db.expire` | time.Duration | expiration time of cache [**default: 1h**] |
| `DB_DEFAULT_LANGUAGE` | `--db.lang` | string | default language to use for geolocation [**default: en**] |

#### DNS Lookup Options
| Environment vars | Flags | Type | Description |
| --- | --- | --- | --- |
| `DNS_RESOLVERS` | `--dns.resolver` | []string | resolver (in host:port form) to use for dns lookups (doesn't work with windows and plan9) (can be used multiple times) |
| `DNS_LOCAL` | `--dns.uselocal` | bool | adds local (system) resolvers to the list of resolvers to use |
| `DNS_CACHE_SIZE` | `--dns.size` | int | total number of lookups to keep in ARC cache (50% most recent, 50% most requested) [**default: 500**] |
| `DNS_CACHE_EXPIRE` | `--dns.expire` | time.Duration | expiration time of cache [**default: 1h**] |
| `DNS_TIMEOUT` | `--dns.timeout` | time.Duration | timeout for dns lookups (longer = better results but longer request duration) [**default: 4s**] |

#### Logging Options
| Environment vars | Flags | Type | Description |
| --- | --- | --- | --- |
| `LOG_QUIET` | `--log.quiet` | bool | disable logging to stdout (also: see levels) |
| `LOG_LEVEL` | `--log.level` | string | logging level [**default: info**] [**choices: debug, info, warn, error, fatal**] |
| `LOG_JSON` | `--log.json` | bool | output logs in JSON format |
| `LOG_PRETTY` | `--log.pretty` | bool | output logs in a pretty colored format (cannot be easily parsed) |
| `LOG_PATH` | `--log.path` | string | path to log file (disables stdout logging) |
