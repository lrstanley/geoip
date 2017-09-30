## geoip -- GeoIP lookup service.

## Table of Contents
- [Installation](#installation)
  - [Ubuntu/Debian](#ubuntudebian)
  - [CentOS/Redhat](#centosredhat)
  - [Manual Install](#manual-install)
  - [Build from source](#build-from-source)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Installation

TODO

### Ubuntu/Debian

TODO

### CentOS/Redhat

TODO

### Manual Install

TODO

### Build From Source

Dependencies (to build from source only):

   * [Go](https://golang.org/doc/install) (1.9 or greater, though latest
   preferred). Ensure your `$GOPATH` is setup.
   * [NodeJS](https://nodejs.org/en/download/) (v6 or greater)

Setup:

```bash
$ go get -d -u github.com/lrstanley/geoip
$ cd $GOPATH/src/github.com/lrstanley/geoip
# this will show you all of the available options (to fetch dependencies,
# run in debug mode, etc.)
$ make help
$ make
$ ./geoip --help
```

For active development:

```bash
$ make fetch # make sure this is ran at least once to fetch all dependencies.
$ make debug
# run this in a different window. this will rebundle the frontend assets on
# change.
$ make generate-watch
```

## Usage

TODO

## Contributing

Please review the [CONTRIBUTING](https://github.com/lrstanley/geoip/blob/master/CONTRIBUTING.md)
doc for submitting issues/a guide on submitting pull requests and helping out.

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
