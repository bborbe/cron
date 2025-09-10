// Copyright (c) 2019 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import (
	"github.com/bborbe/run"
	libtime "github.com/bborbe/time"
	"github.com/golang/glog"
)

// NewCronJob creates a new cron job with automatic strategy selection based on parameters.
// Uses one-time execution if oneTime is true, expression-based scheduling if expression is provided,
// or duration-based intervals if wait duration is specified.
func NewCronJob(
	oneTime bool,
	expression Expression,
	wait libtime.Duration,
	action run.Runnable,
) run.Runnable {
	return NewCronJobWithOptions(
		oneTime,
		expression,
		wait,
		action,
		DefaultOptions(),
	)
}

// NewCronJobWithOptions creates a new cron job with configurable options.
// Applies the same strategy selection as NewCronJob but with additional wrappers for
// timeout, metrics, and parallel execution control based on the provided options.
func NewCronJobWithOptions(
	oneTime bool,
	expression Expression,
	wait libtime.Duration,
	action run.Runnable,
	options Options,
) run.Runnable {
	if oneTime {
		glog.V(2).Infof("create one-time cron")
		return NewOneTimeCronWithOptions(
			action,
			options,
		)
	} else if len(expression) > 0 {
		glog.V(2).Infof("create cron with expression %s", expression)
		return NewExpressionCronWithOptions(
			expression,
			action,
			options,
		)
	} else {
		glog.V(2).Infof("create cron with wait %v", wait)
		return NewIntervalCron(
			wait,
			action,
		)
	}
}
