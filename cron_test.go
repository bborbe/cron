package cron

import (
	"context"
	. "github.com/bborbe/assert"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	counter := 0
	b := New(true, time.Second, func(ctx context.Context) error {
		counter++
		return nil
	})
	err := b.Run(context.Background())
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(counter, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestRunContinuous(t *testing.T) {
	counter := 0
	b := New(false, time.Microsecond, func(ctx context.Context) error {
		counter++
		return nil
	})
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()
	err := b.Run(ctx)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(counter, Gt(1)); err != nil {
		t.Fatal(err)
	}
}
