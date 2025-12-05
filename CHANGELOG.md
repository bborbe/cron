# Changelog

All notable changes to this project will be documented in this file.

Please choose versions by [Semantic Versioning](http://semver.org/).

* MAJOR version when you make incompatible API changes,
* MINOR version when you add functionality in a backwards-compatible manner, and
* PATCH version when you make backwards-compatible bug fixes.

## v1.8.0

- update go and deps

## v1.7.3
- Update dependencies: bborbe/run v1.8.1 → v1.8.2
- Update dependencies: google/osv-scanner v2.2.4 → v2.3.0
- Update dependencies: incu6us/goimports-reviser v3.10.0 → v3.11.0
- Update dependencies: getsentry/sentry-go v0.36.0 → v0.36.2

## v1.7.2
- Update Go version from 1.25.2 to 1.25.4
- Update dependencies: bborbe/errors, bborbe/run, bborbe/sentry, bborbe/service, bborbe/time
- Update dependencies: google/osv-scanner, onsi/ginkgo, securego/gosec

## v1.7.1
- Improve error handling by switching from pkg/errors to bborbe/errors
- Add context parameter to error wrapping in cron expression parsing

## v1.7.0
- Add golangci-lint configuration with comprehensive linter settings
- Add security scanning tools: gosec, osv-scanner, trivy
- Update Go version from 1.24.5 to 1.25.2
- Enhance Makefile with lint, gosec, osv-scanner, and trivy targets
- Apply golines formatting for consistent line length (max 100 characters)
- Add Trivy installation step to CI workflow
- Update dependencies including bborbe/sentry, bborbe/time, and onsi/ginkgo/v2

## v1.6.1
- Update Go version from 1.24.5 to 1.25.1
- Update dependencies to latest versions including bborbe/run, bborbe/sentry, bborbe/time, prometheus/client_golang, onsi/ginkgo/v2, and others

## v1.6.0

- **BREAKING FIX**: Fix timeout context deadline issue in `NewIntervalCronWithOptions` and `NewExpressionCronWithOptions`
- **IMPORTANT**: Timeout now applies to individual action executions instead of entire cron lifecycle
- This prevents premature termination of long-running interval/expression-based crons
- Add comprehensive Go documentation for all public interfaces, structs, functions, and methods
- Add package-level documentation explaining the three execution strategies
- All existing timeout configurations will now work correctly for repeated executions

## v1.5.2

- migrate all interval functions to use libtime.Duration consistently
- fix WrapWithTimeout to accept libtime.Duration directly
- improve API consistency across all factory functions

## v1.5.2

- Important fix of ExpressionCron and IntervalCron with Options

## v1.5.1

- rename CronJobOptions -> Options 

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
