package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
		},
		[]string{"pattern", "method", "status"}, //label
	)

	HttpRequestsDurationHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_requests_duration_seconds_histogram",
			Buckets: []float64{
				0.1,  //100ms
				0.2,  //200ms
				0.25, //250ms
				0.5,  //500ms
				1,    //1s
			},
		},
		[]string{"pattern", "method"},
	)

	HttpRequestsDurationSummary = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_requests_duration_seconds_summary",
			Objectives: map[float64]float64{
				0.99: 0.001, // 0.99 +- 0.001
				0.95: 0.01,  // 0.95 +- 0.01
				0.5:  0.05,  // 0.5 +- 0.05
			},
		},
		[]string{"pattern", "method"},
	)
)
