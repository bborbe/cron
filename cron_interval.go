// Copyright (c) 2019 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import (
	"context"
	"time"

	"github.com/bborbe/errors"
	"github.com/bborbe/run"
	libtime "github.com/bborbe/time"
	"github.com/golang/glog"
)

// NewWaitCron
// Deprecated: use NewIntervalCron instead
func NewWaitCron(
	wait libtime.Duration,
	action run.Runnable,
) run.Runnable {
	return NewIntervalCron(
		wait,
		action,
	)
}

// NewIntervalCron creates a cron job that executes at fixed time intervals.
// The job runs continuously with the specified wait duration between executions.
func NewIntervalCron(
	wait libtime.Duration,
	action run.Runnable,
) run.Runnable {
	return &intervalCron{
		action: action,
		wait:   wait,
	}
}

// NewIntervalCronWithOptions creates an interval-based cron job with configurable options.
// Applies timeout, metrics, and parallel execution controls to individual action executions.
func NewIntervalCronWithOptions(
	wait libtime.Duration,
	action run.Runnable,
	options Options,
) run.Runnable {
	return NewIntervalCron(
		wait,
		WrapWithOptions(
			action,
			options,
		),
	)
}

type intervalCron struct {
	action run.Runnable
	wait   libtime.Duration
}

func (c *intervalCron) Run(ctx context.Context) error {
	for {
		glog.V(4).Infof("run cron action started")
		if err := c.action.Run(ctx); err != nil {
			return errors.Wrapf(ctx, err, "run cron action failed")
		}
		glog.V(4).Infof("run cron action completed")
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.NewTimer(c.wait.Duration()).C:
			glog.V(3).Infof("wait for %v completed", c.wait)
		}
	}
}
