// Copyright (c) 2019 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import (
	"context"
	"time"

	"github.com/bborbe/errors"
	"github.com/bborbe/run"
	"github.com/golang/glog"
)

func NewWaitCron(
	wait time.Duration,
	action run.Runnable,
) CronJob {
	return &cronWait{
		action: action,
		wait:   wait,
	}
}

type cronWait struct {
	action run.Runnable
	wait   time.Duration
}

func (c *cronWait) Run(ctx context.Context) error {
	for {
		glog.V(4).Infof("run cron action started")
		if err := c.action.Run(ctx); err != nil {
			return errors.Wrapf(ctx, err, "run cron action failed")
		}
		glog.V(4).Infof("run cron action completed")
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.NewTimer(c.wait).C:
			glog.V(3).Infof("wait for %v completed", c.wait)
		}
	}
}
