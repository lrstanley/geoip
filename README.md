<!-- template:begin:header -->
<!-- do not edit anything in this "template" block, its auto-generated -->

<p align="center">geoip -- :globe_with_meridians: Geolocation API service -- Run it yourself! | alternative to freegeoip.net</p>
<p align="center">
  <a href="https://github.com/lrstanley/geoip/releases">
    <img title="Release Downloads" src="https://img.shields.io/github/downloads/lrstanley/geoip/total?style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/geoip/tags">
    <img title="Latest Semver Tag" src="https://img.shields.io/github/v/tag/lrstanley/geoip?style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/geoip/commits/master">
    <img title="Last commit" src="https://img.shields.io/github/last-commit/lrstanley/geoip?style=flat-square">
  </a>




  <a href="https://github.com/lrstanley/geoip/actions?query=workflow%3Atest+event%3Apush">
    <img title="GitHub Workflow Status (test @ master)" src="https://img.shields.io/github/workflow/status/lrstanley/geoip/test/master?label=test&style=flat-square&event=push">
  </a>

  <a href="https://codecov.io/gh/lrstanley/geoip">
    <img title="Code Coverage" src="https://img.shields.io/codecov/c/github/lrstanley/geoip/master?style=flat-square">
  </a>

  <a href="https://pkg.go.dev/github.com/lrstanley/geoip">
    <img title="Go Documentation" src="https://pkg.go.dev/badge/github.com/lrstanley/geoip?style=flat-square">
  </a>
  <a href="https://goreportcard.com/report/github.com/lrstanley/geoip">
    <img title="Go Report Card" src="https://goreportcard.com/badge/github.com/lrstanley/geoip?style=flat-square">
  </a>
</p>
<p align="center">
  <a href="https://github.com/lrstanley/geoip/issues?q=is:open+is:issue+label:bug">
    <img title="Bug reports" src="https://img.shields.io/github/issues/lrstanley/geoip/bug?label=issues&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/geoip/issues?q=is:open+is:issue+label:enhancement">
    <img title="Feature requests" src="https://img.shields.io/github/issues/lrstanley/geoip/enhancement?label=feature%20requests&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/geoip/pulls">
    <img title="Open Pull Requests" src="https://img.shields.io/github/issues-pr/lrstanley/geoip?label=prs&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/geoip/releases">
    <img title="Latest Semver Release" src="https://img.shields.io/github/v/release/lrstanley/geoip?style=flat-square">
    <img title="Latest Release Date" src="https://img.shields.io/github/release-date/lrstanley/geoip?label=date&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/geoip/discussions/new?category=q-a">
    <img title="Ask a Question" src="https://img.shields.io/badge/support-ask_a_question!-blue?style=flat-square">
  </a>
  <a href="https://liam.sh/chat"><img src="https://img.shields.io/badge/discord-bytecord-blue.svg?style=flat-square" title="Discord Chat"></a>
</p>
<!-- template:end:header -->

<!-- template:begin:toc -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :link: Table of Contents

  - [Installation](#computer-installation)
    - [Container Images (ghcr)](#whale-container-images-ghcr)
    - [Build From Source](#toolbox-build-from-source)
  - [Usage](#gear-usage)
    - [Example](#example)
  - [Support &amp; Assistance](#raising_hand_man-support--assistance)
  - [Contributing](#handshake-contributing)
  - [License](#balance_scale-license)
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
<!-- do not edit anything in this "template" block, its auto-generated -->
## :raising_hand_man: Support & Assistance

* :heart: Please review the [Code of Conduct](.github/CODE_OF_CONDUCT.md) for
     guidelines on ensuring everyone has the best experience interacting with
     the community.
* :raising_hand_man: Take a look at the [support](.github/SUPPORT.md) document on
     guidelines for tips on how to ask the right questions.
* :lady_beetle: For all features/bugs/issues/questions/etc, [head over here](https://github.com/lrstanley/geoip/issues/new/choose).
<!-- template:end:support -->

<!-- template:begin:contributing -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :handshake: Contributing

* :heart: Please review the [Code of Conduct](.github/CODE_OF_CONDUCT.md) for guidelines
     on ensuring everyone has the best experience interacting with the
    community.
* :clipboard: Please review the [contributing](.github/CONTRIBUTING.md) doc for submitting
     issues/a guide on submitting pull requests and helping out.
* :old_key: For anything security related, please review this repositories [security policy](https://github.com/lrstanley/geoip/security/policy).
<!-- template:end:contributing -->

<!-- template:begin:license -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :balance_scale: License

```
MIT License

Copyright (c) 2015 Liam Stanley <me@liamstanley.io>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

_Also located [here](LICENSE)_
<!-- template:end:license -->
