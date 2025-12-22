package metrics

import (
	"study/pkg/dbmetrics"

	"github.com/prometheus/client_golang/prometheus"
)

func Init() {
	prometheus.MustRegister(
		dbmetrics.HttpRequestsTotal,
		dbmetrics.HttpErrorsTotal,
		dbmetrics.ApiDuration,
		dbmetrics.QueryMaxDuration,
	)
}
