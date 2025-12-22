package dbmetrics

import "github.com/prometheus/client_golang/prometheus"

// API 전체 응답 시간
var ApiDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "api",
		Subsystem: "http",
		Name:      "duration_seconds",
		Help:      "HTTP API request duration.",
		Buckets:   []float64{0.05, 0.1, 0.2, 0.5, 1, 2},
	},
	[]string{"api"},
)

// DB 최대 쿼리 시간 (기존 것)
var QueryMaxDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "db",
		Subsystem: "query",
		Name:      "max_duration_seconds",
		Help:      "Max DB query duration per request.",
		Buckets:   []float64{0.01, 0.05, 0.1, 0.2, 0.5, 1},
	},
	[]string{"api"},
)

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
