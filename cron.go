package cron

import (
	"context"
	"github.com/golang/glog"
	"time"
)

type Cron interface {
	Run(ctx context.Context) error
}

type action func(ctx context.Context) error

type cron struct {
	action  action
	oneTime bool
	wait    time.Duration
}

func New(
	oneTime bool,
	wait time.Duration,
	action action,
) *cron {
	c := new(cron)
	c.action = action
	c.oneTime = oneTime
	c.wait = wait
	return c
}

func (c *cron) Run(ctx context.Context) error {
	for {
		glog.V(4).Infof("run cron action started")
		if err := c.action(ctx); err != nil {
			glog.V(2).Infof("action failed -> exit")
			return err
		}
		glog.V(4).Infof("run cron action finished")
		if c.oneTime {
			glog.V(2).Infof("one-time => exit")
			return nil
		}
		select {
		case <-ctx.Done():
			glog.V(2).Infof("context done -> exit")
			return nil
		case <-c.sleep():
			glog.V(4).Infof("sleep completed")
		}
	}
	return nil
}

func (c *cron) sleep() <-chan time.Time {
	glog.V(0).Infof("sleep for %v", c.wait)
	return time.After(c.wait)
}
