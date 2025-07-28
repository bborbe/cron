# Cron

A Go library for robust cron job scheduling with multiple execution strategies, observability, and reliability features.

## Installation

```bash
go get github.com/bborbe/cron
```

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/bborbe/cron"
    "github.com/bborbe/run"
)

func main() {
    ctx := context.Background()
    
    // Create an action to run
    action := run.Func(func(ctx context.Context) error {
        fmt.Println("Job executed!")
        return nil
    })
    
    // Expression-based scheduling
    cronJob := cron.NewCronJob(false, cron.Expression("@every 1h"), 0, action)
    cronJob.Run(ctx)
}
```

### Enhanced Usage with Options

```go
package main

import (
    "context"
    "time"
    
    "github.com/bborbe/cron"
    "github.com/bborbe/run"
)

func main() {
    ctx := context.Background()
    
    action := run.Func(func(ctx context.Context) error {
        // Your job logic here
        return nil
    })
    
    // Configure options
    options := cron.CronJobOptions{
        Name:          "my-important-job",
        EnableMetrics: true,
        Timeout:       5 * time.Minute,
        ParallelSkip:  true,
    }
    
    // Create job with options
    cronJob := cron.NewCronJobWithOptions(
        false,                      // not one-time
        cron.Expression("0 */15 * * * ?"), // every 15 minutes
        0,                          // no wait duration
        action,
        options,
    )
    
    cronJob.Run(ctx)
}
```

## Execution Strategies

The library automatically selects the appropriate execution strategy based on the parameters provided to `NewCronJob` or `NewCronJobWithOptions`:

### 1. Expression-Based (Cron Expressions)

Uses standard cron expressions or special formats:

```go
// Every hour
cronJob := cron.NewCronJob(false, cron.Expression("@every 1h"), 0, action)

// Every 15 minutes using cron syntax
cronJob := cron.NewCronJob(false, cron.Expression("0 */15 * * * ?"), 0, action)

// Daily at 2:30 AM
cronJob := cron.NewCronJob(false, cron.Expression("0 30 2 * * ?"), 0, action)
```

### 2. Duration-Based (Intervals)

Uses Go `time.Duration` for simple intervals:

```go
// Every 30 seconds
cronJob := cron.NewCronJob(false, "", 30*time.Second, action)

// Every 5 minutes
cronJob := cron.NewCronJob(false, "", 5*time.Minute, action)
```

### 3. One-Time Execution

Executes the job once immediately:

```go
// Run once
cronJob := cron.NewCronJob(true, "", 0, action)
```

## Configuration Options

The `CronJobOptions` struct provides fine-grained control over job behavior:

```go
type CronJobOptions struct {
    Name          string        // Job name for metrics and logging
    EnableMetrics bool          // Enable Prometheus metrics collection  
    Timeout       time.Duration // Execution timeout (0 = disabled)
    ParallelSkip  bool          // Prevent concurrent executions
}
```

### Default Options

```go
options := cron.DefaultCronJobOptions()
// Returns:
// {
//     Name:          "unnamed-cron", 
//     EnableMetrics: false,
//     Timeout:       0,
//     ParallelSkip:  false,
// }
```

## Wrapper Functions

For maximum flexibility, you can use wrapper functions directly:

### Metrics Wrapper

```go
import "github.com/prometheus/client_golang/prometheus"

// Wrap any action with metrics
wrappedAction := cron.WrapWithMetrics("job-name", originalAction)
```

**Metrics Collected:**
- `cron_job_started{name="job-name"}` - Number of job starts
- `cron_job_completed{name="job-name"}` - Number of successful completions  
- `cron_job_failed{name="job-name"}` - Number of failures
- `cron_job_last_success{name="job-name"}` - Timestamp of last success

### Timeout Wrapper

```go
// Add 5-minute timeout
wrappedAction := cron.WrapWithTimeout("job-name", 5*time.Minute, originalAction)

// Disable timeout (≤0 duration)
wrappedAction := cron.WrapWithTimeout("job-name", 0, originalAction)
```

### Chaining Wrappers

```go
// Chain multiple wrappers (order matters)
action := cron.WrapWithTimeout("my-job", 5*time.Minute, originalAction)
action = cron.WrapWithMetrics("my-job", action)
```

## Monitoring and Observability

### Prometheus Integration

When metrics are enabled, the library automatically registers Prometheus collectors. Expose them via HTTP:

```go
import (
    "net/http"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
    // Your cron jobs with EnableMetrics: true
    
    // Expose metrics endpoint
    http.Handle("/metrics", promhttp.Handler())
    http.ListenAndServe(":8080", nil)
}
```

### Custom Metrics

Access the metrics interface directly:

```go
metrics := cron.NewCronMetrics()
metrics.IncreaseStarted("custom-job")
metrics.IncreaseCompleted("custom-job") 
metrics.SetLastSuccessToCurrent("custom-job")
```

## Error Handling

The library provides proper error propagation through all wrapper layers:

```go
action := run.Func(func(ctx context.Context) error {
    return errors.New("job failed")
})

cronJob := cron.NewCronJobWithOptions(true, "", 0, action, options)
err := cronJob.Run(ctx)
if err != nil {
    // Handle error - will include context from wrappers
    log.Printf("Cron job failed: %v", err)
}
```

## Testing

The library provides mock interfaces for testing:

```go
//go:generate go run -mod=mod github.com/maxbrunsfeld/counterfeiter/v6 -generate

import "github.com/bborbe/cron/mocks"

func TestMyCronJob(t *testing.T) {
    mockCronJob := &mocks.CronJob{}
    mockMetrics := &mocks.CronMetrics{}
    
    // Set up expectations and test your code
}
```

## Examples

### Web Scraper with Monitoring

```go
options := cron.CronJobOptions{
    Name:          "web-scraper",
    EnableMetrics: true,
    Timeout:       30 * time.Second,
    ParallelSkip:  true,
}

scraper := run.Func(func(ctx context.Context) error {
    // Scrape website logic
    return scrapeWebsite(ctx)
})

cronJob := cron.NewCronJobWithOptions(
    false,
    cron.Expression("0 */5 * * * ?"), // Every 5 minutes
    0,
    scraper,
    options,
)
```

### Database Cleanup Job

```go
cleanup := run.Func(func(ctx context.Context) error {
    return cleanupOldRecords(ctx)
})

options := cron.CronJobOptions{
    Name:          "db-cleanup",
    EnableMetrics: true,
    Timeout:       10 * time.Minute,
}

cronJob := cron.NewCronJobWithOptions(
    false,
    cron.Expression("0 0 2 * * ?"), // Daily at 2 AM
    0,
    cleanup,
    options,
)
```

## License

BSD-style license. See LICENSE file for details.
