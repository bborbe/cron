// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import (
	"time"
)

type CronJobOptions struct {
	Name          string
	EnableMetrics bool
	Timeout       time.Duration
	ParallelSkip  bool
}

func DefaultCronJobOptions() CronJobOptions {
	return CronJobOptions{
		Name:          "unnamed-cron",
		EnableMetrics: false,
		Timeout:       0, // disabled
		ParallelSkip:  false,
	}
}
