package metrics

import "github.com/prometheus/client_golang/prometheus"

// HTTP 요청 수
var HttpRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	},
	[]string{"method", "path", "status"},
)

// HTTP 에러 수 (4xx, 5xx)
var HttpErrorsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_errors_total",
		Help: "Total number of HTTP error responses",
	},
	[]string{"method", "path", "status"},
)
