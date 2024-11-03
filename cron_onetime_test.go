// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron_test

import (
	"context"
	"errors"

	"github.com/bborbe/run"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	libcron "github.com/bborbe/cron"
)

var _ = Describe("OneTimeCron", func() {
	var ctx context.Context
	var err error
	BeforeEach(func() {
		ctx = context.Background()
	})
	Context("success", func() {
		var counter int
		BeforeEach(func() {
			counter = 0
		})
		JustBeforeEach(func() {
			b := libcron.NewOneTimeCron(run.Func(func(ctx context.Context) error {
				counter++
				return nil
			}))
			err = b.Run(ctx)
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("run one time", func() {
			Expect(counter).To(Equal(1))
		})
	})
	Context("fail", func() {
		JustBeforeEach(func() {
			b := libcron.NewOneTimeCron(run.Func(func(ctx context.Context) error {
				return errors.New("fail")
			}))
			err = b.Run(ctx)
		})
		It("return error", func() {
			Expect(err).NotTo(BeNil())
		})
	})
})
