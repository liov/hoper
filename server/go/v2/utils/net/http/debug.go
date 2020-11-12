package http2

import (
	_ "expvar"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"
	"runtime/debug"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Debug() http.Handler {
	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/pprof", pprof.Handler("debug"))
	http.Handle("/debug", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		debug.PrintStack()
	}))
	return http.DefaultServeMux
}
