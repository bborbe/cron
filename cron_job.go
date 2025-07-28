// Copyright (c) 2019 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import (
	"github.com/bborbe/run"
	libtime "github.com/bborbe/time"
	"github.com/golang/glog"
)

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
