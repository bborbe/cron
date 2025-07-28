// Copyright (c) 2019 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import (
	"context"
	"time"

	"github.com/bborbe/run"
	"github.com/golang/glog"
)

//go:generate go run -mod=mod github.com/maxbrunsfeld/counterfeiter/v6 -o mocks/cron-job.go --fake-name CronJob . CronJob
type CronJob interface {
	Run(ctx context.Context) error
}

func NewCronJob(
	oneTime bool,
	expression Expression,
	wait time.Duration,
	action run.Runnable,
) CronJob {
	return &cronJob{
		oneTime:    oneTime,
		expression: expression,
		wait:       wait,
		action:     action,
		options:    DefaultCronJobOptions(),
	}
}

func NewCronJobWithOptions(
	oneTime bool,
	expression Expression,
	wait time.Duration,
	action run.Runnable,
	options CronJobOptions,
) CronJob {
	return &cronJob{
		oneTime:    oneTime,
		expression: expression,
		wait:       wait,
		action:     action,
		options:    options,
	}
}

type cronJob struct {
	oneTime    bool
	expression Expression
	wait       time.Duration
	action     run.Runnable
	options    CronJobOptions
}

func (c *cronJob) Run(ctx context.Context) error {
	// Apply wrappers to the action based on options
	wrappedAction := c.action

	// Apply timeout wrapper first (innermost)
	if c.options.Timeout > 0 {
		wrappedAction = WrapWithTimeout(c.options.Name, c.options.Timeout, wrappedAction)
	}

	// Apply metrics wrapper (outermost)
	if c.options.EnableMetrics {
		wrappedAction = WrapWithMetrics(c.options.Name, wrappedAction)
	}

	// Apply parallel skip wrapper if enabled
	if c.options.ParallelSkip {
		parallelSkipper := run.NewParallelSkipper()
		wrappedAction = parallelSkipper.SkipParallel(wrappedAction.Run)
	}

	var runner Cron
	if c.oneTime {
		glog.V(2).Infof("create one-time cron")
		runner = NewOneTimeCron(wrappedAction)
	} else if len(c.expression) > 0 {
		glog.V(2).Infof("create cron with expression %s", c.expression)
		runner = NewExpressionCron(
			c.expression,
			wrappedAction,
		)
	} else {
		glog.V(2).Infof("create cron with wait %v", c.wait)
		runner = NewIntervalCron(
			c.wait,
			wrappedAction,
		)
	}
	return runner.Run(ctx)
}
