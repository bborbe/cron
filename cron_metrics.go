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
)

func init() {
	prometheus.DefaultRegisterer.MustRegister(
		started,
		completed,
		failed,
		lastSuccess,
	)
}

//counterfeiter:generate -o ../mocks/cron-metrics.go --fake-name CronMetrics . Metrics
type Metrics interface {
	IncreaseStarted(name string)
	IncreaseFailed(name string)
	IncreaseCompleted(name string)
	SetLastSuccessToCurrent(name string)
}

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
