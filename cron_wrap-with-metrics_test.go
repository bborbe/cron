// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron_test

import (
	"context"
	"errors"

	"github.com/bborbe/run"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cron"
)

var _ = Describe("WrapWithMetrics", func() {
	var ctx context.Context
	var err error
	var actionError error
	var actionCalled bool

	BeforeEach(func() {
		ctx = context.Background()
		actionError = nil
		actionCalled = false
	})

	JustBeforeEach(func() {
		action := run.Func(func(ctx context.Context) error {
			actionCalled = true
			return actionError
		})
		wrappedAction := cron.WrapWithMetrics("test-job", action)
		err = wrappedAction.Run(ctx)
	})

	Context("successful execution", func() {
		BeforeEach(func() {
			actionError = nil
		})

		It("returns no error", func() {
			Expect(err).To(BeNil())
		})

		It("calls the wrapped action", func() {
			Expect(actionCalled).To(BeTrue())
		})
	})

	Context("failed execution", func() {
		BeforeEach(func() {
			actionError = errors.New("test error")
		})

		It("returns the error", func() {
			Expect(err).To(Equal(actionError))
		})

		It("calls the wrapped action", func() {
			Expect(actionCalled).To(BeTrue())
		})
	})

	Context("with different job names", func() {
		It("accepts different job names", func() {
			action := run.Func(func(ctx context.Context) error {
				return nil
			})

			wrappedAction1 := cron.WrapWithMetrics("job-1", action)
			wrappedAction2 := cron.WrapWithMetrics("job-2", action)

			Expect(wrappedAction1.Run(ctx)).To(BeNil())
			Expect(wrappedAction2.Run(ctx)).To(BeNil())
		})
	})
})
