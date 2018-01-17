# Changelog
All notable changes to this project will be documented in this file.

## [Unreleased]

## [0.1.1] - 2018-01-17
### Added
- "Source on Github" sidebar link.
- `--http.throttle` flag, allowing the limitation of concurrent requests.

### Changed
- fix: `Content-Type` header for filtered requests was improperly being set
to `application/json` (should have been `text/plain`.)
- fix: spinner icon for each entry on the bulk lookup page was misaligned.

## [0.1.0] - 2017-10-13

- Initial release.

[Unreleased]: https://github.com/lrstanley/nagios-notify-irc/compare/v0.1.1...HEAD
[0.1.1]: https://github.com/lrstanley/nagios-notify-irc/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/lrstanley/geoip/tree/v0.1.0
