package main

import (
	"expvar"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests received.",
		},
		[]string{"method", "endpoint", "status"},
	)

	responseDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_duration_seconds",
			Help:    "Duration of HTTP responses in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	totalRequests = expvar.NewInt("total_requests")
	errorRequests = expvar.NewInt("error_requests")
)

// Register metrics
func init() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(responseDuration)
}
