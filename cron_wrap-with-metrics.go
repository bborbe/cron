// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import (
	"context"

	"github.com/bborbe/run"
)

func WrapWithMetrics(name string, fn run.Runnable) run.Runnable {
	metrics := NewMetrics()
	return run.Func(func(ctx context.Context) error {
		metrics.IncreaseStarted(name)
		if err := fn.Run(ctx); err != nil {
			metrics.IncreaseFailed(name)
			return err
		}
		metrics.IncreaseCompleted(name)
		metrics.SetLastSuccessToCurrent(name)
		return nil
	})
}
