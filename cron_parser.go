// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import "github.com/robfig/cron/v3"

// CreateDefaultParser creates a cron expression parser with second-level precision.
// The parser supports the standard cron format with optional seconds field.
func CreateDefaultParser() cron.Parser {
	return cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
}
