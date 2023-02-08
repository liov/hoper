package httpi

import (
	_ "expvar"
	"net/http"
	_ "net/http/pprof"
	"runtime/debug"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartPprof() {
	http.ListenAndServe("0.0.0.0:8080", nil)
}

func Debug() http.Handler {
	http.Handle("/debug/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(debug.Stack())
	}))
	return http.DefaultServeMux
}

func PromHandler() http.Handler {
	http.Handle("/metrics", promhttp.Handler())
	return http.DefaultServeMux
}
