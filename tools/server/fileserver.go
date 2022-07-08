package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("/")))
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			log.Printf("%s: %s", k, v)
		}
		w.Write([]byte("test"))
	})

	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		log.Println(err)
	}
}
