package httpi

import (
	_ "expvar"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"
	"runtime/debug"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Debug() http.Handler {
	http.Handle("/debug/metrics", promhttp.Handler())
	http.Handle("/debug/pprof", pprof.Handler("debug"))
	http.Handle("/debug/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(debug.Stack())
	}))
	return http.DefaultServeMux
}
