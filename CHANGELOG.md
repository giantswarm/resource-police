# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Add team label in resources.

## [1.3.1] - 2023-11-30

### Fixed

- Delete unused policy exception.

## [1.3.0] - 2023-11-29

### Changed

- Bump to Go 1.21 and newer dependencies
- Switch to `batch/v1` apiVersion for the Cronjob.

### Added

- Add policy exception.

## [1.2.0] - 2022-05-18

### Added

- Add provider label.

## [1.1.1] - 2021-05-03

- Fix the age grouping.
- Add a `--dry-run` flag to simplify testing in development.

## [1.0.0] - 2021-04-12

- Changed the entire data fetching logic to use Prometheus/Cortex (https://giantswarm.grafana.net/) data. [#43](https://github.com/giantswarm/resource-police/pull/43)

## [0.2.5] - 2021-01-14

### Fixed

- Set service account name for cronjob.

## [0.2.4] - 2020-08-07

## [0.2.3] - 2020-08-07

## [0.2.1] - 2020-07-28

### Fixed

- Moved report template to go source file so it can be read when running in a container.

[Unreleased]: https://github.com/giantswarm/resource-police/compare/v1.3.1...HEAD
[1.3.1]: https://github.com/giantswarm/resource-police/compare/v1.3.0...v1.3.1
[1.3.0]: https://github.com/giantswarm/resource-police/compare/v1.2.0...v1.3.0
[1.2.0]: https://github.com/giantswarm/resource-police/compare/v1.1.1...v1.2.0
[1.1.1]: https://github.com/giantswarm/resource-police/compare/v1.0.0...v1.1.1
[1.0.0]: https://github.com/giantswarm/resource-police/compare/v0.2.5...v1.0.0
[0.2.5]: https://github.com/giantswarm/resource-police/compare/v0.2.4...v0.2.5
[0.2.4]: https://github.com/giantswarm/resource-police/compare/v0.2.3...v0.2.4
[0.2.3]: https://github.com/giantswarm/resource-police/compare/v0.2.1...v0.2.3
[0.2.1]: https://github.com/giantswarm/resource-police/compare/v0.0.0...v0.2.1
