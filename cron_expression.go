// Copyright (c) 2019 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import (
	"context"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	robfig_cron "github.com/robfig/cron/v3"
)

type cronExpression struct {
	ext    string
	action action
}

func NewExpressionCron(
	ext string,
	action action,
) CronJob {
	c := new(cronExpression)
	c.ext = ext
	c.action = action
	return c
}

func (c *cronExpression) Run(ctx context.Context) error {
	glog.V(4).Infof("register cron actions")
	errChan := make(chan error)
	cron := robfig_cron.New()
	parser := robfig_cron.NewParser(robfig_cron.Second | robfig_cron.Minute | robfig_cron.Hour | robfig_cron.Dom | robfig_cron.Month | robfig_cron.Dow | robfig_cron.Descriptor)
	schedule, err := parser.Parse(c.ext)
	if err != nil {
		return errors.Wrapf(err, "parse cron expression '%s' failed", c.ext)
	}

	cron.Start()
	job := robfig_cron.FuncJob(func() {
		glog.V(4).Infof("run cron action started")
		if err := c.action(ctx); err != nil {
			glog.V(2).Infof("action failed -> exit")
			errChan <- err
		}
		glog.V(4).Infof("run cron action finished")
	})
	id := cron.Schedule(schedule, job)
	glog.V(3).Infof("scheduled job: %v", id)

	select {
	case err = <-errChan:
	case <-ctx.Done():
		err = nil
	}
	glog.V(2).Infof("stopping cron started")
	select {
	case err = <-errChan:
	case <-cron.Stop().Done():
		glog.V(2).Infof("stopping cron completed")
	}
	return err
}
