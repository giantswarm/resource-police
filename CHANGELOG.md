# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.1.1] - 2021-05-03

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

[Unreleased]: https://github.com/giantswarm/resource-police/compare/v1.1.1...HEAD
[1.1.1]: https://github.com/giantswarm/resource-police/compare/v1.0.0...v1.1.1
[1.0.0]: https://github.com/giantswarm/resource-police/compare/v0.2.5...v1.0.0
[0.2.5]: https://github.com/giantswarm/resource-police/compare/v0.2.4...v0.2.5
[0.2.4]: https://github.com/giantswarm/resource-police/compare/v0.2.3...v0.2.4
[0.2.3]: https://github.com/giantswarm/resource-police/compare/v0.2.1...v0.2.3
[0.2.1]: https://github.com/giantswarm/resource-police/compare/v0.0.0...v0.2.1
