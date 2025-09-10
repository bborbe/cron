// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import (
	libtime "github.com/bborbe/time"
)

// Options configures behavior for cron jobs with wrappers applied.
type Options struct {
	// Name identifies the cron job for logging and metrics.
	Name string
	// EnableMetrics enables duration and execution metrics collection.
	EnableMetrics bool
	// Timeout sets the maximum duration allowed for individual action executions.
	// A value of 0 disables timeout enforcement.
	Timeout libtime.Duration
	// ParallelSkip prevents multiple instances of the same cron from running concurrently.
	ParallelSkip bool
}

// DefaultOptions returns a new Options with default values.
// The default configuration has metrics disabled, no timeout, and allows parallel execution.
func DefaultOptions() Options {
	return Options{
		Name:          "unnamed-cron",
		EnableMetrics: false,
		Timeout:       0, // disabled
		ParallelSkip:  false,
	}
}
