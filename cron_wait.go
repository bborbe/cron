// Copyright (c) 2019 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import (
	"context"
	"time"

	"github.com/golang/glog"
)

func NewWaitCron(
	wait time.Duration,
	action action,
) CronJob {
	return &cronWait{
		action: action,
		wait:   wait,
	}
}

type cronWait struct {
	action action
	wait   time.Duration
}

func (c *cronWait) Run(ctx context.Context) error {
	for {
		glog.V(4).Infof("run cron action started")
		if err := c.action(ctx); err != nil {
			glog.V(2).Infof("action failed -> exit")
			return err
		}
		select {
		case <-ctx.Done():
			glog.V(2).Infof("context done -> exit")
			return nil
		case <-c.sleep():
			glog.V(4).Infof("sleep completed")
		}
	}
}

func (c *cronWait) sleep() <-chan time.Time {
	glog.V(0).Infof("sleep for %v", c.wait)
	return time.After(c.wait)
}
