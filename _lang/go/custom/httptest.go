package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {
	port := flag.String("p", "8080", "端口")
	http.DefaultServeMux.HandleFunc("/a/b/{version}/{name}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"hello\": \"world\"}"))
	})
	fmt.Println("server listening on:" + *port)
	http.ListenAndServe(":"+*port, http.DefaultServeMux)
}
