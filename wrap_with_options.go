// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import (
	"github.com/bborbe/run"
)

// WrapWithOptions applies all configured wrappers to an action based on the provided options.
// Wrappers are applied in this order (innermost to outermost):
// 1. Timeout wrapper (if timeout > 0)
// 2. Metrics wrapper (if enabled)
// 3. Parallel skip wrapper (if enabled)
func WrapWithOptions(action run.Runnable, options Options) run.Runnable {
	wrappedAction := action

	// Apply timeout wrapper first (innermost)
	if options.Timeout.Duration() > 0 {
		wrappedAction = WrapWithTimeout(options.Name, options.Timeout.Duration(), wrappedAction)
	}

	// Apply metrics wrapper (outermost)
	if options.EnableMetrics {
		wrappedAction = WrapWithMetrics(options.Name, wrappedAction)
	}

	// Apply parallel skip wrapper if enabled
	if options.ParallelSkip {
		parallelSkipper := run.NewParallelSkipper()
		wrappedAction = parallelSkipper.SkipParallel(wrappedAction.Run)
	}

	return wrappedAction
}
