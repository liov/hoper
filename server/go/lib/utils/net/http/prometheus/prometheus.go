package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"math/rand"
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
	[]string{"method", "path"},
)

var HttpDurationsHistogram = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_durations_histogram_millisecond",
		Buckets: []float64{30, 60, 100, 200, 300, 500, 1000},
	},
	[]string{"method", "path"},
)
var HttpDurations = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Name:       "http_durations_millisecond",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{"method", "path"},
)

func init() {
	prometheus.MustRegister(AccessCounter)
	prometheus.MustRegister(QueueGauge)
	prometheus.MustRegister(HttpDurationsHistogram)
	prometheus.MustRegister(HttpDurations)
}

func Prometheus(r *http.Request) {
	labels := prometheus.Labels{
		"method": r.Method,
		"path":   r.RequestURI,
	}
	AccessCounter.With(labels).Add(1)
	QueueGauge.With(labels).Set(float64(rand.Intn(1000)))
	HttpDurationsHistogram.With(labels).Observe(float64(rand.Intn(1000)))
	HttpDurations.With(labels).Observe(float64(rand.Intn(1000)))
}
