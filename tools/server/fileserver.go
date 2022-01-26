package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("/home")))

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Println(err)
	}
}
