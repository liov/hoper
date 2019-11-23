package main

import (
	"net/http"
	_ "net/http/pprof"
)

//http://localhost:8080/debug/pprof/
func main() {
	http.ListenAndServe("0.0.0.0:8080", nil)
	select {}
}
