// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import (
	libtime "github.com/bborbe/time"
)

type Options struct {
	Name          string
	EnableMetrics bool
	Timeout       libtime.Duration
	ParallelSkip  bool
}

func DefaultOptions() Options {
	return Options{
		Name:          "unnamed-cron",
		EnableMetrics: false,
		Timeout:       0, // disabled
		ParallelSkip:  false,
	}
}
