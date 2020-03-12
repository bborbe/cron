// Copyright (c) 2019 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import (
	"context"
)

type Cron interface {
	Run(ctx context.Context) error
}

type action func(ctx context.Context) error
