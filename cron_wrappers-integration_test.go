// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron_test

import (
	"context"
	"errors"
	"time"

	"github.com/bborbe/run"
	libtime "github.com/bborbe/time"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cron"
)

var _ = Describe("Wrapper Integration", func() {
	var ctx context.Context
	var actionCalled bool
	var actionError error

	BeforeEach(func() {
		ctx = context.Background()
		actionCalled = false
		actionError = nil
	})

	Describe("Timeout + Metrics combination", func() {
		Context("successful execution", func() {
			It("applies both wrappers correctly", func() {
				action := run.Func(func(ctx context.Context) error {
					actionCalled = true
					return actionError
				})

				// Apply timeout first, then metrics
				wrappedAction := cron.WrapWithTimeout("integration-job", 1*libtime.Second, action)
				wrappedAction = cron.WrapWithMetrics("integration-job", wrappedAction)

				err := wrappedAction.Run(ctx)

				Expect(err).To(BeNil())
				Expect(actionCalled).To(BeTrue())
			})
		})

		Context("timeout exceeded", func() {
			It("metrics captures the timeout error", func() {
				action := run.Func(func(ctx context.Context) error {
					select {
					case <-time.After(200 * time.Millisecond):
						actionCalled = true
						return nil
					case <-ctx.Done():
						return ctx.Err()
					}
				})

				// Apply timeout first, then metrics
				wrappedAction := cron.WrapWithTimeout("timeout-job", 50*libtime.Millisecond, action)
				wrappedAction = cron.WrapWithMetrics("timeout-job", wrappedAction)

				err := wrappedAction.Run(ctx)

				Expect(err).NotTo(BeNil())
				Expect(actionCalled).To(BeFalse())
			})
		})

		Context("action error", func() {
			BeforeEach(func() {
				actionError = errors.New("action failed")
			})

			It("metrics captures the action error", func() {
				action := run.Func(func(ctx context.Context) error {
					actionCalled = true
					return actionError
				})

				// Apply timeout first, then metrics
				wrappedAction := cron.WrapWithTimeout("error-job", 1*libtime.Second, action)
				wrappedAction = cron.WrapWithMetrics("error-job", wrappedAction)

				err := wrappedAction.Run(ctx)

				Expect(err).To(Equal(actionError))
				Expect(actionCalled).To(BeTrue())
			})
		})
	})

	Describe("Reverse order: Metrics + Timeout", func() {
		It("works correctly with different wrapper order", func() {
			action := run.Func(func(ctx context.Context) error {
				actionCalled = true
				return actionError
			})

			// Apply metrics first, then timeout
			wrappedAction := cron.WrapWithMetrics("reverse-job", action)
			wrappedAction = cron.WrapWithTimeout("reverse-job", 1*libtime.Second, wrappedAction)

			err := wrappedAction.Run(ctx)

			Expect(err).To(BeNil())
			Expect(actionCalled).To(BeTrue())
		})
	})

	Describe("Multiple timeout layers", func() {
		It("respects the innermost timeout", func() {
			action := run.Func(func(ctx context.Context) error {
				select {
				case <-time.After(150 * time.Millisecond):
					actionCalled = true
					return nil
				case <-ctx.Done():
					return ctx.Err()
				}
			})

			// Apply multiple timeout layers - inner one should win
			wrappedAction := cron.WrapWithTimeout(
				"multi-timeout-job",
				50*libtime.Millisecond,
				action,
			)
			wrappedAction = cron.WrapWithTimeout(
				"multi-timeout-job",
				200*libtime.Millisecond,
				wrappedAction,
			)

			err := wrappedAction.Run(ctx)

			Expect(err).NotTo(BeNil())
			Expect(actionCalled).To(BeFalse())
		})
	})

	Describe("Multiple metrics layers", func() {
		It("applies metrics multiple times", func() {
			action := run.Func(func(ctx context.Context) error {
				actionCalled = true
				return actionError
			})

			// Apply multiple metrics layers - both should work
			wrappedAction := cron.WrapWithMetrics("inner-metrics-job", action)
			wrappedAction = cron.WrapWithMetrics("outer-metrics-job", wrappedAction)

			err := wrappedAction.Run(ctx)

			Expect(err).To(BeNil())
			Expect(actionCalled).To(BeTrue())
		})
	})

	Describe("Complex wrapper chains", func() {
		It("handles timeout -> metrics -> timeout -> metrics chain", func() {
			action := run.Func(func(ctx context.Context) error {
				actionCalled = true
				return actionError
			})

			// Complex chain: timeout -> metrics -> timeout -> metrics
			wrappedAction := cron.WrapWithTimeout("complex-job", 1*libtime.Second, action)
			wrappedAction = cron.WrapWithMetrics("complex-job-1", wrappedAction)
			wrappedAction = cron.WrapWithTimeout("complex-job", 2*libtime.Second, wrappedAction)
			wrappedAction = cron.WrapWithMetrics("complex-job-2", wrappedAction)

			err := wrappedAction.Run(ctx)

			Expect(err).To(BeNil())
			Expect(actionCalled).To(BeTrue())
		})
	})

	Describe("Disabled timeout in chain", func() {
		It("handles disabled timeout gracefully", func() {
			action := run.Func(func(ctx context.Context) error {
				actionCalled = true
				return actionError
			})

			// Mix enabled and disabled timeouts
			wrappedAction := cron.WrapWithTimeout("disabled-timeout-job", 0, action) // disabled
			wrappedAction = cron.WrapWithMetrics("disabled-timeout-job", wrappedAction)
			wrappedAction = cron.WrapWithTimeout(
				"disabled-timeout-job",
				1*libtime.Second,
				wrappedAction,
			) // enabled

			err := wrappedAction.Run(ctx)

			Expect(err).To(BeNil())
			Expect(actionCalled).To(BeTrue())
		})
	})
})
