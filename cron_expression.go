package cron

import (
	"context"
	"github.com/golang/glog"
	robfig_cron "github.com/robfig/cron"
)

type cronExpression struct {
	action     action
	expression string
}

func NewExpressionCron(
	expression string,
	action action,
) *cronWait {
	c := new(cronWait)
	c.action = action
	return c
}

func (c *cronExpression) Run(ctx context.Context) error {
	glog.V(4).Infof("run cron action started")
	errChan := make(chan error)
	cron := robfig_cron.New()
	cron.Start()
	defer cron.Stop()
	cron.AddFunc(c.expression, func() {
		glog.V(4).Infof("run cron action started")
		if err := c.action(ctx); err != nil {
			glog.V(2).Infof("action failed -> exit")
			errChan <- err
		}
	})
	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return nil
	}
}
