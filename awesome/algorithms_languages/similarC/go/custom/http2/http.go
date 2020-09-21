package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

//浏览器访问降级1.1
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "Hello world")
	})
	h2Handler := h2c.NewHandler(mux, &http2.Server{})
	server := &http.Server{Addr: ":3999", Handler: h2Handler}
	server.ListenAndServe()
}
