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

var HttpDurationsHistogram = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_durations_histogram_seconds",
		Buckets: []float64{0.2, 0.5, 1, 2, 5, 10, 30},
	},
	[]string{"path"},
)
var HttpDurations = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Name:       "http_durations_seconds",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{"path"},
)

func init() {
	prometheus.MustRegister(AccessCounter)
	prometheus.MustRegister(QueueGauge)
	prometheus.MustRegister(HttpDurationsHistogram)
	prometheus.MustRegister(HttpDurations)
}

func Prometheus(r *http.Request) {
	AccessCounter.With(prometheus.Labels{
		"method": r.Method,
		"path":   r.RequestURI,
	}).Add(1)
}
