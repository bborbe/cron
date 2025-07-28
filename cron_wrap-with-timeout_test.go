// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron_test

import (
	"context"
	libtime "github.com/bborbe/time"
	"time"

	"github.com/bborbe/run"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/cron"
)

var _ = Describe("WrapWithTimeout", func() {
	var ctx context.Context
	var err error
	var sleep libtime.Duration
	var timeout libtime.Duration
	var counter int
	BeforeEach(func() {
		ctx = context.Background()
		counter = 0
	})
	JustBeforeEach(func() {
		fn := cron.WrapWithTimeout(
			"test-cron",
			timeout,
			run.Func(func(ctx context.Context) error {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.NewTimer(sleep.Duration()).C:
					counter++
				}
				return nil
			}),
		)
		err = fn.Run(ctx)
	})
	Context("timeout exceeded", func() {
		BeforeEach(func() {
			sleep = libtime.Second
			timeout = libtime.Nanosecond
		})
		It("returns error", func() {
			Expect(err).NotTo(BeNil())
		})
		It("does not execute function fully", func() {
			Expect(counter).To(Equal(0))
		})
	})
	Context("timeout not exceeded", func() {
		BeforeEach(func() {
			sleep = libtime.Nanosecond
			timeout = 2 * libtime.Second
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("executes function", func() {
			Expect(counter).To(BeNumerically(">=", 1))
		})
	})
	Context("timeout disabled", func() {
		BeforeEach(func() {
			sleep = libtime.Nanosecond
			timeout = 0
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("executes function", func() {
			Expect(counter).To(BeNumerically(">=", 1))
		})
	})
	Context("negative timeout", func() {
		BeforeEach(func() {
			sleep = libtime.Nanosecond
			timeout = -1 * libtime.Second
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("executes function", func() {
			Expect(counter).To(BeNumerically(">=", 1))
		})
	})
})
