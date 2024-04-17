// Copyright (c) 2019 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/bborbe/assert"
	"github.com/bborbe/run"
)

func TestRunOneTime(t *testing.T) {
	counter := 0
	b := NewOneTimeCron(run.Func(func(ctx context.Context) error {
		counter++
		return nil
	}))
	err := b.Run(context.Background())
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(counter, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestRunOneTimeError(t *testing.T) {
	b := NewOneTimeCron(run.Func(func(ctx context.Context) error {
		return errors.New("fail")
	}))
	err := b.Run(context.Background())
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestRunContinuous(t *testing.T) {
	counter := 0
	b := NewWaitCron(time.Microsecond, run.Func(func(ctx context.Context) error {
		counter++
		return nil
	}))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()
	err := b.Run(ctx)
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(err, Is(context.Canceled)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(counter, Gt(1)); err != nil {
		t.Fatal(err)
	}
}

func TestRunContinuousError(t *testing.T) {
	b := NewWaitCron(time.Microsecond, run.Func(func(ctx context.Context) error {
		return errors.New("fail")
	}))
	err := b.Run(context.Background())
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestRunContinuousCancel(t *testing.T) {
	b := NewWaitCron(time.Microsecond, run.Func(func(ctx context.Context) error {
		<-ctx.Done()
		return nil
	}))
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()
	err := b.Run(ctx)
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(err, Is(context.Canceled)); err != nil {
		t.Fatal(err)
	}
}

func TestExpression(t *testing.T) {
	var counter int64
	b := NewExpressionCron("* * * * * ?", run.Func(func(ctx context.Context) error {
		atomic.AddInt64(&counter, 1)
		return nil
	}))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		time.Sleep(1000 * time.Millisecond)
		cancel()
	}()
	err := b.Run(ctx)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(atomic.LoadInt64(&counter), Ge(int64(1))); err != nil {
		t.Fatal(err)
	}
}

func TestExpressionCancel(t *testing.T) {
	b := NewExpressionCron("* * * * * ?", run.Func(func(ctx context.Context) error {
		<-ctx.Done()
		return nil
	}))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()
	err := b.Run(ctx)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestExpressionError(t *testing.T) {
	b := NewExpressionCron("* * * * * ?", run.Func(func(ctx context.Context) error {
		return errors.New("failed")
	}))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		time.Sleep(1000 * time.Millisecond)
		cancel()
	}()
	err := b.Run(ctx)
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
