// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import (
	"context"

	"github.com/bborbe/run"
	libtime "github.com/bborbe/time"
	"github.com/golang/glog"
)

func WrapWithTimeout(name string, timeout libtime.Duration, fn run.Runnable) run.Runnable {
	if timeout <= 0 {
		glog.V(3).Infof("timeout is disabled for cron '%s'", name)
		return fn
	}
	return run.Func(func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, timeout.Duration())
		defer cancel()
		glog.V(3).Infof("add timeout %v to cron '%s'", timeout, name)
		return fn.Run(ctx)
	})
}
