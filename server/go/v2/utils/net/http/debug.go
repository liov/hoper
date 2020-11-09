package http

import (
	_ "expvar"
	"net/http"
	_ "net/http/pprof"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Debug() http.Handler {
	http.Handle("/metrics", promhttp.Handler())
	return http.DefaultServeMux
}
