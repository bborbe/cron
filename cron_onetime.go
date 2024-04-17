// Copyright (c) 2019 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import (
	"context"

	"github.com/bborbe/errors"
	"github.com/bborbe/run"
	"github.com/golang/glog"
)

func NewOneTimeCron(
	action run.Runnable,
) CronJob {
	return &cronOneTime{
		action: action,
	}
}

type cronOneTime struct {
	action run.Runnable
}

func (c *cronOneTime) Run(ctx context.Context) error {
	glog.V(4).Infof("run cron action started")
	if err := c.action.Run(ctx); err != nil {
		return errors.Wrapf(ctx, err, "run cron action failed")
	}
	glog.V(4).Infof("run cron action completed")
	return nil
}
