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

// Expression of the cron
// cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor
// every second:
// * * * * * ?
// every minute
// 0 * * * * ?
// every 15 minute
// 0 */15 * * * ?
// every hour:
// 0 0 * * * ?
// every hour on sunday:
// 0 0 * * * 0
type Expression string

func (e Expression) String() string {
	return string(e)
}

func (e Expression) Ptr() *Expression {
	return &e
}

func (e Expression) Bytes() []byte {
	return []byte(e)
}

func CreateDefaultParser() cron.Parser {
	return cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
}

func NewExpressionCron(
	expression Expression,
	action run.Runnable,
) CronJob {
	return &cronExpression{
		expression: expression,
		action:     action,
		parser:     CreateDefaultParser(),
	}
}

type cronExpression struct {
	expression Expression
	action     run.Runnable
	parser     cron.Parser
}

func (c *cronExpression) Run(ctx context.Context) error {
	glog.V(4).Infof("register cron actions")
	schedule, err := c.parser.Parse(c.expression.String())
	if err != nil {
		return errors.Wrapf(err, "parse cron expression '%s' failed", c.expression)
	}

	cronJob := cron.New()
	cronJob.Start()
	errChan := make(chan error)
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
