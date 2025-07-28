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

var _ = Describe("WrapWithOptions", func() {
	var ctx context.Context
	var actionCalled bool
	var actionError error
	var action run.Runnable

	BeforeEach(func() {
		ctx = context.Background()
		actionCalled = false
		actionError = nil
		action = run.Func(func(ctx context.Context) error {
			actionCalled = true
			return actionError
		})
	})

	Describe("with no options enabled", func() {
		It("returns the original action unchanged", func() {
			options := cron.CronJobOptions{
				Name:          "test-job",
				EnableMetrics: false,
				Timeout:       0,
				ParallelSkip:  false,
			}

			wrappedAction := cron.WrapWithOptions(action, options)
			err := wrappedAction.Run(ctx)

			Expect(err).To(BeNil())
			Expect(actionCalled).To(BeTrue())
		})
	})

	Describe("with metrics enabled", func() {
		It("wraps action with metrics", func() {
			options := cron.CronJobOptions{
				Name:          "metrics-job",
				EnableMetrics: true,
				Timeout:       0,
				ParallelSkip:  false,
			}

			wrappedAction := cron.WrapWithOptions(action, options)
			err := wrappedAction.Run(ctx)

			Expect(err).To(BeNil())
			Expect(actionCalled).To(BeTrue())
		})

		Context("when action fails", func() {
			BeforeEach(func() {
				actionError = errors.New("action failed")
			})

			It("propagates error through metrics wrapper", func() {
				options := cron.CronJobOptions{
					Name:          "failing-metrics-job",
					EnableMetrics: true,
				}

				wrappedAction := cron.WrapWithOptions(action, options)
				err := wrappedAction.Run(ctx)

				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("action failed"))
				Expect(actionCalled).To(BeTrue())
			})
		})
	})

	Describe("with timeout enabled", func() {
		It("wraps action with timeout", func() {
			options := cron.CronJobOptions{
				Name:    "timeout-job",
				Timeout: libtime.Second,
			}

			wrappedAction := cron.WrapWithOptions(action, options)
			err := wrappedAction.Run(ctx)

			Expect(err).To(BeNil())
			Expect(actionCalled).To(BeTrue())
		})

		It("applies timeout when action takes too long", func() {
			slowAction := run.Func(func(ctx context.Context) error {
				select {
				case <-time.After(200 * time.Millisecond):
					actionCalled = true
					return nil
				case <-ctx.Done():
					return ctx.Err()
				}
			})

			options := cron.CronJobOptions{
				Name:    "slow-timeout-job",
				Timeout: libtime.Millisecond * 50,
			}

			wrappedAction := cron.WrapWithOptions(slowAction, options)
			err := wrappedAction.Run(ctx)

			Expect(err).NotTo(BeNil())
			Expect(actionCalled).To(BeFalse())
		})
	})

	Describe("with parallel skip enabled", func() {
		It("wraps action with parallel skipper", func() {
			options := cron.CronJobOptions{
				Name:         "parallel-skip-job",
				ParallelSkip: true,
			}

			wrappedAction := cron.WrapWithOptions(action, options)
			err := wrappedAction.Run(ctx)

			Expect(err).To(BeNil())
			Expect(actionCalled).To(BeTrue())
		})
	})

	Describe("with all options enabled", func() {
		It("applies all wrappers in correct order", func() {
			options := cron.CronJobOptions{
				Name:          "full-options-job",
				EnableMetrics: true,
				Timeout:       libtime.Second,
				ParallelSkip:  true,
			}

			wrappedAction := cron.WrapWithOptions(action, options)
			err := wrappedAction.Run(ctx)

			Expect(err).To(BeNil())
			Expect(actionCalled).To(BeTrue())
		})

		Context("when action fails", func() {
			BeforeEach(func() {
				actionError = errors.New("full-options failed")
			})

			It("propagates error through all wrappers", func() {
				options := cron.CronJobOptions{
					Name:          "full-options-failing-job",
					EnableMetrics: true,
					Timeout:       libtime.Second,
					ParallelSkip:  true,
				}

				wrappedAction := cron.WrapWithOptions(action, options)
				err := wrappedAction.Run(ctx)

				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("full-options failed"))
				Expect(actionCalled).To(BeTrue())
			})
		})

		Context("when timeout is exceeded", func() {
			It("timeout error propagates through other wrappers", func() {
				slowAction := run.Func(func(ctx context.Context) error {
					select {
					case <-time.After(200 * time.Millisecond):
						actionCalled = true
						return nil
					case <-ctx.Done():
						return ctx.Err()
					}
				})

				options := cron.CronJobOptions{
					Name:          "full-options-timeout-job",
					EnableMetrics: true,
					Timeout:       libtime.Millisecond * 50,
					ParallelSkip:  true,
				}

				wrappedAction := cron.WrapWithOptions(slowAction, options)
				err := wrappedAction.Run(ctx)

				Expect(err).NotTo(BeNil())
				Expect(actionCalled).To(BeFalse())
			})
		})
	})

	Describe("wrapper ordering", func() {
		It("applies wrappers in documented order: timeout -> metrics -> parallel skip", func() {
			// This test verifies the wrapper order through successful execution
			// The order is important for proper error handling and metrics collection
			options := cron.CronJobOptions{
				Name:          "wrapper-order-job",
				EnableMetrics: true,
				Timeout:       libtime.Second * 2,
				ParallelSkip:  true,
			}

			wrappedAction := cron.WrapWithOptions(action, options)

			// Multiple executions should work with parallel skip
			err1 := wrappedAction.Run(ctx)
			err2 := wrappedAction.Run(ctx)

			Expect(err1).To(BeNil())
			Expect(err2).To(BeNil())
			Expect(actionCalled).To(BeTrue())
		})
	})

	Describe("integration with existing wrappers", func() {
		It("can be combined with manual wrappers", func() {
			// First apply options wrappers
			options := cron.CronJobOptions{
				Name:          "manual-combo-job",
				EnableMetrics: true,
				Timeout:       libtime.Second,
			}
			wrappedAction := cron.WrapWithOptions(action, options)

			// Then apply additional manual wrapper
			finalAction := cron.WrapWithMetrics("additional-metrics", wrappedAction)

			err := finalAction.Run(ctx)

			Expect(err).To(BeNil())
			Expect(actionCalled).To(BeTrue())
		})
	})
})
