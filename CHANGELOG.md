# Changelog

All notable changes to this project will be documented in this file.

Please choose versions by [Semantic Versioning](http://semver.org/).

* MAJOR version when you make incompatible API changes,
* MINOR version when you add functionality in a backwards-compatible manner, and
* PATCH version when you make backwards-compatible bug fixes.

## v1.5.0

- Add `WrapWithOptions()` function for centralized wrapper management
- Add explicit factory functions: `NewExpressionCronWithOptions()`, `NewIntervalCronWithOptions()`, `NewOneTimeCronWithOptions()`
- Add `CronJobOptions` struct with `github.com/bborbe/time.Duration` for timeout configuration

## v1.4.0

- Add Prometheus metrics integration with `WrapWithMetrics()` 
- Add timeout wrapper with `WrapWithTimeout()`
- Add parallel execution prevention support
- Add `github.com/prometheus/client_golang` dependency
- Maintain full backward compatibility

## v1.3.1

- rename NewWaitCron -> NewIntervalCron and add alias for old
- add github workflows
- add gitignores

## v1.3.0

- remove vendor
- go mod update

## v1.2.4

- go mod update

## v1.2.3

- go mod update

## v1.2.2

- go mod update

## v1.2.1

- go mod update

## v1.2.0

- add cron expression type
- add cmd for test cron expression easy
- refactor tests

## v1.1.0

- refactor
- return context cancel error
- use run.Runnable as action

## v1.0.1

- refactor
- go mod update

## v1.0.0

- Initial Version
