package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
)

var AccessCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "api_requests_total",
	},
	[]string{"method", "path"},
)

var QueueGauge = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "queue_num_total",
	},
	[]string{"name"},
)

func init() {
	prometheus.MustRegister(AccessCounter)
}

func Prometheus(r *http.Request) {
	AccessCounter.With(prometheus.Labels{
		"method": r.Method,
		"path":   r.RequestURI,
	}).Add(1)
}
