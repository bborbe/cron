// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron_test

import (
	"context"
	"sync"
	"time"

	"github.com/bborbe/cron"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

var _ = Describe("", func() {
	var ctx context.Context
	var action func(ctx context.Context) error
	var cancel context.CancelFunc
	var testCompleted context.Context
	var testCompletedCancel context.CancelFunc

	var mux sync.Mutex
	var err error
	var running bool

	BeforeEach(func() {
		testCompleted, testCompletedCancel = context.WithCancel(context.Background())
		ctx, cancel = context.WithCancel(context.Background())
	})
	JustBeforeEach(func() {
		var started sync.WaitGroup
		started.Add(1)
		actionCopy := action
		expressionCron := cron.NewExpressionCron("@every 0s", func(ctx context.Context) error {
			cancel()
			started.Done()
			return actionCopy(ctx)
		})

		mux.Lock()
		running = true
		err = nil
		mux.Unlock()

		go func() {
			result := expressionCron.Run(ctx)

			mux.Lock()
			defer mux.Unlock()

			err = result
			running = false
		}()

		started.Wait()
		time.Sleep(100 * time.Millisecond)
	})
	AfterEach(func() {
		testCompletedCancel()
	})
	Context("cancel is stuck", func() {
		BeforeEach(func() {
			done := testCompleted.Done()
			action = func(ctx context.Context) error {
				select {
				case <-time.NewTimer(time.Hour).C:
				case <-done:
				}
				return nil
			}
		})
		It("cron is running", func() {
			mux.Lock()
			defer mux.Unlock()

			Expect(running).To(BeTrue())
		})
		It("returns no error", func() {
			mux.Lock()
			defer mux.Unlock()

			Expect(err).To(BeNil())
		})
	})
	Context("cancel not stuck", func() {
		BeforeEach(func() {
			action = func(ctx context.Context) error {
				return nil
			}
		})
		It("returns no error", func() {
			mux.Lock()
			defer mux.Unlock()

			Expect(err).To(BeNil())
		})
		It("cron has completed", func() {
			mux.Lock()
			defer mux.Unlock()

			Expect(running).To(BeFalse())
		})
	})
	Context("cancel not stuck", func() {
		BeforeEach(func() {
			action = func(ctx context.Context) error {
				<-ctx.Done()
				return ctx.Err()
			}
		})
		It("returns error", func() {
			mux.Lock()
			defer mux.Unlock()

			Expect(errors.Cause(err)).To(Equal(context.Canceled))
		})
		It("cron has completed", func() {
			mux.Lock()
			defer mux.Unlock()

			Expect(running).To(BeFalse())
		})
	})
})
