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

	Describe("ObserveDuration", func() {
		It("does not panic when called", func() {
			Expect(func() {
				metrics.ObserveDuration("test-job", 1.5)
			}).NotTo(Panic())
		})

		It("accepts different job names", func() {
			Expect(func() {
				metrics.ObserveDuration("job-1", 0.1)
				metrics.ObserveDuration("job-2", 2.3)
				metrics.ObserveDuration("job-3", 10.0)
			}).NotTo(Panic())
		})

		It("accepts various duration values", func() {
			Expect(func() {
				metrics.ObserveDuration("duration-job", 0.001)  // 1ms
				metrics.ObserveDuration("duration-job", 1.0)    // 1s
				metrics.ObserveDuration("duration-job", 60.0)   // 1min
				metrics.ObserveDuration("duration-job", 3600.0) // 1h
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
				metrics.ObserveDuration("success-job", 1.23)
				metrics.IncreaseCompleted("success-job")
				metrics.SetLastSuccessToCurrent("success-job")
			}).NotTo(Panic())
		})

		It("supports typical failure flow", func() {
			Expect(func() {
				metrics.IncreaseStarted("failure-job")
				metrics.ObserveDuration("failure-job", 0.45)
				metrics.IncreaseFailed("failure-job")
			}).NotTo(Panic())
		})

		It("supports multiple calls for same job", func() {
			Expect(func() {
				for i := 0; i < 5; i++ {
					metrics.IncreaseStarted("repeated-job")
					metrics.ObserveDuration("repeated-job", float64(i)*0.1)
					metrics.IncreaseCompleted("repeated-job")
					metrics.SetLastSuccessToCurrent("repeated-job")
				}
			}).NotTo(Panic())
		})
	})
})
