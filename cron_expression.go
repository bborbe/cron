// Copyright (c) 2019 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import (
	"context"

	"github.com/bborbe/run"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
)

func NewExpressionCron(
	expression string,
	action run.Runnable,
) CronJob {
	return &cronExpression{
		expression: expression,
		action:     action,
	}
}

type cronExpression struct {
	expression string
	action     run.Runnable
}

func (c *cronExpression) Run(ctx context.Context) error {
	glog.V(4).Infof("register cron actions")
	errChan := make(chan error)
	cronJob := cron.New()
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	schedule, err := parser.Parse(c.expression)
	if err != nil {
		return errors.Wrapf(err, "parse cron expression '%s' failed", c.expression)
	}

	cronJob.Start()
	job := cron.FuncJob(func() {
		glog.V(4).Infof("run cron action started")
		if err := c.action.Run(ctx); err != nil {
			errChan <- err
		}
		glog.V(4).Infof("run cron action completed")
	})
	id := cronJob.Schedule(schedule, job)
	glog.V(3).Infof("scheduled job: %v", id)

	select {
	case err = <-errChan:
	case <-ctx.Done():
		err = nil
	}
	glog.V(2).Infof("stopping cron started")
	stopContext := cronJob.Stop()
	select {
	case err = <-errChan:
	case <-stopContext.Done():
		glog.V(2).Infof("stopping cron completed")
	}
	return err
}
