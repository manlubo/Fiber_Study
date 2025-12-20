package dbmetrics

import "github.com/prometheus/client_golang/prometheus"

var QueryMaxDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "db_query_max_duration_seconds",
		Help: "Max DB query duration per API request",
		Buckets: []float64{
			0.01, // 10ms
			0.05, // 50ms
			0.1,  // 100ms
			0.25, // 250ms
			0.5,  // 500ms
			1,    // 1s
			2,    // 2s
			5,    // 5s
		},
	},
	[]string{"api"},
)
