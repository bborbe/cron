// Copyright (c) 2016, Benjamin Borbe <bborbe@rocketnews.de>.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cron provides scheduling functionality for executing tasks at specified intervals.
//
// This library offers three execution strategies:
//   - Expression-based scheduling using cron expressions (e.g., "@every 1h", "0 */15 * * * ?")
//   - Duration-based intervals for repeated execution at fixed time intervals
//   - One-time execution for tasks that should run only once
//
// The package includes configurable options for timeouts, metrics collection, and parallel execution control.
// All schedulers support graceful context-based cancellation.
package cron
