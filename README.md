## geoip -- GeoIP lookup service.

## Table of Contents
- [Installation](#installation)
  - [Ubuntu/Debian](#ubuntudebian)
  - [CentOS/Redhat](#centosredhat)
  - [Manual Install](#manual-install)
  - [Build from source](#build-from-source)
- [Usage](#usage)
  - [Example](#example)
- [Contributing](#contributing)
- [License](#license)

## Installation

Check out the [releases](https://github.com/lrstanley/geoip/releases)
page for prebuilt versions. Below are example commands of how you would install
the utility. Some of the more popular OS/distro steps are provided below, but
there are more released versions on the releases page previously mentioned.

### Ubuntu/Debian

```console
$ wget https://liam.sh/ghr/geoip_0.1.1_linux_amd64.deb
$ dpkg -i geoip_0.1.1_linux_amd64.deb
```

### CentOS/Redhat

```console
$ yum localinstall https://liam.sh/ghr/geoip_0.1.1_linux_amd64.rpm
```

### Manual Install

```console
$ wget https://liam.sh/ghr/geoip_0.1.1_linux_amd64.tar.gz
$ tar -C /usr/bin/ -xzvf geoip_0.1.1_linux_amd64.tar.gz geoip
$ chmod +x /usr/bin/geoip
```

### Build From Source

Dependencies (to build from source only):

   * [Go](https://golang.org/doc/install) (1.9 or greater, though latest
   preferred). Ensure your `$GOPATH` is setup.
   * [NodeJS](https://nodejs.org/en/download/) (v6 or greater)

Setup:

```console
$ go get -d -u github.com/lrstanley/geoip
$ cd $GOPATH/src/github.com/lrstanley/geoip
# this will show you all of the available options (to fetch dependencies,
# run in debug mode, etc.)
$ make help
$ make
$ ./geoip --help
```

For active development:

```console
$ make fetch # make sure this is ran at least once to fetch all dependencies.
$ make debug
# run this in a different window. this will rebundle the frontend assets on
# change.
$ make generate-watch
```

## Usage

```console
$ geoip --help
Usage:
  geoip [OPTIONS]

Application Options:
  -d, --debug          enable exception display and pprof endpoints (warn: dangerous)
  -q, --quiet          disable verbose output
      --db=            path to read/store Maxmind DB (default: geoip.db)
      --interval=      interval of time between database update checks (default: 2h)
      --update-url=    maxmind database file download location (must be gzipped) (default: http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz)
  -v, --version        print the version and compilation date

Cache Options:
      --cache.size=    total number of lookups to keep in ARC cache (50% most recent, 50% most requested) (default: 500)
      --cache.expire=  expiration time of cache (default: 20m)

HTTP Options:
  -b, --http.bind=     address and port to bind to (default: :8080)
      --http.proxy     obey X-Forwarded-For headers (warn: dangerous, make sure to only bind to localhost)
      --http.throttle= limit total max concurrent requests across all connections
      --http.limit=    number of requests/ip/hour (default: 2000)
      --http.cors=     cors origin domain to allow with https?:// prefix (empty => '*'; use flag multiple times)

TLS Options:
      --http.tls.use   enable tls
      --http.tls.cert= path to ssl certificate
      --http.tls.key=  path to ssl key

DNS Lookup Options:
      --dns.timeout=   max allowed duration when looking up hostnames (may cause queries to be slow) (default: 2s)
      --dns.resolver=  resolver (in host:port form) to use for dns lookups (doesn't work with windows and plan9) (can be used multiple times)
      --dns.uselocal   adds local (system) resolvers to the list of resolvers to use

Help Options:
  -h, --help           Show this help message

```

### Example

```console
$ geoip --cache.size 1000 --http.bind "localhost:8080" --http.proxy --http.limit 15000 --dns.resolver 8.8.8.8 --dns.resolver 8.8.4.4
```

## Contributing

Please review the [CONTRIBUTING](CONTRIBUTING.md) doc for submitting issues/a guide
on submitting pull requests and helping out.

## License

    LICENSE: The MIT License (MIT)
    Copyright (c) 2015 Liam Stanley <me@liamstanley.io>

    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:

    The above copyright notice and this permission notice shall be included in
    all copies or substantial portions of the Software.

    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
    SOFTWARE.
