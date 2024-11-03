// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	stderrors "errors"

	"github.com/bborbe/sentry"
)

type Options struct {
	ExcludeErrors sentry.ExcludeErrors
}

type OptionsFn func(option *Options)

func NewOptions(fns ...OptionsFn) Options {
	options := Options{
		ExcludeErrors: sentry.ExcludeErrors{
			func(err error) bool {
				return stderrors.Is(err, context.Canceled)
			},
			func(err error) bool {
				return stderrors.Is(err, context.DeadlineExceeded)
			},
		},
	}
	for _, fn := range fns {
		fn(&options)
	}
	return options
}
