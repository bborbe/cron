// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/robfig/cron/v3"

	libcron "github.com/bborbe/cron"
)

var _ = Describe("DefaultParser", func() {
	var parser cron.Parser
	BeforeEach(func() {
		parser = libcron.CreateDefaultParser()
	})
	DescribeTable("Parse",
		func(expression libcron.Expression, expectError bool) {
			schedule, err := parser.Parse(expression.String())
			if expectError {
				Expect(schedule).To(BeNil())
				Expect(err).NotTo(BeNil())
			} else {
				Expect(schedule).NotTo(BeNil())
				Expect(err).To(BeNil())
			}
		},
		Entry("every second", libcron.Expression("* * * * * ?"), false),
		Entry("every minute", libcron.Expression("0 * * * * ?"), false),
		Entry("every hour", libcron.Expression("0 0 * * * ?"), false),
		Entry("every hour first of month", libcron.Expression("0 0 1 * * ?"), false),
		Entry("every hour in december", libcron.Expression("0 0 * * 12 ?"), false),
		Entry("every hour sundays", libcron.Expression("0 0 * * * 0"), false),
	)
})
