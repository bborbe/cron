// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"os"

	"github.com/bborbe/run"
	libsentry "github.com/bborbe/sentry"
	"github.com/bborbe/service"
	"github.com/golang/glog"

	"github.com/bborbe/cron"
)

func main() {
	app := &application{}
	os.Exit(service.Main(context.Background(), app, &app.SentryDSN, &app.SentryProxy))
}

type application struct {
	SentryDSN      string `required:"false" arg:"sentry-dsn" env:"SENTRY_DSN" usage:"SentryDSN" display:"length"`
	SentryProxy    string `required:"false" arg:"sentry-proxy" env:"SENTRY_PROXY" usage:"Sentry Proxy"`
	CronExpression string `required:"true" arg:"cron-expression" env:"CRON_EXPRESSION" usage:"Cron expression to determine when service is run" default:"@every 1h"`
}

func (a *application) Run(ctx context.Context, sentryClient libsentry.Client) error {
	return cron.NewExpressionCron(
		cron.Expression(a.CronExpression),
		run.Func(func(ctx context.Context) error {
			glog.V(2).Infof("cron executed")
			return nil
		}),
	).Run(ctx)
}
