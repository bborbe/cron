// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron_test

import (
	libtime "github.com/bborbe/time"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cron"
)

var _ = Describe("Options", func() {
	Describe("DefaultOptions", func() {
		It("returns default values", func() {
			options := cron.DefaultOptions()

			Expect(options.Name).To(Equal("unnamed-cron"))
			Expect(options.EnableMetrics).To(BeFalse())
			Expect(options.Timeout).To(Equal(libtime.Duration(0)))
			Expect(options.ParallelSkip).To(BeFalse())
		})
	})

	Describe("Options creation", func() {
		It("allows custom configuration", func() {
			options := cron.Options{
				Name:          "test-job",
				EnableMetrics: true,
				Timeout:       libtime.Minute * 5,
				ParallelSkip:  true,
			}

			Expect(options.Name).To(Equal("test-job"))
			Expect(options.EnableMetrics).To(BeTrue())
			Expect(options.Timeout).To(Equal(libtime.Minute * 5))
			Expect(options.ParallelSkip).To(BeTrue())
		})

		It("allows partial configuration", func() {
			options := cron.Options{
				Name:          "partial-job",
				EnableMetrics: true,
			}

			Expect(options.Name).To(Equal("partial-job"))
			Expect(options.EnableMetrics).To(BeTrue())
			Expect(options.Timeout).To(Equal(libtime.Duration(0)))
			Expect(options.ParallelSkip).To(BeFalse())
		})
	})
})
