# Change Log

All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]

Nothing yet.

## [1.0.1] - 2016-04-05

### Added

- `-W` option for timeout
- `-net` option for network (for example `tcp6`)

### Changed

- `Ping` calls `net.DialTimeout` instead of `net.Dial`
- `Ping` takes `host, port string` instead of `host string, port int`
- Improved project organization
- Unit tests use random TCP port instead of 1234

## [1.0.0] - 2016-03-28

- Ping TCP ports using `./portping host port`
- `-c` flag for ping count

[Unreleased]: https://github.com/janosgyerik/portping/compare/v1.0.1...HEAD
[1.0.1]: https://github.com/janosgyerik/portping/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/janosgyerik/portping/compare/065f1d5659af522502d4632085322b1ab65c009f...v1.0.0
