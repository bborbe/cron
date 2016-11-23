package cron

import (
	"github.com/golang/glog"
	"time"
	"context"
)

type Cron interface {
	Run(ctx context.Context) error
}

type action func() error

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
		glog.V(4).Infof("backup cleanup started")
		if err := c.action(); err != nil {
			return err
		}
		glog.V(4).Infof("backup cleanup finished")

		if c.oneTime {
			glog.V(3).Infof("one-time => exit")
			return nil
		}

		glog.V(4).Infof("wait %v", c.wait)
		time.Sleep(c.wait)
		glog.V(4).Infof("sleep done")
	}
	return nil
}
