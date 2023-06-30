<!-- template:define:options
{
  "nodescription": true
}
-->
![logo](https://liam.sh/-/gh/svg/lrstanley/geoip?icon=fluent-emoji-flat%3Aglobe-showing-americas&icon.height=80&layout=left&bgcolor=rgba%282%2C+12%2C+18%2C+1%29)

<!-- template:begin:header -->
<!-- do not edit anything in this "template" block, its auto-generated -->

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
    <img title="GitHub Workflow Status (test @ master)" src="https://img.shields.io/github/actions/workflow/status/lrstanley/geoip/test.yml?branch=master&label=test&style=flat-square">
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
    - [🐳 Container Images (ghcr)](#whale-container-images-ghcr)
    - [Build From Source](#toolbox-build-from-source)
  - [Usage](#gear-usage)
    - [Example](#example)
  - [🙋‍♂️ Support &amp; Assistance](#raising_hand_man-support--assistance)
  - [Contributing](#handshake-contributing)
  - [⚖️ License](#balance_scale-license)
<!-- template:end:toc -->

## :computer: Installation

Check out the [releases](https://github.com/lrstanley/geoip/releases)
page for prebuilt versions.

### :whale: Container Images (ghcr)

```console
$ docker run -it --rm -p 8080:80 --env-file .env -v $PWD/:/data ghcr.io/lrstanley/geoip:latest
$ curl -I http://localhost:8080
HTTP/1.1 200 OK
Content-Type: text/html
Date: Thu, 06 Aug 2020 00:55:21 GMT
```

### :toolbox: Build From Source

Dependencies (to build from source only):

- [Go](https://golang.org/doc/install) (latest)
- [NodeJS](https://nodejs.org/en/download/) (v17)

Setup:

```console
git clone <repo>
cd geoip
make
./geoip --help
```

## :gear: Usage

Take a look at the [CLI usage options here](./USAGE.md).

### Example

```console
geoip --http.bind-addr "localhost:8080" --http.limit 15000 --dns.resolver 8.8.8.8 --dns.resolver 8.8.4.4
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
