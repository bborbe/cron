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

var _ = Describe("CronJob", func() {
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

	Describe("NewCronJob", func() {
		It("creates a cron job with default options", func() {
			cronJob := cron.NewCronJob(false, cron.Expression("@every 1s"), 0, action)
			Expect(cronJob).NotTo(BeNil())
		})

		Context("one-time execution", func() {
			It("executes action once", func() {
				cronJob := cron.NewCronJob(true, "", 0, action)
				err := cronJob.Run(ctx)

				Expect(err).To(BeNil())
				Expect(actionCalled).To(BeTrue())
			})
		})

		Context("expression-based execution", func() {
			It("uses expression when provided", func() {
				cronJob := cron.NewCronJob(false, cron.Expression("@every 1s"), 0, action)

				// We can't easily test the actual scheduling without waiting,
				// but we can verify the job is created properly
				Expect(cronJob).NotTo(BeNil())
			})
		})

		Context("interval-based execution", func() {
			It("uses wait duration when no expression", func() {
				cronJob := cron.NewCronJob(false, "", libtime.Second, action)

				// We can't easily test the actual scheduling without waiting,
				// but we can verify the job is created properly
				Expect(cronJob).NotTo(BeNil())
			})
		})
	})

	Describe("NewCronJobWithOptions", func() {
		var options cron.Options

		BeforeEach(func() {
			options = cron.Options{
				Name:          "test-job",
				EnableMetrics: false,
				Timeout:       0,
				ParallelSkip:  false,
			}
		})

		It("creates a cron job with custom options", func() {
			cronJob := cron.NewCronJobWithOptions(
				false,
				cron.Expression("@every 1s"),
				0,
				action,
				options,
			)
			Expect(cronJob).NotTo(BeNil())
		})

		Context("with metrics enabled", func() {
			BeforeEach(func() {
				options.EnableMetrics = true
			})

			It("wraps action with metrics", func() {
				cronJob := cron.NewCronJobWithOptions(true, "", 0, action, options)
				err := cronJob.Run(ctx)

				Expect(err).To(BeNil())
				Expect(actionCalled).To(BeTrue())
			})
		})

		Context("with timeout enabled", func() {
			BeforeEach(func() {
				options.Timeout = libtime.Millisecond * 100
			})

			It("applies timeout to action", func() {
				// Create an action that takes longer than the timeout
				slowAction := run.Func(func(ctx context.Context) error {
					select {
					case <-time.After(200 * time.Millisecond):
						actionCalled = true
						return nil
					case <-ctx.Done():
						return ctx.Err()
					}
				})

				cronJob := cron.NewCronJobWithOptions(true, "", 0, slowAction, options)
				err := cronJob.Run(ctx)

				Expect(err).NotTo(BeNil())
				Expect(actionCalled).To(BeFalse())
			})
		})

		Context("with parallel skip enabled", func() {
			BeforeEach(func() {
				options.ParallelSkip = true
			})

			It("wraps action with parallel skipper", func() {
				cronJob := cron.NewCronJobWithOptions(true, "", 0, action, options)
				err := cronJob.Run(ctx)

				Expect(err).To(BeNil())
				Expect(actionCalled).To(BeTrue())
			})
		})

		Context("with multiple wrappers", func() {
			BeforeEach(func() {
				options.EnableMetrics = true
				options.Timeout = libtime.Second
				options.ParallelSkip = true
			})

			It("applies all wrappers correctly", func() {
				cronJob := cron.NewCronJobWithOptions(true, "", 0, action, options)
				err := cronJob.Run(ctx)

				Expect(err).To(BeNil())
				Expect(actionCalled).To(BeTrue())
			})
		})

		Context("error handling", func() {
			BeforeEach(func() {
				actionError = errors.New("test error")
				options.EnableMetrics = true
			})

			It("propagates errors through wrappers", func() {
				cronJob := cron.NewCronJobWithOptions(true, "", 0, action, options)
				err := cronJob.Run(ctx)

				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("test error"))
				Expect(actionCalled).To(BeTrue())
			})
		})
	})
})
