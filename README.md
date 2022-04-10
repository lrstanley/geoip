<!-- template:begin:header -->
<!-- template:end:header -->

<!-- template:begin:toc -->
<!-- template:end:toc -->

## :computer: Installation

Check out the [releases](https://github.com/lrstanley/geoip/releases)
page for prebuilt versions.

### :whale: Container Images (ghcr)

```console
$ docker run -it --rm -p 8080:80 ghcr.io/lrstanley/geoip:latest geoip --http.bind 0.0.0.0:80 --db /data/geoip.db
$ curl -I http://localhost:8080
HTTP/1.1 200 OK
Content-Type: text/html
Date: Thu, 06 Aug 2020 00:55:21 GMT
```

### :toolbox: Build From Source

Dependencies (to build from source only):

   * [Go](https://golang.org/doc/install) (latest)
   * [NodeJS](https://nodejs.org/en/download/) (v8)

Setup:

```console
$ git clone <repo>
$ cd geoip
# this will show you all of the available options (to fetch dependencies, run in debug mode, etc.)
$ make help
$ make
$ ./geoip --help
```

For active development:

```console
$ make debug
# run this in a different window. this will rebundle the frontend assets on
# change.
$ make frontend-watch
```

## :gear: Usage

```console
$ geoip --help
Usage:
  geoip [OPTIONS]

Application Options:
  -d, --debug          enable exception display and pprof endpoints (warn: dangerous) [$DEBUG]
  -q, --quiet          disable verbose output [$QUIET]
      --db=            path to read/store Maxmind DB (default: geoip.db) [$DB_PATH]
      --interval=      interval of time between database update checks (default: 12h) [$UPDATE_INTERVAL]
      --update-url=    maxmind database file download location (must be gzipped) (default: https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=%s&suffix=tar.gz)
                       [$MAXMIND_UPDATE_URL]
      --license-key=   maxmind license key (must register for a maxmind account) [$MAXMIND_LICENSE_KEY]
  -v, --version        print the version and compilation date

Cache Options:
      --cache.size=    total number of lookups to keep in ARC cache (50% most recent, 50% most requested) (default: 500) [$CACHE_SIZE]
      --cache.expire=  expiration time of cache (default: 20m) [$CACHE_EXPIRE]

HTTP Options:
  -b, --http.bind=     address and port to bind to (default: :8080) [$HTTP_BIND]
      --http.proxy     obey X-Forwarded-For headers (warn: dangerous, make sure to only bind to localhost) [$HTTP_BEHIND_PROXY]
      --http.throttle= limit total max concurrent requests across all connections [$HTTP_THROTTLE]
      --http.limit=    number of requests/ip/hour (default: 2000) [$HTTP_LIMIT]
      --http.cors=     cors origin domain to allow with https?:// prefix (empty => '*'; use flag multiple times) [$HTTP_CORS]

TLS Options:
      --http.tls.use   enable tls [$TLS_USE]
      --http.tls.cert= path to ssl certificate [$TLS_CERT]
      --http.tls.key=  path to ssl key [$TLS_KEY]

DNS Lookup Options:
      --dns.timeout=   max allowed duration when looking up hostnames (may cause queries to be slow) (default: 2s) [$DNS_TIMEOUT]
      --dns.resolver=  resolver (in host:port form) to use for dns lookups (doesn't work with windows and plan9) (can be used multiple times) [$DNS_RESOLVERS]
      --dns.uselocal   adds local (system) resolvers to the list of resolvers to use [$DNS_LOCAL]

Help Options:
  -h, --help           Show this help message

```

### Example

```console
$ geoip --cache.size 1000 --http.bind "localhost:8080" --http.proxy --http.limit 15000 --dns.resolver 8.8.8.8 --dns.resolver 8.8.4.4
```

<!-- template:begin:support -->
<!-- template:end:support -->

<!-- template:begin:contributing -->
<!-- template:end:contributing -->

<!-- template:begin:license -->
<!-- template:end:license -->
