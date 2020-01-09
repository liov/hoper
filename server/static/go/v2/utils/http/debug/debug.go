package debug

import (
	_ "expvar"
	"net/http"
	_ "net/http/pprof"
)

func Debug() http.Handler {
	return http.DefaultServeMux
}
