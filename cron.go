// Copyright (c) 2019 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import (
	"context"
)

//counterfeiter:generate -o ../mocks/cron.go --fake-name Cron . Cron

// Cron represents a scheduled task that can be executed with context-based cancellation.
type Cron interface {
	// Run executes the cron job until the context is cancelled or an error occurs.
	// For recurring jobs, this method blocks and continues scheduling executions.
	// For one-time jobs, this method returns after a single execution.
	Run(ctx context.Context) error
}
