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

//go:generate go run -mod=vendor github.com/maxbrunsfeld/counterfeiter/v6 -o mocks/cron-job.go --fake-name CronJob . CronJob
type CronJob interface {
	Run(ctx context.Context) error
}

func NewCronJob(
	oneTime bool,
	expression string,
	wait time.Duration,
	action run.Runnable,
) CronJob {
	return &cronJob{
		oneTime:    oneTime,
		expression: expression,
		wait:       wait,
		action:     action,
	}
}

type cronJob struct {
	oneTime    bool
	expression string
	wait       time.Duration
	action     run.Runnable
}

func (c *cronJob) Run(ctx context.Context) error {
	var runner Cron
	if c.oneTime {
		glog.V(2).Infof("create one-time cron")
		runner = NewOneTimeCron(c.action)
	} else if len(c.expression) > 0 {
		glog.V(2).Infof("create cron with expression %s", c.expression)
		runner = NewExpressionCron(
			c.expression,
			c.action,
		)
	} else {
		glog.V(2).Infof("create cron with wait %v", c.wait)
		runner = NewWaitCron(
			c.wait,
			c.action,
		)
	}
	return runner.Run(ctx)
}
