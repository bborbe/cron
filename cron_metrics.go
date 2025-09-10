// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cron

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	started = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "cron",
		Subsystem: "job",
		Name:      "started",
		Help:      "Number of times cron job was started",
	}, []string{"name"})
	completed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "cron",
		Subsystem: "job",
		Name:      "completed",
		Help:      "Number of times cron job completed successfully",
	}, []string{"name"})
	failed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "cron",
		Subsystem: "job",
		Name:      "failed",
		Help:      "Number of times cron job failed",
	}, []string{"name"})
	lastSuccess = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "cron",
		Subsystem: "job",
		Name:      "last_success",
		Help:      "Timestamp of last successful run",
	}, []string{"name"})
	duration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "cron",
		Subsystem: "job",
		Name:      "duration_seconds",
		Help:      "Duration of cron job execution in seconds",
		Buckets:   prometheus.DefBuckets,
	}, []string{"name"})
)

func init() {
	prometheus.DefaultRegisterer.MustRegister(
		started,
		completed,
		failed,
		lastSuccess,
		duration,
	)
}

//counterfeiter:generate -o ../mocks/cron-metrics.go --fake-name CronMetrics . Metrics

// Metrics provides methods for collecting and reporting cron job execution statistics.
type Metrics interface {
	// IncreaseStarted increments the counter for cron job start events.
	IncreaseStarted(name string)
	// IncreaseFailed increments the counter for cron job failure events.
	IncreaseFailed(name string)
	// IncreaseCompleted increments the counter for cron job completion events.
	IncreaseCompleted(name string)
	// SetLastSuccessToCurrent records the timestamp of the last successful execution.
	SetLastSuccessToCurrent(name string)
	// ObserveDuration records the execution duration in seconds for the named cron job.
	ObserveDuration(name string, durationSeconds float64)
}

// NewMetrics creates a new Metrics instance that reports to Prometheus.
func NewMetrics() Metrics {
	return &metrics{}
}

type metrics struct {
}

func (c *metrics) IncreaseStarted(name string) {
	started.With(prometheus.Labels{"name": name}).Inc()
}

func (c *metrics) IncreaseFailed(name string) {
	failed.With(prometheus.Labels{"name": name}).Inc()
}

func (c *metrics) IncreaseCompleted(name string) {
	completed.With(prometheus.Labels{"name": name}).Inc()
}

func (c *metrics) SetLastSuccessToCurrent(name string) {
	lastSuccess.With(prometheus.Labels{"name": name}).SetToCurrentTime()
}

func (c *metrics) ObserveDuration(name string, durationSeconds float64) {
	duration.With(prometheus.Labels{"name": name}).Observe(durationSeconds)
}
