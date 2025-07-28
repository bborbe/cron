// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cron"
)

var _ = Describe("CronMetrics", func() {
	var metrics cron.Metrics

	BeforeEach(func() {
		metrics = cron.NewMetrics()
	})

	Describe("NewCronMetrics", func() {
		It("creates a new metrics instance", func() {
			Expect(metrics).NotTo(BeNil())
		})
	})

	Describe("IncreaseStarted", func() {
		It("does not panic when called", func() {
			Expect(func() {
				metrics.IncreaseStarted("test-job")
			}).NotTo(Panic())
		})

		It("accepts different job names", func() {
			Expect(func() {
				metrics.IncreaseStarted("job-1")
				metrics.IncreaseStarted("job-2")
				metrics.IncreaseStarted("job-3")
			}).NotTo(Panic())
		})
	})

	Describe("IncreaseFailed", func() {
		It("does not panic when called", func() {
			Expect(func() {
				metrics.IncreaseFailed("test-job")
			}).NotTo(Panic())
		})

		It("accepts different job names", func() {
			Expect(func() {
				metrics.IncreaseFailed("job-1")
				metrics.IncreaseFailed("job-2")
				metrics.IncreaseFailed("job-3")
			}).NotTo(Panic())
		})
	})

	Describe("IncreaseCompleted", func() {
		It("does not panic when called", func() {
			Expect(func() {
				metrics.IncreaseCompleted("test-job")
			}).NotTo(Panic())
		})

		It("accepts different job names", func() {
			Expect(func() {
				metrics.IncreaseCompleted("job-1")
				metrics.IncreaseCompleted("job-2")
				metrics.IncreaseCompleted("job-3")
			}).NotTo(Panic())
		})
	})

	Describe("SetLastSuccessToCurrent", func() {
		It("does not panic when called", func() {
			Expect(func() {
				metrics.SetLastSuccessToCurrent("test-job")
			}).NotTo(Panic())
		})

		It("accepts different job names", func() {
			Expect(func() {
				metrics.SetLastSuccessToCurrent("job-1")
				metrics.SetLastSuccessToCurrent("job-2")
				metrics.SetLastSuccessToCurrent("job-3")
			}).NotTo(Panic())
		})
	})

	Describe("Interface compliance", func() {
		It("implements CronMetrics interface", func() {
			var _ cron.Metrics = metrics
		})
	})

	Describe("Typical usage patterns", func() {
		It("supports typical success flow", func() {
			Expect(func() {
				metrics.IncreaseStarted("success-job")
				metrics.IncreaseCompleted("success-job")
				metrics.SetLastSuccessToCurrent("success-job")
			}).NotTo(Panic())
		})

		It("supports typical failure flow", func() {
			Expect(func() {
				metrics.IncreaseStarted("failure-job")
				metrics.IncreaseFailed("failure-job")
			}).NotTo(Panic())
		})

		It("supports multiple calls for same job", func() {
			Expect(func() {
				for i := 0; i < 5; i++ {
					metrics.IncreaseStarted("repeated-job")
					metrics.IncreaseCompleted("repeated-job")
					metrics.SetLastSuccessToCurrent("repeated-job")
				}
			}).NotTo(Panic())
		})
	})
})
