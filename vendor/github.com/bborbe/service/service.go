// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package service

import (
	"context"

	"github.com/bborbe/errors"
	libsentry "github.com/bborbe/sentry"
	"github.com/getsentry/sentry-go"
	"github.com/golang/glog"
)

//counterfeiter:generate -o mocks/service.go --fake-name Service . Service
type Service interface {
	Run(ctx context.Context) error
}

func NewService(
	sentryClient libsentry.Client,
	app Application,
) Service {
	return &service{
		app:          app,
		sentryClient: sentryClient,
	}
}

type service struct {
	sentryClient libsentry.Client
	app          Application
}

func (s *service) Run(ctx context.Context) error {
	if err := s.app.Run(ctx, s.sentryClient); err != nil {
		s.sentryClient.CaptureException(
			err,
			&sentry.EventHint{
				Context:           ctx,
				OriginalException: err,
			},
			sentry.NewScope(),
		)
		return errors.Wrapf(ctx, err, "application failed")
	}
	glog.V(4).Infof("run finished without error")
	return nil
}
