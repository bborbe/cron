// Copyright (c) 2019 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import (
	"context"

	"github.com/golang/glog"
)

func NewOneTimeCron(
	action action,
) CronJob {
	c := new(cronOneTime)
	c.action = action
	return c
}

type cronOneTime struct {
	action action
}

func (c *cronOneTime) Run(ctx context.Context) error {
	glog.V(4).Infof("run cron action started")
	if err := c.action(ctx); err != nil {
		glog.V(2).Infof("action failed -> exit")
		return err
	}
	return nil
}
