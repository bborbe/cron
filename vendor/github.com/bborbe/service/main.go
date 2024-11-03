// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/bborbe/argument/v2"
	libsentry "github.com/bborbe/sentry"
	"github.com/getsentry/sentry-go"
	"github.com/golang/glog"
)

//counterfeiter:generate -o mocks/service-application.go --fake-name ServiceApplication . Application
type Application interface {
	Run(ctx context.Context, sentryClient libsentry.Client) error
}

func Main(
	ctx context.Context,
	app Application,
	sentryDSN *string,
	sentryProxy *string,
	fns ...OptionsFn,
) int {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	runtime.GOMAXPROCS(runtime.NumCPU())
	_ = flag.Set("logtostderr", "true")
	_ = flag.Set("v", "2")

	time.Local = time.UTC
	glog.V(2).Infof("set global timezone to UTC")

	if err := argument.Parse(ctx, app); err != nil {
		glog.Errorf("parse app failed: %v", err)
		return 4
	}

	options := NewOptions(fns...)

	if sentryDSN == nil {
		glog.Errorf("sentryDSN args missing")
		return 3
	}
	httpTransport := http.DefaultTransport
	if sentryProxy != nil {
		httpTransport = libsentry.NewProxyRoundTripper(
			httpTransport,
			*sentryProxy,
		)
		glog.V(2).Infof("use sentryProxy %s", *sentryProxy)
	}
	sentryClient, err := libsentry.NewClient(
		ctx,
		sentry.ClientOptions{
			Dsn:              *sentryDSN,
			TracesSampleRate: 1.0,
			HTTPTransport:    httpTransport,
		},
		options.ExcludeErrors...,
	)
	if err != nil {
		glog.Errorf("setting up Sentry failed: %+v", err)
		return 2
	}
	defer func() {
		_ = sentryClient.Flush(2 * time.Second)
		_ = sentryClient.Close()
	}()

	service := NewService(
		sentryClient,
		app,
	)

	glog.V(0).Infof("application started")
	if err := service.Run(contextWithSig(ctx)); err != nil {
		glog.Error(err)
		return 1
	}
	glog.V(0).Infof("application finished")
	return 0
}

func contextWithSig(ctx context.Context) context.Context {
	ctxWithCancel, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()

		signalCh := make(chan os.Signal, 1)
		defer close(signalCh)

		signal.Notify(signalCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		select {
		case signal, ok := <-signalCh:
			if !ok {
				glog.V(2).Infof("signal channel closed => cancel context ")
				return
			}
			glog.V(2).Infof("got signal %s => cancel context ", signal)
		case <-ctx.Done():
		}
	}()

	return ctxWithCancel
}
