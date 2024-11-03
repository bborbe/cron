// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"errors"

	"github.com/bborbe/run"
)

func Run(ctx context.Context, funcs ...run.Func) error {
	for i, fn := range funcs {
		funcs[i] = run.LogErrors(
			run.CatchPanic(
				FilterErrors(
					fn,
					context.Canceled,
				),
			),
		)
	}
	return run.CancelOnFirstFinishWait(ctx, funcs...)
}

// FilterErrors for the given func
func FilterErrors(fn run.Func, filteredErrors ...error) run.Func {
	return func(ctx context.Context) error {
		if err := fn(ctx); err != nil {
			for _, filteredError := range filteredErrors {
				if errors.Is(err, filteredError) {
					return nil
				}
			}
			return err
		}
		return nil
	}
}
