// Copyright (c) 2019 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import (
	"time"

	"github.com/bborbe/run"
	"github.com/golang/glog"
)

func NewCronJob(
	oneTime bool,
	expression Expression,
	wait time.Duration,
	action run.Runnable,
) run.Runnable {
	return NewCronJobWithOptions(
		oneTime,
		expression,
		wait,
		action,
		DefaultCronJobOptions(),
	)
}

func NewCronJobWithOptions(
	oneTime bool,
	expression Expression,
	wait time.Duration,
	action run.Runnable,
	options CronJobOptions,
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
