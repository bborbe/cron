// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron_test

import (
	"context"
	"errors"
	stdtime "time"

	"github.com/bborbe/run"
	libtime "github.com/bborbe/time"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"

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

var _ = Describe("WrapWithMetrics duration clock", func() {
	var ctx context.Context
	var originalNow func() stdtime.Time

	BeforeEach(func() {
		ctx = context.Background()
		originalNow = libtime.Now
	})
	AfterEach(func() {
		libtime.Now = originalNow
	})

	// Regression: start and elapsed must read the SAME clock. If one of the
	// two call sites still used the real clock, the observed duration would be
	// ~0 instead of the 42s the fake clock advances by.
	It("measures the duration with the injected clock", func() {
		current := stdtime.Date(2026, 8, 13, 12, 0, 0, 0, stdtime.UTC)
		libtime.Now = func() stdtime.Time { return current }

		action := run.Func(func(ctx context.Context) error {
			current = current.Add(42 * stdtime.Second)
			return nil
		})

		before := histogramSum("clock-job")
		Expect(cron.WrapWithMetrics("clock-job", action).Run(ctx)).To(BeNil())
		Expect(histogramSum("clock-job") - before).To(BeNumerically("~", 42.0, 0.001))
	})
})

// histogramSum reads the cron_job_duration_seconds sum for the given job name
// straight off the default registry.
func histogramSum(name string) float64 {
	families, err := prometheus.DefaultGatherer.Gather()
	Expect(err).To(BeNil())
	for _, family := range families {
		if family.GetName() != "cron_job_duration_seconds" {
			continue
		}
		for _, metric := range family.GetMetric() {
			for _, label := range metric.GetLabel() {
				if label.GetName() == "name" && label.GetValue() == name {
					return metric.GetHistogram().GetSampleSum()
				}
			}
		}
	}
	return 0
}
