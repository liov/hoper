package server

import (
	"fmt"
	"net/http"
)

func Http() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(w, "hello")
	})
	return mux
}
