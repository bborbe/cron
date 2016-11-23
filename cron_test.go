package cron

import (
	"testing"
	. "github.com/bborbe/assert"
	"time"
	"context"
)

func TestRun(t *testing.T) {
	counter := 0
	b := New(true, time.Second, func() error {
		counter++
		return nil
	})
	b.Run(context.Background())
	if err := AssertThat(counter, Is(1)); err != nil {
		t.Fatal(err)
	}
}
