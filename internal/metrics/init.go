package metrics

import (
	"study/pkg/dbmetrics"

	"github.com/prometheus/client_golang/prometheus"
)

func Init() {
	prometheus.MustRegister(HttpRequestsTotal)
	prometheus.MustRegister(HttpErrorsTotal)

	prometheus.MustRegister(dbmetrics.QueryMaxDuration)
}
