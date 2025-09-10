// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import (
	"context"
	"time"

	"github.com/bborbe/run"
)

// WrapWithMetrics wraps a runnable with Prometheus metrics collection.
// Records start, completion, failure counts and execution duration for the named job.
func WrapWithMetrics(name string, fn run.Runnable) run.Runnable {
	metrics := NewMetrics()
	return run.Func(func(ctx context.Context) error {
		start := time.Now()
		metrics.IncreaseStarted(name)

		err := fn.Run(ctx)
		duration := time.Since(start)
		metrics.ObserveDuration(name, duration.Seconds())

		if err != nil {
			metrics.IncreaseFailed(name)
			return err
		}
		metrics.IncreaseCompleted(name)
		metrics.SetLastSuccessToCurrent(name)
		return nil
	})
}
