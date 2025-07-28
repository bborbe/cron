// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import (
	"context"
	"time"

	"github.com/bborbe/run"
	"github.com/golang/glog"
)

func WrapWithTimeout(name string, timeout time.Duration, fn run.Runnable) run.Runnable {
	if timeout <= 0 {
		glog.V(3).Infof("timeout is disabled for cron '%s'", name)
		return fn
	}
	return run.Func(func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		glog.V(3).Infof("add timeout %v to cron '%s'", timeout, name)
		return fn.Run(ctx)
	})
}
