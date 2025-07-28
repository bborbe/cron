// Copyright (c) 2024 Benjamin Borbe All rights reserved.
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

	libcron "github.com/bborbe/cron"
)

var _ = Describe("CronWait", func() {
	var ctx context.Context
	var err error
	BeforeEach(func() {
		ctx = context.Background()
	})
	Context("RunContinous", func() {
		var counter int
		BeforeEach(func() {
			counter = 0
		})
		JustBeforeEach(func() {
			b := libcron.NewWaitCron(libtime.Microsecond, run.Func(func(ctx context.Context) error {
				counter++
				return nil
			}))
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			go func() {
				time.Sleep(10 * libtime.Millisecond.Duration())
				cancel()
			}()
			err = b.Run(ctx)
		})

		It("returns error", func() {
			Expect(err).NotTo(BeNil())
			Expect(errors.Is(err, context.Canceled)).To(BeTrue())
		})
		It("increase counter", func() {
			Expect(counter).To(BeNumerically(">=", 1))
		})
	})
	Context("RunContinuousError", func() {
		JustBeforeEach(func() {
			b := libcron.NewWaitCron(libtime.Microsecond, run.Func(func(ctx context.Context) error {
				return errors.New("fail")
			}))
			err = b.Run(ctx)
		})
		It("returns error", func() {
			Expect(err).NotTo(BeNil())
		})
	})
	Context("RunContinuousCancel", func() {
		JustBeforeEach(func() {
			b := libcron.NewWaitCron(libtime.Microsecond, run.Func(func(ctx context.Context) error {
				<-ctx.Done()
				return nil
			}))
			ctx, cancel := context.WithCancel(context.Background())
			go func() {
				time.Sleep(10 * time.Millisecond)
				cancel()
			}()
			err = b.Run(ctx)
		})
		It("returns error", func() {
			Expect(err).NotTo(BeNil())
			Expect(errors.Is(err, context.Canceled)).To(BeTrue())
		})
	})
})
